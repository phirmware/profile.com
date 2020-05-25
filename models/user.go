package models

import (
	"errors"
	"fmt"
)

var (
	// ErrNameMissing is returned when the user fails to input a name
	ErrNameMissing = errors.New("models: Please provide your name")
	// ErrEmailMissing is returned when the user fails to input an email
	ErrEmailMissing = errors.New("models: Please input your email")
	// ErrPasswordTooShort is returned when user inputs short password
	ErrPasswordTooShort = errors.New("models: The password you provided is too short, minimum of 8 characters")
	// ErrPasswordNotProvided is returned when user doesnt provide a pasword
	ErrPasswordNotProvided = errors.New("models: Please provide a password")
)

// User defines the shape of the user db
type User struct {
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}

// UserDB defines the shape of the userdb interface
type UserDB interface {
	Create(user *User) error
}

// UserService defines the shape of the userservice
type UserService struct {
	UserDB
}
type userValidation struct {
	UserDB
}
type userGorm struct{}

// NewUserService returns the userservice struct
func NewUserService() *UserService {
	uv := newUserValidation()
	return &UserService{
		UserDB: uv,
	}
}

func newUserValidation() *userValidation {
	return &userValidation{}
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

func (uv *userValidation) checkPasswordLength(user *User) error {
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidation) Create(user *User) error {
	if err := runUserValFns(user,
		uv.checkForName,
		uv.checkForEmail,
		uv.checkForPassword,
		uv.checkPasswordLength,
	); err != nil {
		return err
	}
	fmt.Println(user)
	return nil
}

// ##################### User Gorm ################################ //
