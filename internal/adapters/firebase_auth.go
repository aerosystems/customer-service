package adapters

import (
	"context"
	"firebase.google.com/go/auth"
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

func (f *FirebaseAuthAdapter) SetClaimCustomerUUID(ctx context.Context, uid string, customerUUID uuid.UUID) error {
	claims := map[string]interface{}{"customer_uuid": customerUUID.String()}
	return f.Client.SetCustomUserClaims(ctx, uid, claims)
}
