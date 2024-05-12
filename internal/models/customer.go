package models

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	Uuid      uuid.UUID
	CreatedAt time.Time
}
