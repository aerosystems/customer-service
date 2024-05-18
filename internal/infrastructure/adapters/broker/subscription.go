package broker

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aerosystems/customer-service/internal/models"
	PubSub "github.com/aerosystems/customer-service/pkg/pubsub"
	"github.com/google/uuid"
	"time"
)

const (
	defaultTimeout = 2 * time.Second
)

type SubscriptionEventsAdapter struct {
	pubsubClient *PubSub.Client
	topicId      string
}

func NewSubscriptionEventsAdapter(pubsubClient *PubSub.Client, topicId string) *SubscriptionEventsAdapter {
	return &SubscriptionEventsAdapter{
		pubsubClient: pubsubClient,
		topicId:      topicId,
	}
}

type CreateSubscriptionEvent struct {
	CustomerUuid         string `json:"customerUuid"`
	SubscriptionType     string `json:"subscriptionType"`
	SubscriptionDuration string `json:"subscriptionDuration"`
}

func (s SubscriptionEventsAdapter) PublishCreateSubscriptionEvent(
	customerUuid uuid.UUID,
	subscriptionType models.SubscriptionType,
	subscriptionDuration models.SubscriptionDuration,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	topic := s.pubsubClient.Client.Topic(s.topicId)
	ok, err := topic.Exists(ctx)
	defer topic.Stop()
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}
	if !ok {
		if _, err := s.pubsubClient.Client.CreateTopic(ctx, s.topicId); err != nil {
			return fmt.Errorf("failed to create topic: %w", err)
		}
	}

	event := CreateSubscriptionEvent{
		CustomerUuid:         customerUuid.String(),
		SubscriptionType:     subscriptionType.String(),
		SubscriptionDuration: subscriptionDuration.String(),
	}
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal create subscription event: %w", err)
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: eventData,
	})
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("failed to publish create subscription event: %w", err)
	}

	return nil
}
