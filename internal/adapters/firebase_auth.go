package adapters

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
)

type FirebaseAuthAdapter struct {
	Client *auth.Client
}

func NewFirebaseAuthAdapter(client *auth.Client) *FirebaseAuthAdapter {
	return &FirebaseAuthAdapter{
		Client: client,
	}
}

func (f *FirebaseAuthAdapter) SetCustomUserClaims(ctx context.Context, uid string, customerUUID uuid.UUID) error {
	claims := map[string]interface{}{
		"user_uuid": customerUUID.String(),
		"role":      "customer",
	}
	return f.Client.SetCustomUserClaims(ctx, uid, claims)
}
