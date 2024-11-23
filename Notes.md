### NOTES.md

```markdown
# Notes and Learnings

## Producer Acknowledgment (`acks`) Levels
- **`acks=0`**:  
  The producer sends the message without waiting for acknowledgment.  
  - **Throughput**: Highest  
  - **Reliability**: Lowest (messages might be lost)

- **`acks=1`**:  
  The producer waits for the leader of the partition to acknowledge the message.  
  - **Reliability**: Ensures the message is stored on the leader but may not yet be replicated.

- **`acks=all`**:  
  The producer waits for all in-sync replicas to acknowledge the message.  
  - **Reliability**: Highest, ensures durability and fault tolerance.

---

## Creating a Kafka Topic
To create a topic, use the Kafka CLI or Kafka-UI. For CLI, follow the command below:
```bash
/bin/kafka-topics --create --bootstrap-server kafka1:19092,kafka2:19093,kafka3:19094 \
  --topic example-topic \
  --partitions 3 \
  --replication-factor 3 \
  --config min.insync.replicas=2
```

---

## Key Rules for Kafka Partitions and Consumers
1. **Partition Consumption Rule**:  
   Each partition in a Kafka topic can be consumed by only one consumer within a consumer group at any given time.

2. **Maximum Consumers per Consumer Group**:  
   The maximum number of consumers in a group equals the number of partitions in the topic(s) being consumed.

### Scenarios
- **Consumers ≤ Partitions**:  
  If consumers are fewer or equal to the partitions:
  - **Example**: Topic has 10 partitions, consumer group has 5 consumers.  
    - Each consumer is assigned 2 partitions.

- **Consumers > Partitions**:  
  If consumers exceed the number of partitions:
  - **Example**: Topic has 4 partitions, consumer group has 6 consumers.  
    - Kafka assigns 4 consumers to partitions, and the remaining 2 are idle.

- **Consumers = Partitions**:  
  If consumers equal partitions:
  - **Example**: Topic has 8 partitions, consumer group has 8 consumers.  
    - Each consumer is assigned 1 partition.

---

## Producer and Consumer Workflow in Kafka
1. **Producer Sends Message**:  
   - A producer sends a message to a Kafka topic, where it is assigned an offset and stored in one of the topic’s partitions.

2. **Consumer Reads Message**:  
   - A consumer subscribed to the topic fetches the message and processes it. Offsets can be committed automatically or manually once the message is processed.

3. **Scaling with Consumer Groups**:  
   - Multiple consumers in a group can read from the same topic. Kafka ensures each partition is processed by only one consumer in the group for parallel processing.

### Message Order Guarantee
- Kafka guarantees message ordering **within a partition**.
- To maintain order for related messages, ensure that messages with the same key are routed to the same partition.

---

## Kafka Message Delivery Semantics
Kafka provides three levels of delivery guarantees for messages exchanged between producers and consumers:

### 1. At-Most-Once Delivery
- **Guarantee**: Messages may be lost but will never be delivered more than once.
- **Performance**: Fastest, as there are no retries.
- **Use Cases**:  
  - Logging or non-critical monitoring where occasional data loss is acceptable.
- **How it Works**:  
  - The producer sends a message once without retrying.  
  - The consumer processes the message and immediately commits the offset.

### 2. At-Least-Once Delivery
- **Guarantee**: Every message will be delivered, but duplicates are possible in case of retries.
- **Default**: Kafka’s default behavior.
- **Use Cases**:  
  - Applications like payment processing, order systems, or microservices where no message loss is acceptable.
- **How it Works**:  
  - Producer retries sending messages if no acknowledgment is received.  
  - Consumers manually commit offsets after processing. Retries may cause duplicates.

---

### 3. Exactly-Once Delivery
- **Guarantee**: Each message is delivered exactly once.
- **How it Works**:
  - Enable **`enable.idempotence=true`** for the producer to ensure retries do not result in duplicate messages.
  - Use **`transactional.id`** for managing transactions, ensuring either full commit or rollback of all messages.

#### Configurations for Exactly-Once Delivery:
- **`enable.idempotence=true`**:  
  - Ensures a message is written only once, even with retries.  
  - **Default**: False, must be explicitly enabled.

- **`transactional.id`**:  
  - A unique identifier for the producer’s transactions.  
  - Guarantees transactional consistency by committing or aborting messages in bulk.

#### Example:
- **Transactional ID**: `txn-1`  
  - Each producer must use a unique ID.  
  - If the producer crashes, it resumes from its transactional state using this ID.

---

### Summary of Delivery Guarantees:
| **Delivery Guarantee** | **Behavior**                        | **Use Case**                                             |
|-------------------------|-------------------------------------|---------------------------------------------------------|
| **At-Most-Once**        | No retries; potential message loss | Low-priority data, fast performance                    |
| **At-Least-Once**       | Retries with potential duplicates  | Reliable systems where duplicates are acceptable       |
| **Exactly-Once**        | Guaranteed unique delivery         | Critical systems requiring no duplicates or data loss  |