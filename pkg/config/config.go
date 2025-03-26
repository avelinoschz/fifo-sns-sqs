package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/joho/godotenv"
)

var (
	// Configuration settings
	UseLocalstack bool
	GroupSpecific bool

	NumGroups        int
	MessagesPerGroup int

	// Consumer configurations
	Group1NumConsumers  int // consumers for group-1 in group-specific mode
	GeneralNumConsumers int // total consumers when not group-specific

	// AWS-related settings
	QueueURL string
	TopicARN string
	Region   string
	Profile  string
)

// Initialize configuration with default values and environment overrides
func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Set defaults then override with environment variables if available
	UseLocalstack = getEnvBool("USE_LOCALSTACK", true)
	GroupSpecific = getEnvBool("GROUP_SPECIFIC", false)

	NumGroups = getEnvInt("NUM_GROUPS", 6)
	MessagesPerGroup = getEnvInt("MESSAGES_PER_GROUP", 10)

	Group1NumConsumers = getEnvInt("GROUP1_NUM_CONSUMERS", 1)
	GeneralNumConsumers = getEnvInt("GENERAL_NUM_CONSUMERS", 3)

	Region = getEnvString("AWS_REGION", "us-east-1")

	if UseLocalstack {
		QueueURL = getEnvString("LOCALSTACK_QUEUE_URL", "http://localhost:4566/000000000000/demo-queue.fifo")
		TopicARN = getEnvString("LOCALSTACK_TOPIC_ARN", "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo")
	} else {
		QueueURL = os.Getenv("AWS_QUEUE_URL")
		TopicARN = os.Getenv("AWS_TOPIC_ARN")
		Profile = os.Getenv("AWS_PROFILE")
	}
}

// Helper functions for environment variable handling
func getEnvString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, fallback)
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "1" || value == "yes"
	}
	return fallback
}

// LoadAWSConfig loads AWS configuration based on current settings
func LoadAWSConfig(ctx context.Context) (aws.Config, error) {
	var cfg aws.Config
	var err error

	if UseLocalstack {
		fmt.Println("Using Localstack for development")
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(Region),
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
			config.WithSharedConfigProfile(Profile),
		)
	}

	return cfg, err
}

// GenerateGroups creates a slice of group IDs
func GenerateGroups() []string {
	groups := make([]string, NumGroups)
	for i := 0; i < NumGroups; i++ {
		groups[i] = fmt.Sprintf("group-%d", i+1)
	}
	return groups
}
