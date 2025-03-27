.PHONY: help run-producer run-consumer run-ensure-conn setup-localstack check-docker

help:
	@echo "Available commands:"
	@echo "  make run-producer     - Run the FIFO producer"
	@echo "  make run-consumer     - Run the FIFO consumer"
	@echo "  make run-ensure-conn  - List AWS topics and queues"
	@echo "  make setup-localstack - Setup Localstack resources (FIFO queue, topic, and subscription)"
	@echo "  make check-docker     - Verify Docker and Localstack container are running"

run-producer:
	@echo "Running FIFO producer..."
	go run cmd/producer/main.go

run-consumer:
	@echo "Running FIFO consumer..."
	go run cmd/consumer/*.go

run-ensure-conn:
	@echo "Checking AWS connections..."
	go run cmd/ensure-conn/list-topics-queues.go

check-docker:
	@echo "Checking Docker status..."
	@if ! docker info > /dev/null 2>&1; then \
		echo "❌ Docker is not running. Please start Docker first!"; \
		exit 1; \
	else \
		echo "✅ Docker is running"; \
	fi
	
	@echo "Checking Localstack container status..."
	@if ! docker ps | grep localstack > /dev/null; then \
		echo "❌ Localstack container is not running."; \
		echo "   Start with: docker run --rm -it -p 4566:4566 -p 4510-4559:4510-4559 localstack/localstack"; \
		exit 1; \
	else \
		echo "✅ Localstack container is running"; \
		echo "   Container info:"; \
		docker ps --filter name=localstack --format "   ID: {{.ID}}\n   Name: {{.Names}}\n   Status: {{.Status}}\n   Ports: {{.Ports}}"; \
	fi
	@echo "Docker environment is ready!"

setup-localstack: check-docker
	@echo "Setting up Localstack resources..."
	@echo "1. Creating FIFO queue 'demo-queue.fifo'..."
	aws --endpoint-url=http://localhost:4566 --region=us-east-1 --no-sign-request --no-paginate \
		sqs create-queue \
		--queue-name demo-queue.fifo \
		--attributes FifoQueue=true,ContentBasedDeduplication=true
	
	@echo "2. Creating FIFO topic 'demo-topic.fifo'..."
	aws --endpoint-url=http://localhost:4566 --region=us-east-1 --no-sign-request --no-paginate \
		sns create-topic \
		--name demo-topic.fifo \
		--attributes FifoTopic=true,ContentBasedDeduplication=true
	
	@echo "3. Subscribing queue to topic..."
	aws --endpoint-url=http://localhost:4566 --region=us-east-1 --no-sign-request --no-paginate \
		sns subscribe \
		--topic-arn arn:aws:sns:us-east-1:000000000000:demo-topic.fifo \
		--protocol sqs \
		--notification-endpoint arn:aws:sqs:us-east-1:000000000000:demo-queue.fifo \
		--attributes '{"RawMessageDelivery":"true"}'
	
	@echo "✅ Localstack setup complete! You can verify resources with 'make run-ensure-conn'"