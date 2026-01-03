# Event-Driven Architecture Tutorial – Episode 2

YouTube ***(Hungarian)***: todo

This repository contains the source code for Episode 2 of my YouTube tutorial series on Event-Driven Architecture (EDA) using Go and RabbitMQ.

---

## Overview

In this episode, we build a complete **event-driven application** using a **microservices architecture**.

The project includes a Bash script that spins up all required Docker containers, allowing the system to run locally in a way that closely resembles a real cloud deployment.

---


## Application Workflow

### Producers:
- HTTP API Service Sends an OrderCreated event when a new order is created.

### Consumers:


- Store Order Service
  - Consumes: `OrderCreated`
  - Stores the order
  - Updates stock
  - Produces: `OrderPersisted`
- Payment Service
  - Consumes: `OrderPersisted`
  - Processes payment
  - Produces
    - `PaymentSucceeded` on success
    - `paymentFailed` on failure
- Cancel Order Service
  - Consumes: `paymentFailed`
  - Cancels the order  
  - Restores stock
- Send Cancel Order Email Service
  - Consumes: `paymentFailed`
- Invoice Service
  - Consumes: PaymentSucceeded
  - Creates an invoice
  - Produces: InvoiceCreated
- Send Invoice Email Service  
 - Consumes: InvoiceCreated
- Log message Created
  - Store log in database, any other consumer or producer can send this event on error or info 


---

## Usage:

Create database tables:
```bash
migrator migrate
```

> Note: if you did not follow the basic course, here is how to setup the migrator: https://github.com/olbrichattila/godbmigrator_cmd



Migrate database:

**Important:** First, build and start the base Docker services.
```bash
docker compose up -d
```

Then build all microservices:
```
./build.sh
```

For a detailed walkthrough, see the accompanying YouTube video.

---


## What next
In this episode, you’ll notice that if an error occurs at the wrong moment:
- Orders may be duplicated
- Emails may be sent multiple times
- Stock levels may become incorrect

In the next episode, we will introduce idempotency and show how to handle these problems correctly in an event-driven system.