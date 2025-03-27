package main

import (
	"fmt"
	"sort"
	"sync"
)

// MessageStats tracks the number of messages processed by each consumer/group
type MessageStats struct {
	mu     sync.Mutex
	counts map[int]map[string]int // map[consumerID]map[groupID]count

}

func NewMessageStats() *MessageStats {
	return &MessageStats{
		counts: make(map[int]map[string]int),
	}
}

// IncrementCount increments the message count for a specific consumer and group
func (ms *MessageStats) IncrementCount(consumerID int, groupID string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.counts[consumerID]; !exists {
		ms.counts[consumerID] = make(map[string]int)
	}
	ms.counts[consumerID][groupID]++
}

func (ms *MessageStats) PrintStats() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	fmt.Println("\nMessage Processing Statistics")
	fmt.Println("By Consumer/Group:")

	groupTotals := make(map[string]int)
	consumerTotals := make(map[int]int)
	totalMessages := 0

	// calculate totals
	for consumerID, groups := range ms.counts {
		consumerTotal := 0
		for groupID, count := range groups {
			groupTotals[groupID] += count
			consumerTotal += count
		}
		consumerTotals[consumerID] = consumerTotal
		totalMessages += consumerTotal
	}

	for consumerID, groups := range ms.counts {
		fmt.Printf("Consumer %d:\n", consumerID)
		for groupID, count := range groups {
			fmt.Printf("  - Group %s: %d messages\n", groupID, count)
		}
		fmt.Printf("  Total: %d messages\n", consumerTotals[consumerID])
	}

	fmt.Println("\nBy Group:")

	// sort groups for consistent output
	sortedGroups := make([]string, 0, len(groupTotals))
	for group := range groupTotals {
		sortedGroups = append(sortedGroups, group)
	}
	sort.Strings(sortedGroups)

	for _, groupID := range sortedGroups {
		count := groupTotals[groupID]
		percentage := 0.0
		if totalMessages > 0 {
			percentage = float64(count) / float64(totalMessages) * 100
		}
		fmt.Printf("  - Group %s: %d messages (%.1f%%)\n", groupID, count, percentage)
	}

	fmt.Println("\nProcessing Load Distribution:")
	fmt.Println("----------------------------------")

	// sort consumers for consistent output
	sortedConsumers := make([]int, 0, len(consumerTotals))
	for consumer := range consumerTotals {
		sortedConsumers = append(sortedConsumers, consumer)
	}
	sort.Ints(sortedConsumers)

	for _, consumerID := range sortedConsumers {
		count := consumerTotals[consumerID]
		percentage := 0.0
		if totalMessages > 0 {
			percentage = float64(count) / float64(totalMessages) * 100
		}
		fmt.Printf("Consumer %d: %d messages (%.1f%%)\n", consumerID, count, percentage)
	}

	fmt.Printf("\nTotal Messages Processed: %d\n", totalMessages)
	fmt.Println("-----------------------------------")
}
