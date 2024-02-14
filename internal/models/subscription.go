package models

import (
	"github.com/google/uuid"
	"time"
)

type SubsRPCPayload struct {
	UserUuid   uuid.UUID
	Kind       KindSubscription
	AccessTime time.Time
}

type KindSubscription string

const (
	TrialSubscription    KindSubscription = "trial"
	StartupSubscription  KindSubscription = "startup"
	BusinessSubscription KindSubscription = "business"
)

func (k KindSubscription) String() string {
	return string(k)
}
