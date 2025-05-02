# Task Management Service

## Table of Contents
- [Video Walkthrough](#video-walkthrough)
- [Tech Stack](#tech-stack)
- [Setup Steps](#setup-steps)
- [API Usage Examples](#api-usage-examples)
- [Technical Overview](#technical-overview)

## Video Walkthrough
[Link to video walkthrough]

## Tech Stack
- Golang
- Gin
- Gorm
- Postgres

## Setup Steps
1. Clone the repository:
   ```
   git clone https://github.com/dilawari2008/Golang-Task-Management-Service.git
   ```

2. Navigate to the project directory:
   ```
   cd task-management-system
   ```

3. Generate gRPC code:
   ```
   protoc --proto_path=api/proto \
     --go_out=paths=source_relative:api/proto \
     --go-grpc_out=paths=source_relative:api/proto \
     api/proto/task.proto
   ```

4. Run the application:
   ```
   go run main.go
   ```

## API Usage Examples
### Create a Task
```
curl --location 'localhost:8080/api/tasks/' \
--header 'Content-Type: application/json' \
--data '{
  "topic": "consumer2",
  "data": "consumer2 data"
}'
```

### Get tasks list by filtering and pagination
```curl --location 'localhost:8080/api/tasks?status=completed&page=1&limit=10'```

### Get task by id
```curl --location 'localhost:8080/api/tasks/25'```

### Update task by id
```curl --location --request PUT 'localhost:8080/api/tasks/24' \
--header 'Content-Type: application/json' \
--data '{
  "status": "pending"
}'```

### Delete task by id
```curl --location --request DELETE 'localhost:8080/api/tasks/24' \
--header 'Content-Type: application/json'```


# Technical Overview: Task Processing System

## Core Architecture

1. **Task Processing Lifecycle**
   - Each task is processed in 3 phases: pre-processing, topic-specific processing, and post-processing
   - Pre-processing and post-processing phases are common for all tasks
   - Processing phase is topic-specific and implements the unique business logic

2. **Task Status Management**
   - Status transitions follow a strict hierarchy and can only move forward:
     - **Tier 1:** Pending (default status for new tasks)
     - **Tier 2:** InProgress
     - **Tier 3:** Completed, Failed, or Expired (terminal states)
   - Tasks in non-terminal states (Pending, InProgress) will expire after a defined time period relative to their `updatedAt` timestamp

3. **Topic System**
   - Tasks must have both a topic and associated data for processing
   - To add a new topic:
     - Add the topic to constants
     - Implement its specific processing logic in `services/impl`
   - Topic and data can only be updated if the task is not in a terminal state (Completed, Failed, Expired)

## Implementation Details

4. **Queue Management**
   - Currently using in-memory queue for message processing
   - Design supports horizontal scaling via environment variables that control service behavior:
     - Task creation via REST endpoints
     - Task processing from main task queue
     - Failed task processing from Dead Letter Queue (DLQ)

5. **Error Handling & Callbacks**
   - For external services consuming messages from queues, tasks are marked as completed using callbacks
   - Failed task details and error messages are persisted in the database
   - DLQ handles failed task processing

6. **Data Persistence**
   - Delete operations are implemented as soft deletes to maintain data integrity and audit history
   - Task status changes and error conditions are logged in the database

## Infrastructure & Scaling

7. **Security Considerations**
   - Authentication might be unnecessary if the service operates within a private subnet of a VPC:
     - No connection to internet gateway (blocks inbound traffic)
     - Connected to NAT Gateway for outbound traffic
   - Consider implementing additional security based on deployment environment

8. **External Integrations**
   - External service responses are currently mocked in consumer implementations
   - Designed to be replaced with actual service integrations

## Development Notes

9. **Testing Approach**
   - Service implementations include mocks for external dependencies
   - Unit tests should verify the complete task lifecycle
   - Consider implementing integration tests that validate the end-to-end task processing

10. **Future Enhancements**
    - Replace in-memory queue with durable message queue (e.g., SQS, RabbitMQ, Kafka)
    - Implement metrics and monitoring for task processing
    - Add retry policies with exponential backoff for transient failures
