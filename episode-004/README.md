# Event-Driven Architecture Tutorial – Episode 4

YouTube ***(Hungarian)***: https://youtu.be/p-zJLZp509E

This repository contains the source code for Episode 4 of my YouTube tutorial series on Event-Driven Architecture (EDA) using Go and RabbitMQ.

---

## Overview

In this episode we add Idempotency to the other services by idempotency deduplication key store

---

## Example test payload for de-duplication testing: (copy into RabbitMQ console, Queue: `order-created.order-service`)
```json
{
  "metadata": {
    "action_type": "order-created",
    "action_id": "b0732a43-8216-4821-8240-ac052587e692",
    "causation_id": "",
    "correlation_id": "b0732a43-8216-4821-8240-ac052587e692",
    "index": 0,
    "occurred_at": "2026-01-18T13:18:57.615060924Z"
  },
  "payload": {
    "uuid": "5b1dd483-a380-4ab5-8751-8c0feacccda0",
    "userID": "XDSFDXASAs",
    "email": "joen@company.com",
    "items": [
      { "productID": "A1215", "quantity": 1 },
      { "productID": "A1216", "quantity": 4 }
    ]
  }
}
```

---

## Useful queries for this episode
```sql
-- Statistics, event count per correlation
SELECT correlation_id, count(*) from logs l group by correlation_id;

--  Initial events
SELECT * FROM logs l WHERE causation_id = "";

-- First child events
SELECT * FROM logs l WHERE correlation_id  = causation_id; 

-- Decuplication data
SELECT * from processed_actions pa;

-- Invoices
SELECT * from invoice_heads ih;

-- Orders
SELECT * from order_heads oh; 

-- Get a chain of events (Use your correlation ID for investigation)
SELECT * FROM logs where correlation_id = "d4e5b6cc-07de-41b8-803b-6023c078b20c" order by created_at;

-- Get a chain of error events (Use your correlation ID for investigation)
SELECT * FROM logs where correlation_id = "d4e5b6cc-07de-41b8-803b-6023c078b20c" AND level = 'error' order by created_at;
```
