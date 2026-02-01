# EDA tutorial

This is a work in progress repository for EDA Event Driven tutorial video series (on my ***Hungarian YouTube*** channel)

> Note: The producers and consumers use replace in go.mod to access the shared module. In a real-world setup, this shared module would live in a separate GitHub repository and replace would not be used. It is included here only to keep all the code in one place for the sake of this tutorial.

---

## Episode 1:
Simple **RabbitMQ consumer** and **producer**.
- YouTube ***(Hungarian)***: https://youtu.be/Q1Hk3L1Cyyo
- Code: https://github.com/olbrichattila/edatutorial/tree/main/episode-001

---

## Episode 2:
In this episode, we build a complete **event-driven application** using a **microservices architecture**.
- YouTube ***(Hungarian)***: https://youtu.be/lOQd9QgGmFw
- Code: https://github.com/olbrichattila/edatutorial/tree/main/episode-002

---

## Episode 3:
In this episode, we make the **order-service** idempotent.
- YouTube ***(Hungarian)***: https://youtu.be/ge9yfBp2io8
- Code: https://github.com/olbrichattila/edatutorial/tree/main/episode-003

---

## Episode 4:
In this episode, we make other services idempotent with action store
- YouTube ***(Hungarian)***: https://youtu.be/p-zJLZp509E
- Code: https://github.com/olbrichattila/edatutorial/tree/main/episode-004

---


## Episode 5:
In this episode, we handle database and business level race conditions in stock management
Using Optimistic versioned race condition management
- YouTube ***(Hungarian)***: https://youtu.be/D9dcft8hCwo
- Code: https://github.com/olbrichattila/edatutorial/tree/main/episode-005

---

## Episode 6:
In this episode, I introduce Dead Letter Queues (DLQ) using RabbitMQ. We’ll cover the core concepts and walk through a hands-on demo in the RabbitMQ Management Console. In the next episode, we’ll integrate DLQ handling into our Golang codebase.
- YouTube ***(Hungarian)***: Coming soon
- Code: https://github.com/olbrichattila/edatutorial/tree/main/episode-006

---