package broker

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	PubSub "github.com/aerosystems/customer-service/pkg/pubsub"
	"github.com/google/uuid"
	"time"
)

const (
	defaultTimeout = 2 * time.Second
)

type SubscriptionEventsAdapter struct {
	pubsubClient              *PubSub.Client
	topicId                   string
	subName                   string
	createFreeTrialEndpoint   string
	subscriptionServiceApiKey string
}

func NewSubscriptionEventsAdapter(pubsubClient *PubSub.Client, topicId, subName, createFreeTrialEndpoint, subscriptionServiceApiKey string) *SubscriptionEventsAdapter {
	return &SubscriptionEventsAdapter{
		pubsubClient:              pubsubClient,
		topicId:                   topicId,
		subName:                   subName,
		createFreeTrialEndpoint:   createFreeTrialEndpoint,
		subscriptionServiceApiKey: subscriptionServiceApiKey,
	}
}

type CreateSubscriptionEvent struct {
	CustomerUuid string `json:"customerUuid"`
}

func (s SubscriptionEventsAdapter) PublishCreateFreeTrialEvent(
	customerUuid uuid.UUID,
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

	sub := s.pubsubClient.Client.Subscription(s.subName)
	ok, err = sub.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if subscription exists: %w", err)
	}
	if !ok {
		if _, err := s.pubsubClient.Client.CreateSubscription(ctx, s.subName, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 10 * time.Second,
			PushConfig: pubsub.PushConfig{
				Endpoint: s.createFreeTrialEndpoint,
			},
		}); err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}
	}

	event := CreateSubscriptionEvent{
		CustomerUuid: customerUuid.String(),
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
