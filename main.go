package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	useLocalstack = true
	groupSpecific = false

	numGroups        = 6
	messagesPerGroup = 10

	// configure number of consumers for different modes
	group1numConsumers  = 1 // consumers for group-1 in group-specific mode
	generalNumConsumers = 3 // total consumers when not group-specific
)

var (
	queueURL string
	topicARN string
	region   string
	profile  string
)

func generateGroups() []string {
	groups := make([]string, numGroups)
	for i := 0; i < numGroups; i++ {
		groups[i] = fmt.Sprintf("group-%d", i+1)
	}
	return groups
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	region = "us-east-1"
	if useLocalstack {
		queueURL = "http://localhost:4566/000000000000/demo-queue.fifo"
		topicARN = "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo"
	} else {
		queueURL = os.Getenv("AWS_QUEUE_URL")
		topicARN = os.Getenv("AWS_TOPIC_ARN")
		profile = os.Getenv("AWS_PROFILE")
	}
}

func main() {
	fmt.Println("Starting FIFO Producer and Consumers...")

	ctx := context.Background()

	var cfg aws.Config
	var err error

	if useLocalstack {
		fmt.Println("Using Localstack for development")
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
			config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "http://localhost:4566",
					SigningRegion: region,
				}, nil
			})),
		)
	} else {
		fmt.Println("Connecting to AWS using profile")
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithSharedConfigProfile(profile),
		)
	}

	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)
	snsClient := sns.NewFromConfig(cfg)

	groups := generateGroups()
	// sending messages per group
	for _, groupID := range groups {
		for i := 0; i < messagesPerGroup; i++ {
			sendToSNS(ctx, snsClient, fmt.Sprintf("%s - Message %d", groupID, i+1), groupID)
		}
	}

	var wg sync.WaitGroup

	if groupSpecific {
		// start a specific consumer for each group
		for groupNum := 1; groupNum <= numGroups; groupNum++ {
			groupID := fmt.Sprintf("group-%d", groupNum)

			// for group-1, start configured number of consumers
			if groupNum == 1 {
				for i := 1; i <= group1numConsumers; i++ {
					wg.Add(1)
					go func(id int, gID string) {
						fmt.Printf("[Consumer] Starting consumer %d for %s\n", id, gID)
						consumeOnlyGroup(ctx, sqsClient, id, gID, &wg)
					}(i, groupID)
				}
				continue
			}

			// for other groups, start one consumer each
			wg.Add(1)
			go func(id int, gID string) {
				fmt.Printf("[Consumer] Starting consumer %d for %s\n", id, gID)
				consumeOnlyGroup(ctx, sqsClient, id, gID, &wg)
			}(groupNum, groupID)
		}
	} else {
		// start the non group specific consumers
		for i := 1; i <= generalNumConsumers; i++ {
			wg.Add(1)
			go func(id int) {
				fmt.Printf("[Consumer] Starting general consumer %d (total consumers: %d)\n",
					id, generalNumConsumers)
				consumeAllMessages(ctx, sqsClient, id, &wg)
			}(i)
		}
	}

	wg.Wait()
}

func sendToSNS(ctx context.Context, client *sns.Client, body string, groupID string) {
	fmt.Println("[Producer] Sending message to SNS topic:", topicARN, "GroupId:", groupID, "Body:", body)

	_, err := client.Publish(ctx, &sns.PublishInput{
		Message:                aws.String(body),
		TopicArn:               aws.String(topicARN),
		MessageGroupId:         aws.String(groupID),
		MessageDeduplicationId: aws.String(fmt.Sprintf("%s-%d", groupID, time.Now().UnixNano())),
	})
	if err != nil {
		log.Printf("[Producer] failed to send message to %s: %v", topicARN, err)
	}
}

func consumeAllMessages(ctx context.Context, client *sqs.Client, consumerID int, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()

	for {
		resp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     2,
		})
		if err != nil {
			log.Printf("[Consumer %d] Error receiving: %v", consumerID, err)
			continue
		}
		if len(resp.Messages) == 0 {
			elapsed := time.Since(start)
			log.Printf("[Consumer %d] Finished processing all messages in %v", consumerID, elapsed)
			return
		}

		msg := resp.Messages[0]
		log.Printf("[Consumer %d] Processing: %s", consumerID, *msg.Body)

		// processing the message
		time.Sleep(2 * time.Second)

		_, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queueURL),
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			log.Printf("[Consumer %d] Error deleting message: %v", consumerID, err)
		} else {
			log.Printf("[Consumer %d] Deleted: %s", consumerID, *msg.Body)
		}
	}
}

func consumeOnlyGroup(ctx context.Context, client *sqs.Client, consumerID int, expectedGroup string, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()

	for {
		resp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     2,
		})
		if err != nil {
			log.Printf("[Consumer %d][%s] Error receiving: %v", consumerID, expectedGroup, err)
			continue
		}
		if len(resp.Messages) == 0 {
			elapsed := time.Since(start)
			log.Printf("[Consumer %d][%s] Finished processing all messages in %v", consumerID, expectedGroup, elapsed)
			return // there are no more messages to read
		}

		msg := resp.Messages[0]

		// manual filtering, only processing if the message belongs to the expected group
		if !isGroupMessage(msg, expectedGroup) {
			log.Printf("[Consumer %d][%s] Ignoring message from another group: %s", consumerID, expectedGroup, *msg.Body)
			_, err = client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
				QueueUrl:          aws.String(queueURL),
				ReceiptHandle:     msg.ReceiptHandle,
				VisibilityTimeout: int32(0), // makes it visible again immediately
			})
			if err != nil {
				log.Printf("[Consumer %d][%s] Error changing visibility: %v", consumerID, expectedGroup, err)
			}
			continue
		}

		log.Printf("[Consumer %d][%s] Processing: %s", consumerID, expectedGroup, *msg.Body)

		// processing the message
		time.Sleep(2 * time.Second)

		_, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queueURL),
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			log.Printf("[Consumer %d][%s] Error deleting message: %v", consumerID, expectedGroup, err)
		} else {
			log.Printf("[Consumer %d][%s] Deleted: %s", consumerID, expectedGroup, *msg.Body)
		}
	}
}

// simple check based on group name in message body (group ID is not returned by ReceiveMessage)
// supposedly, it should be on the msg metadata map, but didn't find it. Maybe due to LocalStack limitations.
func isGroupMessage(msg types.Message, group string) bool {
	return msg.Body != nil && len(*msg.Body) > 0 && contains(*msg.Body, group)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
