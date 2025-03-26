package main

import (
	"context"
	"fmt"
	"log"
	"sqs-fifo/pkg/config"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx := context.Background()

	// Initialize connection settings from config
	cfg, err := config.LoadAWSConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	// List SQS queues
	sqsClient := sqs.NewFromConfig(cfg)
	sqsOut, err := sqsClient.ListQueues(ctx, &sqs.ListQueuesInput{})
	if err != nil {
		log.Fatalf("failed to list SQS queues: %v", err)
	}

	fmt.Println("\n🟦 SQS Queues:")
	if len(sqsOut.QueueUrls) == 0 {
		fmt.Println("  (no queues found)")
	} else {
		for _, url := range sqsOut.QueueUrls {
			fmt.Println("  -", url)
		}
	}

	// List SNS topics
	snsClient := sns.NewFromConfig(cfg)
	snsOut, err := snsClient.ListTopics(ctx, &sns.ListTopicsInput{})
	if err != nil {
		log.Fatalf("failed to list SNS topics: %v", err)
	}

	fmt.Println("\n🟨 SNS Topics:")
	if len(snsOut.Topics) == 0 {
		fmt.Println("  (no topics found)")
	} else {
		for _, topic := range snsOut.Topics {
			fmt.Println("  -", *topic.TopicArn)
		}
	}
}
