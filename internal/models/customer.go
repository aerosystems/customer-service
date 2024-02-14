package models

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	Id        int       `json:"-" gorm:"primaryKey;unique;autoIncrement"`
	Uuid      uuid.UUID `json:"uuid" gorm:"unique"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
