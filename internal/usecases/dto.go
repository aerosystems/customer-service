package usecases

import "time"

type SubscriptionDTO struct {
	UUID        string
	Type        string
	AccessTime  time.Time
	AccessCount int64
}
