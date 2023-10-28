package models

import "time"

type Code struct {
	Id        uint      `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	Code      string    `json:"code"`
	UserId    uint      `json:"userId"`
	User      User      `json:"user"`
	Action    string    `json:"action"`
	Data      string    `json:"data"`
	IsUsed    bool      `json:"isUsed"`
	CreatedAt time.Time `json:"createdAt"`
	ExpireAt  time.Time `json:"expireAt"`
}

type CodeRepository interface {
	FindById(Id uint) (*Code, error)
	Create(code *Code) error
	Update(code *Code) error
	Delete(code *Code) error
	GetByCode(value string) (*Code, error)
	GetLastIsActiveCode(UserId uint, Action string) (*Code, error)
	ExtendExpiration(code *Code) error
}
