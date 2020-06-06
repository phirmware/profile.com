package models

import (
	"github.com/jinzhu/gorm"
)

// Services defines the shape of the service struct
type Services struct {
	db   *gorm.DB
	User UserService
}

// NewServices is used to define the service shape
func NewServices(connectionString string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionString)
	db.LogMode(true)
	userService := NewUserService(db)
	if err != nil {
		return nil, err
	}
	return &Services{
		User: userService,
		db:   db,
	}, nil
}

// AutoMigrate creates the tables in the database
func (s *Services) AutoMigrate() error {
	if err := s.db.AutoMigrate(User{}).Error; err != nil {
		return err
	}
	return nil
}

// DestructiveConstruct destroys db and recreates
func (s *Services) DestructiveConstruct() error {
	if err := s.db.DropTableIfExists(User{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}
