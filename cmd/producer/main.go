package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sort"
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

	rand.Seed(time.Now().UnixNano())

	groups := config.GenerateGroups()

	// count messages per group
	groupCounters := make(map[string]int)
	for _, group := range groups {
		groupCounters[group] = 0
	}

	fmt.Printf("Preparing to send %d messages randomly distributed across %d groups...\n",
		config.TotalMessages, config.NumGroups)

	startTime := time.Now()

	for i := 0; i < config.TotalMessages; i++ {
		// randomly select a group
		randomGroupIndex := rand.Intn(len(groups))
		groupID := groups[randomGroupIndex]

		groupCounters[groupID]++

		sendToSNS(ctx, snsClient, fmt.Sprintf("Message %d - Group %s", i+1, groupID), groupID)

		// small delay to avoid overwhelming the SNS service
		time.Sleep(10 * time.Millisecond)
	}

	elapsedTime := time.Since(startTime)

	fmt.Println("All messages have been sent successfully!")
	fmt.Printf("Time elapsed: %s\n", elapsedTime)

	fmt.Println("\nMessage distribution across groups:")
	fmt.Println("----------------------------------")

	// sort groups for consistent output
	sortedGroups := make([]string, 0, len(groupCounters))
	for group := range groupCounters {
		sortedGroups = append(sortedGroups, group)
	}
	sort.Strings(sortedGroups)

	totalSent := 0
	for _, group := range sortedGroups {
		count := groupCounters[group]
		totalSent += count
		percentage := float64(count) / float64(config.TotalMessages) * 100
		fmt.Printf("%s: %d messages (%.1f%%)\n", group, count, percentage)
	}

	fmt.Printf("\nTotal messages sent: %d\n", totalSent)
	fmt.Printf("Total time: %s\n", elapsedTime)
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
