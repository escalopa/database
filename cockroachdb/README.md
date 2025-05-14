- [ ] [serializable-lockless-distributed-isolation-cockroachdb](https://www.cockroachlabs.com/blog/serializable-lockless-distributed-isolation-cockroachdb/)
- [ ] [consensus-made-thrive](https://www.cockroachlabs.com/blog/consensus-made-thrive/)
- [ ] [trust-but-verify-cockroachdb-checks-replication](https://www.cockroachlabs.com/blog/trust-but-verify-cockroachdb-checks-replication/)
- [ ] [living-without-atomic-clocks](https://www.cockroachlabs.com/blog/living-without-atomic-clocks/)

## Setup

Run cluster:

```bash
docker compose up -d
``` 

Init cluster:

```bash
docker exec -it roach1 cockroach --host=roach1:26357 init --insecure
```

Connect to DB:

```bash
docker exec -it roach1 cockroach sql --host=roach1:26257 --insecure
```

Create random data:

```sql
-- Step 1: Create the database
CREATE DATABASE IF NOT EXISTS shop;

-- Step 2: Switch to the `shop` database
USE shop;

-- Step 3: Create the `orders` table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_name STRING NOT NULL,
    item STRING NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    order_date TIMESTAMP NOT NULL DEFAULT now()
);

-- Step 4: Generate and insert 10,000 orders
-- CockroachDB supports `generate_series()` for batch insert

INSERT INTO orders (customer_name, item, quantity, price, order_date)
SELECT
    'Customer_' || gs::STRING,
    ARRAY['Laptop', 'Phone', 'Tablet', 'Headphones', 'Monitor'][((random()*4)::INT + 1)],
    (random()*10)::INT + 1,
    round((random() * 900 + 100)::NUMERIC, 2),
    now() - (random() * INTERVAL '30 days')
FROM generate_series(1, 10000) AS gs;
```

---

## General:

- CockroachDB is Postgres compatible, distributed SQL database.
- Received SQL RPC query is converted into KV operations.
- Splits data into ranges and distributes them across nodes (at least 3 nodes).
- If a node receives a read or write request it cannot directly serve, it finds the node that can handle the request, and communicates with that node.
- Leaseholders are responsible for serving reads and writes for a range of data (raft leader of the range).
- A transaction is unable to complete due to another concurrent or recent transaction attempting to write to the same data. This is also called lock contention.

### Layers

| Layer         | Order | Purpose                                                                                                                                     |
|---------------|-------|---------------------------------------------------------------------------------------------------------------------------------------------|
| SQL           | 1     | Translate client SQL queries to KV operations.                                                                                              |
| Transactional | 2     | Allow atomic changes to multiple KV entries.                                                                                                |
| Distribution  | 3     | Present replicated KV ranges as a single entity.                                                                                            |
| Replication   | 4     | Consistently and synchronously replicate KV ranges across many nodes. This layer also enables consistent reads using a consensus algorithm. |
| Storage       | 5     | Read and write KV data on disk.                                                                                                             |

---

## SQL Layer

| Component            | Description                                                                                     |
|----------------------|-------------------------------------------------------------------------------------------------|
| SQL API              | Forms the user interface.                                                                       |
| Parser               | Converts SQL text into an abstract syntax tree (AST).                                           |
| Cost-based optimizer | Converts the AST into an optimized logical query plan.                                          |
| Physical planner     | Converts the logical query plan into a physical query plan for execution by one or more nodes.  |
| SQL execution engine | Executes the physical plan by making read and write requests to the underlying key-value store. |

Constant folding: replacing expressions with their computed values at compile time (e.g., `1 + 2` becomes `3`).

### Logical Plan

**Example (simplified):**

```
SELECT name FROM users WHERE age > 30
```

```sql
  EXPLAIN SELECT ...;
```

```
→ Filter (age > 30)
  → Scan users
```

### Physical Plan

**Example (simplified):**

```sql
  EXPLAIN (DISTSQL) SELECT ...;
```

```
→ DistSQL Plan:
  - Node 1: Scan users (range 1–10)
  - Node 2: Scan users (range 11–20)
  - Coordinator: Apply filter (age > 30), combine results
```
 
---

## Transaction Layer

### Write Transaction

A write transaction is split into set of operations:

- Write intents
- UnReplicated locks
- Transaction record

> As write intents are created, CockroachDB checks for newer committed values.
> If newer committed values exist, the transaction may be restarted. (dangerous)
> If existing write intents or locks exist on the same keys,
> it is resolved as a transaction conflict.

### Read Transaction

Provided types of reads:

- **Strongly consistent reads:** Reads the latest committed value.
- **Stale reads:** Reads the latest committed value, but may not be the most recent.

Reads have 2 types of locks:

- **Exclusive locks:** Prevents other transactions from writing/reading the same key.
- **Shared locks:** Prevents other transactions from writing/reading(using exclusive locks like `SELECT FOR UPDATE`) the same key, but allows normal reading.

### Time & Hybrid Digital Clock

For reads CockroachDB uses HCL (Hybrid Logical Clock).

HCL is a combination:
- physical component (always close to wall clock time).
- logical component (used to compare events within the same physical component).

> HLC timestamp is always >= wall clock timestamp ([paper](https://cse.buffalo.edu/tech-reports/2014-04.pdf))

Max Clock Offset Enforcement:

When a node detects that its clock is out of sync with at least half of the
other nodes in the cluster by 80% of the maximum offset allowed, 
it crashes immediately. This can be prevented using NTP (Network Time Protocol).

---

## Distribution Layer 

---

## Replication Layer 

---

## Storage Layer

---

## TODO

- Check transaction serializable retries.
- Check DistSQL
