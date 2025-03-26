package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sqs-fifo/pkg/config"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func main() {
	fmt.Println("Starting FIFO Consumers...")

	// Create context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
		cancel()
	}()

	cfg, err := config.LoadAWSConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	var wg sync.WaitGroup

	if config.GroupSpecific {
		// start a specific consumer for each group
		for groupNum := 1; groupNum <= config.NumGroups; groupNum++ {
			groupID := fmt.Sprintf("group-%d", groupNum)

			// for group-1, start configured number of consumers
			if groupNum == 1 {
				for i := 1; i <= config.Group1NumConsumers; i++ {
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
		for i := 1; i <= config.GeneralNumConsumers; i++ {
			wg.Add(1)
			go func(id int) {
				fmt.Printf("[Consumer] Starting general consumer %d (total consumers: %d)\n",
					id, config.GeneralNumConsumers)
				consumeAllMessages(ctx, sqsClient, id, &wg)
			}(i)
		}
	}

	// Wait for all consumers to finish or context to be canceled
	wg.Wait()
	fmt.Println("All consumers have shut down. Exiting.")
}

func consumeAllMessages(ctx context.Context, client *sqs.Client, consumerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Consumer %d] Shutting down due to context cancellation", consumerID)
			return
		default:
			resp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(config.QueueURL),
				MaxNumberOfMessages: 1,
				WaitTimeSeconds:     10, // Long polling
			})
			if err != nil {
				if ctx.Err() != nil {
					return // Context was canceled
				}
				log.Printf("[Consumer %d] Error receiving: %v", consumerID, err)
				time.Sleep(1 * time.Second) // Back off a bit before retrying
				continue
			}

			if len(resp.Messages) == 0 {
				continue // Keep polling for new messages
			}

			msg := resp.Messages[0]
			log.Printf("[Consumer %d] Processing: %s", consumerID, *msg.Body)

			// processing the message
			time.Sleep(2 * time.Second)

			_, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(config.QueueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Printf("[Consumer %d] Error deleting message: %v", consumerID, err)
			} else {
				log.Printf("[Consumer %d] Deleted: %s", consumerID, *msg.Body)
			}
		}
	}
}

func consumeOnlyGroup(ctx context.Context, client *sqs.Client, consumerID int, expectedGroup string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Consumer %d][%s] Shutting down due to context cancellation", consumerID, expectedGroup)
			return
		default:
			resp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(config.QueueURL),
				MaxNumberOfMessages: 1,
				WaitTimeSeconds:     10, // Long polling
			})
			if err != nil {
				if ctx.Err() != nil {
					return // Context was canceled
				}
				log.Printf("[Consumer %d][%s] Error receiving: %v", consumerID, expectedGroup, err)
				time.Sleep(1 * time.Second) // Back off a bit before retrying
				continue
			}

			if len(resp.Messages) == 0 {
				continue // Keep polling for new messages
			}

			msg := resp.Messages[0]

			// manual filtering, only processing if the message belongs to the expected group
			if !isGroupMessage(msg, expectedGroup) {
				log.Printf("[Consumer %d][%s] Ignoring message from another group: %s", consumerID, expectedGroup, *msg.Body)
				_, err = client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
					QueueUrl:          aws.String(config.QueueURL),
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
				QueueUrl:      aws.String(config.QueueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Printf("[Consumer %d][%s] Error deleting message: %v", consumerID, expectedGroup, err)
			} else {
				log.Printf("[Consumer %d][%s] Deleted: %s", consumerID, expectedGroup, *msg.Body)
			}
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
