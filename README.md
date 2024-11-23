### README.md

```markdown
# Kafka Docker Setup and Topic Management

This guide provides step-by-step instructions for setting up and managing Kafka topics in a Docker environment using the provided `docker-compose.yml` file.

---

## Prerequisites
- Docker installed on your system.
- `docker-compose.yml` file configured as shown in this repository.

---

## Getting Started

1. **Start the Docker Compose Setup**  
   Run the following command to start the Kafka and Zookeeper services:
   ```bash
   docker-compose up -d
   ```
   This will start all the services defined in the `docker-compose.yml` file.

2. **Verify Services**  
   Ensure all services are running:
   ```bash
   docker ps
   ```
   You should see containers for Kafka brokers, Zookeeper instances, and Kafka-UI.

3. **Access Kafka UI**  
   Open your browser and navigate to [http://localhost:8080](http://localhost:8080) to access the Kafka UI.

---

## Managing Kafka Topics via CLI

### Step 1: SSH into a Kafka Container
To manage topics using the CLI, access one of the Kafka containers. For example:
```bash
docker exec -it kafka1 /bin/bash
```

### Step 2: Create a Kafka Topic
Inside the Kafka container, use the following command to create a new topic:
```bash
/bin/kafka-topics --create --bootstrap-server kafka1:19092,kafka2:19093,kafka3:19094 \
  --topic atleast-once-semantic \
  --partitions 3 \
  --replication-factor 3 \
  --config min.insync.replicas=2
```

### Step 3: Verify Topic Creation
To ensure the topic was created, list all topics:
```bash
/bin/kafka-topics --list --bootstrap-server kafka1:19092,kafka2:19093,kafka3:19094
```

### Step 4: Describe a Topic
To describe the newly created topic, use:
```bash
/bin/kafka-topics --bootstrap-server kafka1:19092,kafka2:19093,kafka3:19094 \
  --describe \
  --topic atleast-once-semantic
```

### Step 5: Exit the Kafka Container
After completing your tasks, exit the container:
```bash
exit
```

---

## Managing Kafka Topics via UI

1. **Access the Kafka UI**  
   Open [http://localhost:8080](http://localhost:8080) in your browser.

2. **Create a Topic**  
   - Navigate to the "Topics" section.
   - Click on "Create Topic."
   - Fill in the details:
     - Topic Name: `atleast-once-semantic`
     - Partitions: `3`
     - Replication Factor: `3`
   - Add configuration `min.insync.replicas=2` if necessary.

3. **Describe a Topic**  
   - Select the topic from the list.
   - View details about partitions, replication, and configurations.

4. **Delete a Topic**  
   - Select the topic from the list.
   - Click "Delete" if necessary.

---

## Notes
- **UI Management**: Most operations can be performed using the Kafka UI for convenience.
- **Error Handling**:
  - Ensure all containers are running (`docker ps`).
  - If you encounter connection issues, check network configurations and logs.
- **Default Ports**:
  - Kafka UI: `8080`
  - Kafka Brokers: `9092`, `9093`, `9094`
  - Zookeeper: `2181`, `2182`, `2183`