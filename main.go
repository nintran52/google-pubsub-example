package main

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
)

const (
	projectID = "ion-staging"
	topicID   = "example-topic"
)

type Message struct {
	Data string `json:"data"`
}

func main() {
	ctx := context.Background()
	// Create a Pub/Sub client
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	// Create a topic if it doesn't already exist
	topic := client.Topic(topicID)
	ok, err := topic.Exists(ctx)
	if err != nil {
		log.Fatalf("Failed to check if topic exists: %v", err)
	}
	if !ok {
		if _, err := client.CreateTopic(ctx, topicID); err != nil {
			log.Fatalf("Failed to create topic: %v", err)
		}
		log.Printf("Topic %s created.\n", topicID)
	}
	// Start pub message
	msgString := make(map[string]string)
	msgString["field1"] = "data-field1"
	var msg []byte
	msg, _ = json.Marshal(msgString)
	attr := make(map[string]string)
	attr["topic"] = "top-name"
	res := topic.Publish(ctx, &pubsub.Message{Data: msg, Attributes: attr})
	_, err = res.Get(ctx)
	if err != nil {
		log.Printf("result.Get: %v", err)
	}
	log.Printf("Pushed data: %s\n", string(msg))
}
