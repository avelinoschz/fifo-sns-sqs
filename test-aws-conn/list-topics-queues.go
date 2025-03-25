package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		profile = "default"
	}

	// load the AWS config using a named profile
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)
	sqsOut, err := sqsClient.ListQueues(ctx, &sqs.ListQueuesInput{})
	if err != nil {
		log.Fatalf("failed to list SQS queues: %v", err)
	}

	fmt.Println("🟦 SQS Queues:")
	if len(sqsOut.QueueUrls) == 0 {
		fmt.Println("  (no queues found)")
	} else {
		for _, url := range sqsOut.QueueUrls {
			fmt.Println("  -", url)
		}
	}

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
