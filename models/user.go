package models

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"

	"profile.com/hash"
	"profile.com/rand"
)

var (
	// ErrInternalServerError is returned when err cannot be determined
	ErrInternalServerError = errors.New("models: Something went wrong, contact for help")
	// ErrNameMissing is returned when the user fails to input a name
	ErrNameMissing = errors.New("models: Please provide your name")
	// ErrEmailMissing is returned when the user fails to input an email
	ErrEmailMissing = errors.New("models: Please input your email")
	// ErrEmailTaken is returned when email is already in use
	ErrEmailTaken = errors.New("models: This email is already in use")
	// ErrInvalidCredentials is returned after an invalid login attempt
	ErrInvalidCredentials = errors.New("models: Invalid login credentials")
	// ErrPasswordTooShort is returned when user inputs short password
	ErrPasswordTooShort = errors.New("models: The password you provided is too short, minimum of 8 characters")
	// ErrPasswordNotProvided is returned when user doesnt provide a pasword
	ErrPasswordNotProvided = errors.New("models: Please provide a password")
	// ErrPasswordInvalid is returned when a user uses an invalid password
	ErrPasswordInvalid = errors.New("models: Invalid Password, try again")
	// ErrPasswordHashMissing is returned when a password hash is missing
	ErrPasswordHashMissing = errors.New("models: No password hash")
	// ErrRememberMissing is returned when there is no remember field set
	ErrRememberMissing = errors.New("models: Remember is missing")
)

const (
	pepper = "secret-user-pepper"
	key    = "secret-key"
)

// User defines the shape of the user db
type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"not null"`
	RememberHash string
	Title        string
	Summary      string
	Skills       string
}

// UserDB defines the shape of the userdb interface
type UserDB interface {
	Create(user *User) error
	ByEmail(email string) (*User, error)
	Update(user *User) error
	ByRemember(rememberToken string) (*User, error)
	All() (*[]User, error)
}

// UserVal is the user validation interface
type UserVal interface {
	Authenticate(user *User) (*User, error)
	UserDB
}

// UserService defines the shape of the userservice
type UserService interface {
	UserVal
}

type userService struct {
	UserVal
}
type userValidation struct {
	UserDB
	hmac hash.HMAC
}
type userGorm struct {
	db *gorm.DB
}

// NewUserService returns the userservice struct
func NewUserService(db *gorm.DB) UserService {
	ug := newUserGorm(db)
	uv := newUserValidation(ug)
	return &userService{
		UserVal: uv,
	}
}

func newUserValidation(ug *userGorm) *userValidation {
	hmac := hash.NewHMAC(key)
	return &userValidation{
		hmac:   hmac,
		UserDB: ug,
	}
}

func newUserGorm(db *gorm.DB) *userGorm {
	return &userGorm{
		db: db,
	}
}

// ##################### User Service ################################ //

// ##################### User Validation ################################ //

type userValFn func(user *User) error

func runUserValFns(user *User, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidation) checkForName(user *User) error {
	if user.Name == "" {
		return ErrNameMissing
	}
	return nil
}

func (uv *userValidation) checkForEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailMissing
	}
	return nil
}

func (uv *userValidation) checkForPassword(user *User) error {
	if user.Password == "" {
		return ErrPasswordNotProvided
	}
	return nil
}

func (uv *userValidation) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidation) checkDBForEmail(user *User) error {
	_, err := uv.UserDB.ByEmail(user.Email)
	if err != nil {
		return nil
	}
	return ErrEmailTaken
}

func (uv *userValidation) checkPasswordLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidation) hashPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	passwordPepper := user.Password + pepper
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordPepper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(bytes)
	user.Password = ""
	return nil
}

func (uv *userValidation) checkForPasswordHash(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordHashMissing
	}
	return nil
}

func (uv *userValidation) generateRemember(user *User) error {
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidation) rememberHash(user *User) error {
	if user.Remember == "" {
		return ErrRememberMissing
	}
	hash := uv.hmac.Hash(user.Remember)
	user.RememberHash = hash
	return nil
}

func (uv *userValidation) Create(user *User) error {
	if err := runUserValFns(user,
		uv.checkForName,
		uv.checkForEmail,
		uv.checkForPassword,
		uv.checkPasswordLength,
		uv.normalizeEmail,
		uv.checkDBForEmail,
		uv.hashPassword,
		uv.checkForPasswordHash,
		uv.generateRemember,
		uv.rememberHash,
	); err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidation) ByEmail(email string) (*User, error) {
	user := &User{
		Email: email,
	}
	if err := runUserValFns(user,
		uv.checkForEmail,
		uv.normalizeEmail,
	); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

func (uv *userValidation) Update(user *User) error {
	if err := runUserValFns(user,
		uv.checkForName,
		uv.checkForEmail,
		uv.normalizeEmail,
		uv.checkPasswordLength,
		uv.hashPassword,
		uv.checkForPasswordHash,
	); err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidation) Authenticate(user *User) (*User, error) {
	if err := runUserValFns(user,
		uv.checkForEmail,
		uv.normalizeEmail,
		uv.checkForPassword,
	); err != nil {
		return nil, err
	}
	u, err := uv.UserDB.ByEmail(user.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(user.Password+pepper)); err != nil {
		return nil, ErrPasswordInvalid
	}
	return u, nil
}

// ##################### User Gorm ################################ //

func (ug *userGorm) All() (*[]User, error) {
	users := []User{}
	if err := ug.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (ug *userGorm) Create(user *User) error {
	if err := ug.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	user := &User{}
	err := ug.db.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ug *userGorm) ByRemember(rememberToken string) (*User, error) {
	var user User
	if err := ug.db.Where("remember = ?", rememberToken).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}
