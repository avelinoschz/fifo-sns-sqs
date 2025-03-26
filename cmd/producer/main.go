package main

import (
	"context"
	"fmt"
	"log"
	"sqs-fifo/pkg/config"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func main() {
	fmt.Println("Starting FIFO Producer...")

	ctx := context.Background()

	cfg, err := config.LoadAWSConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	snsClient := sns.NewFromConfig(cfg)

	groups := config.GenerateGroups()
	// sending messages per group
	for _, groupID := range groups {
		for i := 0; i < config.MessagesPerGroup; i++ {
			sendToSNS(ctx, snsClient, fmt.Sprintf("%s - Message %d", groupID, i+1), groupID)
		}
	}

	fmt.Println("All messages have been sent successfully!")
}

func sendToSNS(ctx context.Context, client *sns.Client, body string, groupID string) {
	fmt.Println("[Producer] Sending message to SNS topic:", config.TopicARN, "GroupId:", groupID, "Body:", body)

	_, err := client.Publish(ctx, &sns.PublishInput{
		Message:                aws.String(body),
		TopicArn:               aws.String(config.TopicARN),
		MessageGroupId:         aws.String(groupID),
		MessageDeduplicationId: aws.String(fmt.Sprintf("%s-%d", groupID, time.Now().UnixNano())),
	})
	if err != nil {
		log.Printf("[Producer] failed to send message to %s: %v", config.TopicARN, err)
	}
}
