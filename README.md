# Event-Driven Scheduler

## Introduction

This is a simple event-driven scheduler that consits of 2 parts:
1. *go-job-runner*: a simple job runner that triggers worker by sending a message to a queue
2. *go-job-worker*: a simple worker that receives a message from a queue and executes a job

## How to run

### Prerequisites

- Docker
- Docker Compose

### Run

1. Run `docker-compose up -d` to start the scheduler
2. Open `http://localhost:4001/swagger/index.html` to see the swagger documentation
3. Register worker with name `random`
4. Copy worker id from the response
5. Launch a worker with id from step 4
6. Copy job id from the response
7. Check job statuses by job id from step 6

## Architecture

Scheduler consists of 2 parts: job runner and job worker. 
- Job runner is responsible for triggering workers by sending messages to a queue. 
- Job worker is responsible for receiving messages from a queue and executing jobs. 
- RabbitMQ is used as a message broker. 
- PostgreSQL is used as a database. 
- Flyway is used for database migrations. 
- Swagger is used for API documentation.

## How to scale

To scale the scheduler, we can run multiple instances of job runner and job worker. To balance HTTP requests between job runners, we can use a load balancer. We can run as many instances of job worker as we need. Multiple workers will be able to receive messages from a queue and execute jobs. RabbitMQ will balance messages between different workers.
