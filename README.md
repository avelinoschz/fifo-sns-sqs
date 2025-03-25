# FIFO-SNS-SQS

A local FIFO SQS simulator using AWS SDK and LocalStack.

This project sends and consumes messages to a FIFO SNS and SQS with support for multiple message groups and concurrent consumers. Useful for testing message ordering and consumer behavior in a FIFO topic and queue setup.

## Architecture SNS → SQS

This project uses an SNS FIFO topic as the entry point for messages. Instead of publishing directly to the SQS FIFO queue, the producer sends messages to the SNS FIFO topic. The topic is subscribed to the queue.

This simulates a more realistic architecture for distributed systems, where publishers communicate through a centralized messaging hub. FIFO behavior is still preserved thanks to the use of `MessageGroupId` and FIFO-compatible components.

## Goal

Understand the parallelization capabilities of FIFO SQS.

## Characteristics of a FIFO topic & queue

FIFO guarantees message ordering, and has the feature to manage messages by MessageGroupId.

The MessageGroupId needs to be assigned per message, is not a configuration on the producer or consumer level.

## How to Run

1. Install [LocalStack](https://docs.localstack.cloud/get-started/)
2. Create the FIFO queue:

    ```zsh
        aws --endpoint-url=http://localhost:4566 --region=us-east-1 --no-sign-request --no-paginate \
        sqs create-queue \ 
            --queue-name demo-queue.fifo --attributes FifoQueue=true,ContentBasedDeduplication=true
    ```

3. Create the FIFO SNS topic:

    ```zsh
        aws --endpoint-url=http://localhost:4566 --region=us-east-1 --no-sign-request --no-paginate \
        sns create-topic \
            --name demo-topic.fifo --attributes FifoTopic=true,ContentBasedDeduplication=true
    ```

4. Subscribe the SQS queue to the SNS topic:

    ```zsh
        aws --endpoint-url=http://localhost:4566 --region=us-east-1 --no-sign-request --no-paginate \
        sns subscribe \
            --topic-arn arn:aws:sns:us-east-1:000000000000:demo-topic.fifo \
            --protocol sqs \
            --notification-endpoint arn:aws:sqs:us-east-1:000000000000:demo-queue.fifo \
            --attributes '{"RawMessageDelivery":"true"}'
    ```

5. Run the Go app with the producer and consumers:

    ```zsh
        go run main.go
    ```

## Test Scenarios

### Scenario 1: One Consumer Per Group

In this scenario, we use 3 different message group ids, with one consumer assigned to each group.

**Expected Behavior**: Messages should be processed in parallel, with strict ordering within each group.

**Result**:

```text
Starting FIFO Producer and Consumers...
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 1
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 2
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 3
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 4
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 5
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 6
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 7
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 8
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 9
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 10
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 1
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 2
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 3
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 4
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 5
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 6
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 7
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 8
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 9
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 10
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 1
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 2
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 3
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 4
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 5
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 6
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 7
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 8
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 9
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 10
[Consumer] Starting consumer 3 for group-3
[Consumer] Starting consumer 1 for group-1
[Consumer] Starting consumer 2 for group-2
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 1
[Consumer 2][group-2] Processing: group-2 - Message 1
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 1
[Consumer 1][group-1] Processing: group-1 - Message 1
[Consumer 3][group-3] Processing: group-3 - Message 1
[Consumer 3][group-3] Deleted: group-3 - Message 1
[Consumer 1][group-1] Deleted: group-1 - Message 1
[Consumer 2][group-2] Deleted: group-2 - Message 1
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 2
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 2
[Consumer 2][group-2] Processing: group-2 - Message 2
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 2
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 2
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 2
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 2
[Consumer 3][group-3] Processing: group-3 - Message 2
[Consumer 1][group-1] Processing: group-1 - Message 2
[Consumer 2][group-2] Deleted: group-2 - Message 2
[Consumer 2][group-2] Processing: group-2 - Message 3
[Consumer 3][group-3] Deleted: group-3 - Message 2
[Consumer 1][group-1] Deleted: group-1 - Message 2
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 3
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 3
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 3
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 3
[Consumer 3][group-3] Processing: group-3 - Message 3
[Consumer 1][group-1] Processing: group-1 - Message 3
[Consumer 2][group-2] Deleted: group-2 - Message 3
[Consumer 2][group-2] Processing: group-2 - Message 4
[Consumer 3][group-3] Deleted: group-3 - Message 3
[Consumer 1][group-1] Deleted: group-1 - Message 3
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 4
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 4
[Consumer 3][group-3] Processing: group-3 - Message 4
[Consumer 1][group-1] Processing: group-1 - Message 4
[Consumer 2][group-2] Deleted: group-2 - Message 4
[Consumer 2][group-2] Processing: group-2 - Message 5
[Consumer 3][group-3] Deleted: group-3 - Message 4
[Consumer 1][group-1] Deleted: group-1 - Message 4
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 5
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 5
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 5
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 5
[Consumer 3][group-3] Processing: group-3 - Message 5
[Consumer 1][group-1] Processing: group-1 - Message 5
[Consumer 2][group-2] Deleted: group-2 - Message 5
[Consumer 2][group-2] Processing: group-2 - Message 6
[Consumer 1][group-1] Deleted: group-1 - Message 5
[Consumer 3][group-3] Deleted: group-3 - Message 5
[Consumer 3][group-3] Processing: group-3 - Message 6
[Consumer 1][group-1] Processing: group-1 - Message 6
[Consumer 2][group-2] Deleted: group-2 - Message 6
[Consumer 2][group-2] Processing: group-2 - Message 7
[Consumer 1][group-1] Deleted: group-1 - Message 6
[Consumer 3][group-3] Deleted: group-3 - Message 6
[Consumer 1][group-1] Processing: group-1 - Message 7
[Consumer 3][group-3] Processing: group-3 - Message 7
[Consumer 2][group-2] Deleted: group-2 - Message 7
[Consumer 2][group-2] Processing: group-2 - Message 8
[Consumer 3][group-3] Deleted: group-3 - Message 7
[Consumer 1][group-1] Deleted: group-1 - Message 7
[Consumer 1][group-1] Processing: group-1 - Message 8
[Consumer 3][group-3] Processing: group-3 - Message 8
[Consumer 2][group-2] Deleted: group-2 - Message 8
[Consumer 2][group-2] Processing: group-2 - Message 9
[Consumer 3][group-3] Deleted: group-3 - Message 8
[Consumer 1][group-1] Deleted: group-1 - Message 8
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 9
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 9
[Consumer 1][group-1] Processing: group-1 - Message 9
[Consumer 3][group-3] Processing: group-3 - Message 9
[Consumer 2][group-2] Deleted: group-2 - Message 9
[Consumer 2][group-2] Processing: group-2 - Message 10
[Consumer 3][group-3] Deleted: group-3 - Message 9
[Consumer 1][group-1] Deleted: group-1 - Message 9
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 10
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 10
[Consumer 1][group-1] Processing: group-1 - Message 10
[Consumer 3][group-3] Processing: group-3 - Message 10
[Consumer 2][group-2] Deleted: group-2 - Message 10
[Consumer 3][group-3] Deleted: group-3 - Message 10
[Consumer 1][group-1] Deleted: group-1 - Message 10
[Consumer 2][group-2] Finished processing all messages in 22.194683375s
[Consumer 1][group-1] Finished processing all messages in 22.270463333s
[Consumer 3][group-3] Finished processing all messages in 22.270532792s
```

As we can see, the 3 consumers are processing the messages in parallel and in expected order.
One important detail is that `MessageGroupId` is assigned at the **message level**, not the consumer or queue level. This means any consumer can receive any message, regardless of its group.

```text
2025/03/21 19:03:31 [Consumer 2][group-2] Ignoring message from another group: group-3 - Message 10
2025/03/21 19:03:31 [Consumer 3][group-3] Ignoring message from another group: group-2 - Message 10
```

This happens because SQS FIFO does not route messages based on `MessageGroupId`. It is up to the application to implement group-based filtering after receiving messages.

### Scenario 2: Multiple Consumers for One Group

In this scenario, we instantiate 3 consumers for message group id 1, and one consumer each for groups 2 and 3.

**Expected Behavior**: Messages from group-1 are processed sequentially by competing consumers, but only one at a time due to FIFO constraints. Groups 2 and 3 continue to process in parallel.

**Result**:

```text
Starting FIFO Producer and Consumers...
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 1
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 2
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 3
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 4
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 5
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 6
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 7
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 8
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 9
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 10
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 1
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 2
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 3
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 4
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 5
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 6
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 7
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 8
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 9
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 10
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 1
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 2
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 3
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 4
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 5
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 6
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 7
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 8
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 9
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 10
[Consumer] Starting consumer 3 for group-3
[Consumer] Starting consumer 2 for group-1
[Consumer] Starting consumer 3 for group-1
[Consumer] Starting consumer 2 for group-2
[Consumer] Starting consumer 1 for group-1
[Consumer 3][group-3] Ignoring message from another group: group-1 - Message 1
[Consumer 2][group-1] Ignoring message from another group: group-3 - Message 1
[Consumer 2][group-2] Processing: group-2 - Message 1
[Consumer 1][group-1] Processing: group-1 - Message 1
[Consumer 3][group-1] Ignoring message from another group: group-3 - Message 1
[Consumer 3][group-3] Processing: group-3 - Message 1
[Consumer 2][group-2] Deleted: group-2 - Message 1
[Consumer 2][group-1] Ignoring message from another group: group-2 - Message 2
[Consumer 1][group-1] Deleted: group-1 - Message 1
[Consumer 2][group-2] Ignoring message from another group: group-3 - Message 2
[Consumer 3][group-1] Processing: group-1 - Message 2
[Consumer 3][group-3] Deleted: group-3 - Message 1
[Consumer 3][group-3] Processing: group-3 - Message 2
[Consumer 1][group-1] Ignoring message from another group: group-2 - Message 2
[Consumer 2][group-1] Ignoring message from another group: group-2 - Message 2
[Consumer 2][group-2] Processing: group-2 - Message 2
[Consumer 3][group-1] Deleted: group-1 - Message 2
[Consumer 1][group-1] Processing: group-1 - Message 3
[Consumer 3][group-3] Deleted: group-3 - Message 2
[Consumer 2][group-1] Ignoring message from another group: group-3 - Message 3
[Consumer 3][group-1] Ignoring message from another group: group-2 - Message 3
[Consumer 2][group-2] Deleted: group-2 - Message 2
[Consumer 3][group-3] Processing: group-3 - Message 3
[Consumer 2][group-2] Processing: group-2 - Message 3
[Consumer 1][group-1] Deleted: group-1 - Message 3
[Consumer 2][group-1] Processing: group-1 - Message 4
[Consumer 3][group-3] Deleted: group-3 - Message 3
[Consumer 2][group-2] Deleted: group-2 - Message 3
[Consumer 1][group-1] Ignoring message from another group: group-2 - Message 4
[Consumer 3][group-1] Ignoring message from another group: group-3 - Message 4
[Consumer 3][group-3] Processing: group-3 - Message 4
[Consumer 2][group-2] Processing: group-2 - Message 4
[Consumer 2][group-1] Deleted: group-1 - Message 4
[Consumer 3][group-1] Processing: group-1 - Message 5
[Consumer 3][group-3] Deleted: group-3 - Message 4
[Consumer 2][group-2] Deleted: group-2 - Message 4
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 5
[Consumer 2][group-1] Ignoring message from another group: group-2 - Message 5
[Consumer 2][group-2] Ignoring message from another group: group-3 - Message 5
[Consumer 3][group-3] Ignoring message from another group: group-2 - Message 5
[Consumer 2][group-1] Ignoring message from another group: group-2 - Message 5
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 5
[Consumer 3][group-3] Ignoring message from another group: group-2 - Message 5
[Consumer 2][group-2] Ignoring message from another group: group-3 - Message 5
[Consumer 2][group-1] Ignoring message from another group: group-3 - Message 5
[Consumer 1][group-1] Ignoring message from another group: group-2 - Message 5
[Consumer 2][group-2] Ignoring message from another group: group-3 - Message 5
[Consumer 3][group-3] Ignoring message from another group: group-2 - Message 5
[Consumer 1][group-1] Ignoring message from another group: group-2 - Message 5
[Consumer 2][group-1] Ignoring message from another group: group-3 - Message 5
[Consumer 2][group-2] Processing: group-2 - Message 5
[Consumer 3][group-3] Processing: group-3 - Message 5
[Consumer 3][group-1] Deleted: group-1 - Message 5
[Consumer 1][group-1] Processing: group-1 - Message 6
[Consumer 2][group-1] Finished processing all messages in 10.14486175s
[Consumer 2][group-2] Deleted: group-2 - Message 5
[Consumer 3][group-1] Ignoring message from another group: group-2 - Message 6
[Consumer 3][group-3] Deleted: group-3 - Message 5
[Consumer 2][group-2] Ignoring message from another group: group-3 - Message 6
[Consumer 3][group-1] Ignoring message from another group: group-2 - Message 6
[Consumer 3][group-3] Processing: group-3 - Message 6
[Consumer 2][group-2] Processing: group-2 - Message 6
[Consumer 1][group-1] Deleted: group-1 - Message 6
[Consumer 3][group-1] Processing: group-1 - Message 7
[Consumer 3][group-3] Deleted: group-3 - Message 6
[Consumer 1][group-1] Ignoring message from another group: group-3 - Message 7
[Consumer 2][group-2] Deleted: group-2 - Message 6
[Consumer 3][group-3] Processing: group-3 - Message 7
[Consumer 1][group-1] Ignoring message from another group: group-2 - Message 7
[Consumer 2][group-2] Processing: group-2 - Message 7
[Consumer 3][group-1] Deleted: group-1 - Message 7
[Consumer 1][group-1] Processing: group-1 - Message 8
[Consumer 2][group-2] Deleted: group-2 - Message 7
[Consumer 3][group-3] Deleted: group-3 - Message 7
[Consumer 2][group-2] Ignoring message from another group: group-3 - Message 8
[Consumer 3][group-1] Ignoring message from another group: group-2 - Message 8
[Consumer 3][group-3] Processing: group-3 - Message 8
[Consumer 2][group-2] Processing: group-2 - Message 8
[Consumer 1][group-1] Deleted: group-1 - Message 8
[Consumer 3][group-1] Processing: group-1 - Message 9
[Consumer 2][group-2] Deleted: group-2 - Message 8
[Consumer 1][group-1] Ignoring message from another group: group-2 - Message 9
[Consumer 3][group-3] Deleted: group-3 - Message 8
[Consumer 2][group-2] Ignoring message from another group: group-3 - Message 9
[Consumer 1][group-1] Ignoring message from another group: group-2 - Message 9
[Consumer 3][group-3] Processing: group-3 - Message 9
[Consumer 2][group-2] Processing: group-2 - Message 9
[Consumer 3][group-1] Deleted: group-1 - Message 9
[Consumer 1][group-1] Processing: group-1 - Message 10
[Consumer 2][group-2] Deleted: group-2 - Message 9
[Consumer 3][group-3] Deleted: group-3 - Message 9
[Consumer 3][group-1] Ignoring message from another group: group-2 - Message 10
[Consumer 3][group-3] Processing: group-3 - Message 10
[Consumer 2][group-2] Processing: group-2 - Message 10
[Consumer 1][group-1] Deleted: group-1 - Message 10
[Consumer 3][group-3] Deleted: group-3 - Message 10
[Consumer 3][group-1] Finished processing all messages in 20.255861792s
[Consumer 2][group-2] Deleted: group-2 - Message 10
[Consumer 1][group-1] Finished processing all messages in 22.150736083s
[Consumer 3][group-3] Finished processing all messages in 22.261156667s
[Consumer 2][group-2] Finished processing all messages in 22.261213083s
```

At the middle of the run, we see the following log line:

``
[Consumer 2][group-1] Finished processing all messages in 10.14486175s
``

This happened due to starvation. FIFO queue guarantees order, so the messages get locked until processed. Having 3 consumers for same group causes them to race for the message.

## Conclusions

- SQS FIFO queues guarantee strict ordering within a MessageGroupId.
- Only one message per group can be processed at a time, even with multiple consumers.
- Parallel processing is only achieved when multiple distinct MessageGroupIds are used.
- Extra consumers for the same group do not improve throughput.
- Consumers may receive messages for the wrong group and must implement filtering logic.

**Note**: In this test, messages used a naming convention to identify their group. In production, group info should be extracted from message metadata or payload, depending on how it's published.

## Updated version: Scenario 3: One Consumer Per Group going through actual SNS + SQS AWS services

In this scenario, we instantiate one consumer for groups each group and are connecting directly to an actual AWS instance.

**Expected Behavior**: Messages from group-1 are processed sequentially by competing consumers, but only one at a time due to FIFO constraints. Groups 2 and 3 continue to process in parallel.

**Result**:

```text
Starting FIFO Producer and Consumers...
Connecting to AWS using profile
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 1
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 2
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 3
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 4
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 5
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 6
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 7
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 8
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 9
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-1 Body: group-1 - Message 10
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 1
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 2
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 3
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 4
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 5
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 6
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 7
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 8
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 9
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-2 Body: group-2 - Message 10
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 1
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 2
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 3
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 4
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 5
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 6
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 7
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 8
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 9
[Producer] Sending message to SNS topic: arn:aws:sns:us-east-1:000000000000:demo-topic.fifo GroupId: group-3 Body: group-3 - Message 10
[Consumer] Starting consumer 1 for group-1
[Consumer] Starting consumer 3 for group-3
[Consumer] Starting consumer 2 for group-2
[Consumer 1][group-1] Processing: {
  "Type" : "Notification",
  "MessageId" : "a811567b-7f70-5b33-b42e-36275bea37f2",
  "SequenceNumber" : "10000000000000003000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 1",
  "Timestamp" : "2025-03-25T03:26:18.422Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Processing: {
  "Type" : "Notification",
  "MessageId" : "2eebe0d1-5805-50d6-8978-808d7c9ebc17",
  "SequenceNumber" : "10000000000000013000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 1",
  "Timestamp" : "2025-03-25T03:26:19.402Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Processing: {
  "Type" : "Notification",
  "MessageId" : "b511c497-9559-54ca-ae16-a4cbb1c5b80b",
  "SequenceNumber" : "10000000000000023000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 1",
  "Timestamp" : "2025-03-25T03:26:20.237Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Deleted: {
  "Type" : "Notification",
  "MessageId" : "b511c497-9559-54ca-ae16-a4cbb1c5b80b",
  "SequenceNumber" : "10000000000000023000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 1",
  "Timestamp" : "2025-03-25T03:26:20.237Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Deleted: {
  "Type" : "Notification",
  "MessageId" : "2eebe0d1-5805-50d6-8978-808d7c9ebc17",
  "SequenceNumber" : "10000000000000013000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 1",
  "Timestamp" : "2025-03-25T03:26:19.402Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Deleted: {
  "Type" : "Notification",
  "MessageId" : "a811567b-7f70-5b33-b42e-36275bea37f2",
  "SequenceNumber" : "10000000000000003000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 1",
  "Timestamp" : "2025-03-25T03:26:18.422Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "fe5af9c4-f6bf-5cf0-8c75-9998d225a031",
  "SequenceNumber" : "10000000000000004000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 2",
  "Timestamp" : "2025-03-25T03:26:18.639Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "287a640d-c3b8-5e11-956e-0451e72492d0",
  "SequenceNumber" : "10000000000000024000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 2",
  "Timestamp" : "2025-03-25T03:26:20.320Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "0bbca7ef-6f6a-5638-a211-c7b34bf16af9",
  "SequenceNumber" : "10000000000000014000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 2",
  "Timestamp" : "2025-03-25T03:26:19.485Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "8cc3ef25-75da-52d2-9400-8f7951404e13",
  "SequenceNumber" : "10000000000000005000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 3",
  "Timestamp" : "2025-03-25T03:26:18.725Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "6afd3ba5-6a85-5711-bb11-5bf94e1e7bd0",
  "SequenceNumber" : "10000000000000015000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 3",
  "Timestamp" : "2025-03-25T03:26:19.566Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "9fa19103-3012-5d90-9b14-f15f148ef151",
  "SequenceNumber" : "10000000000000025000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 3",
  "Timestamp" : "2025-03-25T03:26:20.403Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "9a7b3680-eae4-5922-b1fb-56dcce7e15e3",
  "SequenceNumber" : "10000000000000006000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 4",
  "Timestamp" : "2025-03-25T03:26:18.808Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "7f7e7446-c3da-5a73-b696-141ae746f6b2",
  "SequenceNumber" : "10000000000000026000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 4",
  "Timestamp" : "2025-03-25T03:26:20.486Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "380202d6-4e9b-538a-9fc1-628a6a002f4b",
  "SequenceNumber" : "10000000000000016000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 4",
  "Timestamp" : "2025-03-25T03:26:19.651Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Processing: {
  "Type" : "Notification",
  "MessageId" : "9ddee4fe-4542-5b2b-a52d-84a261d857ea",
  "SequenceNumber" : "10000000000000007000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 5",
  "Timestamp" : "2025-03-25T03:26:18.890Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "55159054-e096-5681-9c4b-f008c40aea40",
  "SequenceNumber" : "10000000000000027000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 5",
  "Timestamp" : "2025-03-25T03:26:20.570Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "962eccdc-a94a-583e-9ef5-ac808dd174ff",
  "SequenceNumber" : "10000000000000017000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 5",
  "Timestamp" : "2025-03-25T03:26:19.736Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "fb1fa075-1cd5-5765-9bdb-c0293849fba4",
  "SequenceNumber" : "10000000000000028000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 6",
  "Timestamp" : "2025-03-25T03:26:20.653Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "a5b91548-09cf-55d5-82e7-75142d541914",
  "SequenceNumber" : "10000000000000018000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 6",
  "Timestamp" : "2025-03-25T03:26:19.821Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "cf230055-b5a9-588a-9560-60099135c44b",
  "SequenceNumber" : "10000000000000029000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 7",
  "Timestamp" : "2025-03-25T03:26:20.737Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "b5b3f444-64ec-54ec-a11b-94184077b82d",
  "SequenceNumber" : "10000000000000019000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 7",
  "Timestamp" : "2025-03-25T03:26:19.904Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "67a05e37-140f-5ea7-ad33-3e6b3093e266",
  "SequenceNumber" : "10000000000000030000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 8",
  "Timestamp" : "2025-03-25T03:26:20.820Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "bd0c6de4-2f91-57e8-80a0-0db80b0aab29",
  "SequenceNumber" : "10000000000000020000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 8",
  "Timestamp" : "2025-03-25T03:26:19.986Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "a043e5c5-104d-561e-9d3d-859934d14adf",
  "SequenceNumber" : "10000000000000031000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 9",
  "Timestamp" : "2025-03-25T03:26:20.904Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "653c3437-a451-5674-b005-546b80cddab5",
  "SequenceNumber" : "10000000000000021000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 9",
  "Timestamp" : "2025-03-25T03:26:20.070Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "0ec697b3-ebfd-5cf1-8056-f7efcf61f030",
  "SequenceNumber" : "10000000000000032000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-3 - Message 10",
  "Timestamp" : "2025-03-25T03:26:20.988Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "97107e1f-99a6-5cce-b28c-1e43bd3af00f",
  "SequenceNumber" : "10000000000000022000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-2 - Message 10",
  "Timestamp" : "2025-03-25T03:26:20.153Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Deleted: {
  "Type" : "Notification",
  "MessageId" : "9ddee4fe-4542-5b2b-a52d-84a261d857ea",
  "SequenceNumber" : "10000000000000007000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 5",
  "Timestamp" : "2025-03-25T03:26:18.890Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "cc685a2c-6ed7-58cc-b8c3-3a704a4db40d",
  "SequenceNumber" : "10000000000000008000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 6",
  "Timestamp" : "2025-03-25T03:26:18.978Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Ignoring message from another group: {
  "Type" : "Notification",
  "MessageId" : "c9239fed-966d-5eab-9591-ffb1b2f6426a",
  "SequenceNumber" : "10000000000000009000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 7",
  "Timestamp" : "2025-03-25T03:26:19.065Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Processing: {
  "Type" : "Notification",
  "MessageId" : "347fc72b-18ce-5061-9925-a7773c554ca0",
  "SequenceNumber" : "10000000000000010000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 8",
  "Timestamp" : "2025-03-25T03:26:19.150Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 2][group-2] Finished processing all messages in 7.21575425s
[Consumer 1][group-1] Deleted: {
  "Type" : "Notification",
  "MessageId" : "347fc72b-18ce-5061-9925-a7773c554ca0",
  "SequenceNumber" : "10000000000000010000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 8",
  "Timestamp" : "2025-03-25T03:26:19.150Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 3][group-3] Finished processing all messages in 7.27819425s
[Consumer 1][group-1] Processing: {
  "Type" : "Notification",
  "MessageId" : "61015d69-6077-50d5-8bb2-3d3544ec0c14",
  "SequenceNumber" : "10000000000000011000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 9",
  "Timestamp" : "2025-03-25T03:26:19.232Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Deleted: {
  "Type" : "Notification",
  "MessageId" : "61015d69-6077-50d5-8bb2-3d3544ec0c14",
  "SequenceNumber" : "10000000000000011000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 9",
  "Timestamp" : "2025-03-25T03:26:19.232Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Processing: {
  "Type" : "Notification",
  "MessageId" : "b23fa870-c3fe-563c-8fad-c82ce765dafd",
  "SequenceNumber" : "10000000000000012000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 10",
  "Timestamp" : "2025-03-25T03:26:19.317Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Deleted: {
  "Type" : "Notification",
  "MessageId" : "b23fa870-c3fe-563c-8fad-c82ce765dafd",
  "SequenceNumber" : "10000000000000012000",
  "TopicArn" : "arn:aws:sns:us-east-1:000000000000:demo-topic.fifo",
  "Message" : "group-1 - Message 10",
  "Timestamp" : "2025-03-25T03:26:19.317Z",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:demo-topic.fifo:e7677be2-cd0c-4768-b831-3fd213ba0680"
}
[Consumer 1][group-1] Finished processing all messages in 13.719142125s
```

**Note**: In this last test, an actual AWS instance, using actual SNS and SQS services. For security reasons, all the arn values have been replaced to be similar as Localstack. The actual values are not shown in the logs.

