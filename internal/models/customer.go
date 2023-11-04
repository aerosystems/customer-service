package models

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	Id        uint      `json:"-" gorm:"primaryKey;unique;autoIncrement"`
	Uuid      uuid.UUID `json:"uuid" gorm:"unique"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type CustomerRepository interface {
	GetAll() (*[]Customer, error)
	GetById(Id uint) (*Customer, error)
	Create(user *Customer) error
	Update(user *Customer) error
	Delete(user *Customer) error
}
