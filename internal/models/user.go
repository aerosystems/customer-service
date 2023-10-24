package models

import (
	"time"
)

type User struct {
	Id           uint      `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	Email        string    `json:"email" gorm:"unique"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"-"`
	GoogleId     string    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

func NewUserEntity(email, passwordHash, role string) *User {
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	}
}

type UserRepository interface {
	FindAll() (*[]User, error)
	FindById(Id uint) (*User, error)
	FindByEmail(Email string) (*User, error)
	FindByGoogleId(GoogleId string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(user *User) error
}
