package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

const (
	projectID      = "test-project"
	topicName      = "test-topic"
	subscriptionID = "test-subscription"
)

func main() {
	// Set up context
	ctx := context.Background()

	// Create Pub/Sub client pointing to the emulator
	client, err := pubsub.NewClient(ctx, projectID, option.WithEndpoint("localhost:8085"))
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}
	defer client.Close()

	// Create topic and subscription
	if err := setupTopicAndSubscription(ctx, client); err != nil {
		log.Fatalf("Failed to setup topic and subscription: %v", err)
	}

	// Start subscriber in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := subscribeMessages(ctx, client); err != nil {
			log.Printf("Failed to subscribe: %v", err)
		}
	}()

	// Give subscriber time to start
	time.Sleep(2 * time.Second)

	// Publish some test messages
	messages := []string{
		"Hello from Pub/Sub emulator!",
		"This is message number 2",
		"Testing Pub/Sub functionality",
		"Final test message",
	}

	for i, msg := range messages {
		if err := publishMessage(ctx, client, fmt.Sprintf("Message %d: %s", i+1, msg)); err != nil {
			log.Printf("Failed to publish message %d: %v", i+1, err)
		} else {
			fmt.Printf("Published: Message %d\n", i+1)
		}
		time.Sleep(1 * time.Second)
	}

	// Wait a bit for messages to be processed
	time.Sleep(5 * time.Second)

	fmt.Println("Demo completed!")
}

func setupTopicAndSubscription(ctx context.Context, client *pubsub.Client) error {
	// Create topic
	topic := client.Topic(topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %v", err)
	}

	if !exists {
		topic, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			return fmt.Errorf("failed to create topic: %v", err)
		}
		fmt.Printf("Created topic: %s\n", topicName)
	} else {
		fmt.Printf("Topic '%s' already exists\n", topicName)
	}

	// Create subscription
	sub := client.Subscription(subscriptionID)
	exists, err = sub.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if subscription exists: %v", err)
	}

	if !exists {
		_, err = client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
			Topic: topic,
		})
		if err != nil {
			return fmt.Errorf("failed to create subscription: %v", err)
		}
		fmt.Printf("Created subscription: %s\n", subscriptionID)
	} else {
		fmt.Printf("Subscription '%s' already exists\n", subscriptionID)
	}

	return nil
}

func publishMessage(ctx context.Context, client *pubsub.Client, message string) error {
	topic := client.Topic(topicName)

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(message),
		Attributes: map[string]string{
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})

	// Wait for the publish to complete
	_, err := result.Get(ctx)
	return err
}

func subscribeMessages(ctx context.Context, client *pubsub.Client) error {
	sub := client.Subscription(subscriptionID)

	fmt.Println("Starting message subscriber...")

	// Receive messages
	return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))

		// Print attributes if any
		if len(msg.Attributes) > 0 {
			fmt.Printf("  Attributes: %v\n", msg.Attributes)
		}

		// Acknowledge the message
		msg.Ack()
	})
}
