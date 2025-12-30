# Event-Driven Architecture Tutorial – Episode 1

YouTube ***(Hungarian)***: https://youtu.be/Q1Hk3L1Cyyo

This repository contains the source code for Episode 1 of my YouTube tutorial series on Event-Driven Architecture (EDA) using Go and RabbitMQ.

---

## Overview

In this episode, we cover:
**Architectural Designs:**
- Monolith
- Modular Monolith
- Microservices

**Introduction to Event-Driven Architecture (EDA):**
- What EDA is
- Use cases where different architectures are appropriate

**Core Concepts:**
- Queues and messaging patterns
- Technologies like RabbitMQ, Amazon SNS, SQS, and a mention of Kafka (note: Kafka is not a queue)

**Common Pitfalls:**
- The frequently misunderstood pattern called “Passive-Aggressive Events”
**Code in this Episode**
- This episode demonstrates a basic Go implementation:
- A simple RabbitMQ producer
- A simple RabbitMQ consumer

> Note: The code uses replace in go.mod to access the shared module for tutorial purposes. In a real-world setup, the shared module would be in a separate repository.

---

## Next Steps
In the next episode, this code will evolve into a full event-driven ordering system, as explained in the video.

---

## Getting Started
- Make sure you have Go installed (1.25+ recommended).
- Install and run RabbitMQ locally (or use Docker).
- Run the producer and consumer as shown in the video.

---

## References
- RabbitMQ official site: https://www.rabbitmq.com/
- Amazon SNS: https://aws.amazon.com/sns/
- Amazon SQS: https://aws.amazon.com/sqs/
- Kafka: https://kafka.apache.org/
