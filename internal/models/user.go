package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id        uint      `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"-"`
	GoogleId  string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func NewUser(email, password, role string) *User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return &User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}
}

type UserRepository interface {
	FindAll() (*[]User, error)
	FindByID(ID int) (*User, error)
	FindByEmail(Email string) (*User, error)
	FindByGoogleID(GoogleID string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(user *User) error
	ResetPassword(user *User, password string) error
	PasswordMatches(user *User, plainText string) (bool, error)
}
