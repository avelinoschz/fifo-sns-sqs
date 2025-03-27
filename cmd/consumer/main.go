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
	startTime := time.Now() // Record start time

	// camcelable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// message stats tracker
	stats := NewMessageStats()

	// graceful shutdown
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
	totalConsumers := 0

	if config.GroupSpecific {
		// number of consumers for group-specific mode: group-1 config + (num_groups - 1)
		totalConsumers = config.Group1NumConsumers + (config.NumGroups - 1)
		fmt.Printf("[Consumer] Starting %d group-specific consumers...\n", totalConsumers)

		// start a specific consumer for each group
		for groupNum := 1; groupNum <= config.NumGroups; groupNum++ {
			groupID := fmt.Sprintf("group-%d", groupNum)

			// for group-1, start configured number of consumers
			if groupNum == 1 {
				for i := 1; i <= config.Group1NumConsumers; i++ {
					wg.Add(1)
					go func(id int, gID string) {
						fmt.Printf("[Consumer] Starting consumer %d for %s\n", id, gID)
						consumeOnlyGroup(ctx, sqsClient, id, gID, &wg, stats)
					}(i, groupID)
				}
				continue
			}

			// for other groups, start one consumer each
			wg.Add(1)
			go func(id int, gID string) {
				fmt.Printf("[Consumer] Starting consumer %d for %s\n", id, gID)
				consumeOnlyGroup(ctx, sqsClient, id, gID, &wg, stats)
			}(groupNum, groupID)
		}
	} else {
		totalConsumers = config.GeneralNumConsumers
		fmt.Printf("[Consumer] Starting %d general consumers...\n", totalConsumers)

		// start the general consumers
		for i := 1; i <= config.GeneralNumConsumers; i++ {
			wg.Add(1)
			go func(id int) {
				fmt.Printf("[Consumer] Starting general consumer %d\n", id)
				consumeAllMessages(ctx, sqsClient, id, &wg, stats)
			}(i)
		}
	}

	wg.Wait()

	stats.PrintStats()

	uptime := time.Since(startTime).Round(time.Second)
	fmt.Printf("Service was up for: %s\n", uptime)

	fmt.Println("All consumers have shut down. Exiting.")
}

func consumeAllMessages(ctx context.Context, client *sqs.Client, consumerID int, wg *sync.WaitGroup, stats *MessageStats) {
	defer wg.Done()

	consumeMessages(ctx, client, consumerID, "", stats, func(msg types.Message, groupID string) bool {
		// process all messages normally
		return true
	})
}

func consumeOnlyGroup(ctx context.Context, client *sqs.Client, consumerID int, expectedGroup string, wg *sync.WaitGroup, stats *MessageStats) {
	defer wg.Done()

	consumeMessages(ctx, client, consumerID, expectedGroup, stats, func(msg types.Message, groupID string) bool {
		// only process messages from the expected group
		if !isGroupMessage(msg, expectedGroup) {
			log.Printf("[Consumer %d][%s] Ignoring message from another group: %s", consumerID, expectedGroup, *msg.Body)
			return false
		}
		return true
	})
}

func consumeMessages(
	ctx context.Context,
	client *sqs.Client,
	consumerID int,
	groupLabel string, // empty for general consumers, group name for specific consumers
	stats *MessageStats,
	shouldProcess func(types.Message, string) bool, // Function to determine if we should process a message
) {
	logPrefix := fmt.Sprintf("[Consumer %d]", consumerID)
	if groupLabel != "" {
		logPrefix = fmt.Sprintf("[Consumer %d][%s]", consumerID, groupLabel)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("%s Shutting down due to context cancellation", logPrefix)
			return
		default:
			resp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(config.QueueURL),
				MaxNumberOfMessages: 1,
				WaitTimeSeconds:     10, // long polling
			})
			if err != nil {
				if ctx.Err() != nil {
					return // Context was canceled
				}
				log.Printf("%s Error receiving: %v", logPrefix, err)
				time.Sleep(1 * time.Second) // back off a bit before retrying
				continue
			}

			if len(resp.Messages) == 0 {
				continue // keep polling for new messages
			}

			msg := resp.Messages[0]

			// extract group ID from the message body for stats
			groupID := extractGroupID(*msg.Body)

			// depending on the mode (group-specific or general), decide if we should process the message
			if !shouldProcess(msg, groupID) {
				// Return the message to the queue
				_, err = client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
					QueueUrl:          aws.String(config.QueueURL),
					ReceiptHandle:     msg.ReceiptHandle,
					VisibilityTimeout: int32(0), // makes it visible again immediately
				})
				if err != nil {
					log.Printf("%s Error changing visibility: %v", logPrefix, err)
				}
				continue
			}

			log.Printf("%s Processing: %s", logPrefix, *msg.Body)

			// processing the message
			time.Sleep(200 * time.Millisecond)

			_, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(config.QueueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Printf("%s Error deleting message: %v", logPrefix, err)
			} else {
				stats.IncrementCount(consumerID, groupID)
				log.Printf("%s Deleted: %s", logPrefix, *msg.Body)
			}
		}
	}
}

// extractGroupID extracts the group id from the message body based on format
func extractGroupID(body string) string {
	for i := 1; i <= config.NumGroups; i++ {
		groupID := fmt.Sprintf("group-%d", i)
		if strings.Contains(body, groupID) {
			return groupID
		}
	}
	return "unknown"
}

// simple check based on group name in message body (group id is not returned by ReceiveMessage)
// supposedly, it should be on the msg metadata map, but didn't find it. Maybe due to LocalStack limitations.
func isGroupMessage(msg types.Message, group string) bool {
	return msg.Body != nil && len(*msg.Body) > 0 && contains(*msg.Body, group)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
