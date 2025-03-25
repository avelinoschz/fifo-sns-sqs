# SQS-FIFO

A local FIFO SQS simulator using AWS SDK and LocalStack.

This project sends and consumes messages from a FIFO SQS queue with support for multiple message groups and concurrent consumers. Useful for testing message ordering and consumer behavior in a FIFO queue setup.

## Goal

Understand the parallelization capabilities of FIFO SQS.

## Architecture Update: SNS → SQS

This project uses an SNS FIFO topic as the entry point for messages. Instead of publishing directly to the SQS FIFO queue, the producer sends messages to the SNS FIFO topic. The topic is subscribed to the queue.

This simulates a more realistic architecture for distributed systems, where publishers communicate through a centralized messaging hub. FIFO behavior is still preserved thanks to the use of `MessageGroupId` and FIFO-compatible components.

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
