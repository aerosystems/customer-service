package consumer

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	PubSub "github.com/aerosystems/customer-service/pkg/pubsub"
	"github.com/sirupsen/logrus"
)

type AuthSubscription struct {
	log             *logrus.Logger
	pubsubClient    *PubSub.Client
	topicId         string
	subName         string
	customerUsecase CustomerUsecase
}

func NewAuthSubscription(log *logrus.Logger, client *PubSub.Client, topicId, subName string, customerUsecase CustomerUsecase) *AuthSubscription {
	return &AuthSubscription{
		log:             log,
		pubsubClient:    client,
		topicId:         topicId,
		subName:         subName,
		customerUsecase: customerUsecase,
	}
}

type CustomerData struct {
	Uuid string `json:"uuid"`
}

func (s AuthSubscription) Run() error {
	// Create a topic if it doesn't already exist
	topic := s.pubsubClient.Client.Topic(s.topicId)
	ok, err := topic.Exists(s.pubsubClient.Ctx)
	if err != nil {
		s.log.Errorf("Failed to check if topic exists: %v", err)
	}
	if !ok {
		if _, err := s.pubsubClient.Client.CreateTopic(s.pubsubClient.Ctx, s.topicId); err != nil {
			s.log.Errorf("Failed to create topic: %v", err)
		}
		s.log.Infof("Topic %s created.\n", s.topicId)
	}
	// Create a subscription to the topic
	sub := s.pubsubClient.Client.Subscription(s.topicId)
	ok, err = sub.Exists(s.pubsubClient.Ctx)
	if err != nil {
		s.log.Errorf("Failed to check if subscription exists: %v", err)
	}
	if !ok {
		if _, err := s.pubsubClient.Client.CreateSubscription(s.pubsubClient.Ctx, s.subName, pubsub.SubscriptionConfig{
			Topic: topic,
		}); err != nil {
			s.log.Errorf("Failed to create subscription: %v", err)
		}
		s.log.Infof("Subscription %s created.\n", s.subName)
	}
	// Start consuming messages from the subscription
	if err := sub.Receive(s.pubsubClient.Ctx, func(ctx context.Context, msg *pubsub.Message) {
		// Unmarshal the message data into a struct
		var customer CustomerData
		if err := json.Unmarshal(msg.Data, &customer); err != nil {
			s.log.Errorf("Failed to unmarshal message data: %v", err)
			msg.Nack()
			return
		}

		_, err := s.customerUsecase.CreateCustomer(customer.Uuid)
		if err != nil {
			s.log.Errorf("Failed to create customer: %v", err)
		}
		// Acknowledge the message
		msg.Ack()
	}); err != nil {
		s.log.Errorf("Failed to receive messages: %v", err)
	}
	return nil
}
