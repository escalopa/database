## Setup

```bash
go install github.com/nats-io/natscli/nats@latest
```

```bash
docker run -p 4222:4222 -p 8222:8222 -p 65432:65432 -v ./data/nats:/nats -v ./js.conf:/js.conf nats:latest -js -m 8222 -c /js.conf
```
```
```

## Quality Of Service (QOS)

### At most once

Core NATS is a fire-and-forget messaging system.
This means that messages are sent without any guarantee of delivery.
This is the default behavior of NATS and is suitable for use cases where message loss is acceptable.

### At least once

NATS Streaming provides at-least-once delivery semantics.

### Exactly once

NATS JetStream provides exactly-once delivery semantics.

### Subject

At its simplest, a subject is just a string of characters that form a name
the publisher and subscriber can use to find each other.

### Wildcards

- Publishers will always send a message to a fully specified subject
- Subscribers can subscribe to a subject with wildcards

For example see Subject Hierarchy:

### Subject Hierarchy

Below is an example of a subject hierarchy:

```text
time.us
time.us.east
time.us.east.atlanta
time.eu.east
time.eu.east.warsaw
```

A subscriber can subscribe to:

- `time` - all subjects
- `time.us` - all subjects under `time.us` (i.e. `time.us.east`, `time.us.east.atlanta`)
- etc...

#### Recommendations

Use at most

- 16 tokens
- 256 characters length

Example:

```text
time.us.east.atlanta
```

- Tokens: 4 (i.e `time`, `us`, `east`, `atlanta`)
- Length: 20

### Numbers Of Subjects

NATS can manage 10s of millions of subjects.

### Match Tokens

#### Star Token `*`

The star token `*` matches a single token in a subject.

Example:

```text
time.*.east
```

This will match:
- `time.us.east`
- `time.eu.east`

But doesn't match a substring within a token:

```text
time.New*.east
```

#### Greater Than Token `>`

The greater than token `>` matches all remaining tokens in a subject (must appear at the end of the subject).

Example:

```text
time.us.east.>
```
This will match:
- `time.us.east`
- `time.us.east.atlanta`
- `time.us.east.warsaw`


### Allowed Characters

- Alphanumeric characters (`a-z`, `A-Z`, `0-9`)
- ASCII characters (`-`, `_`)
- Token separator (`.`)

### Reserved Characters

- `$` - used for system subjects
- `>` - used for wildcard
- `*` - used for wildcard

### Pedantic Mode

By default, for the sake of efficiency, subject names are not verified during message publishing.

To enable subject name verification, activate pedantic mode in the client connection options.

### Pub/Sub

Subscribe to a subject:

```bash
nats sub <subject> [options]
```

Example usage:

```bash
nats sub time.us.east
```

Publish to a subject

```bash
nats pub <subject> <message> [options]
```

Example usage:

```bash
nats pub time.us.east "Hello World"
```

### Request / Reply

Reply to a subject:

```bash
nats reply <subject> [options]
```

Example usage:

```bash 
    nats reply time.us.east
```

Request to a subject:

```bash
nats request <subject> <message> [options]
```

Example usage:

```bash
nats request time.us.east "Hello World"
```

### Queue Groups

Queue groups are a way to group subscribers together so that only 
one subscriber in the group receives a message.

Consider the following example:

Create 3 subscribers:

```bash
nats sub time.us.east --queue group1
```

```text
nats sub time.us.east --queue group1
```

```text
nats sub time.us.east --queue group1
```

Send 1k messages to the subject:

```bash
nats pub time.us.east "Hello World" --count 1000
```

## JetStream

### Replay Policy

| **Replay Policy**                        | **Description**                                                                                                             |
|------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------|
| **All Messages: Replay Speed: Instant**  | Replay all messages currently stored in the stream, Messages are delivered to the consumer as fast as it can take them.     |
| **All Messages: Replay Speed: Original** | Replay all messages currently stored in the stream, Messages are delivered at the rate they were published into the stream. |
| **Last Message**                         | Replay the last message stored in the stream.                                                                               |
| **Last Message Per Subject**             | Replay the last message for each subject (streams can capture more than one subject).                                       |
| **From Specific Sequence Number**        | Replay starting from a specific sequence number.                                                                            |
| **From Specific Start Time**             | Replay starting from a specific start time.                                                                                 |

### Retention

- Maximum message age.
- Maximum total stream size (in bytes).
- Maximum number of messages in the stream.
- Maximum individual message size.
- You can also set limits on the number of consumers that can be defined for the stream at any given point in time.

### Discard Policy

Discarding messages happens when the stream reaches its retention limits (any of the above).

Available strategies:
- discard old messages
- discard new messages (publisher gets an error that we reached the limit)

### Retention Policy

| **Retention Policy** | **Description**                                                                                                                                                                                                     |
|----------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **limits**           | The default policy to provide a replay of messages in the stream.                                                                                                                                                   |
| **work queue**       | The stream is used as a shared queue, and messages are removed as consumed, providing exactly-once consumption.                                                                                                     |
| **interest**         | Messages are kept in the stream as long as there are consumers that haven't delivered the message yet. This is a variation of work queue that retains messages only if there is interest for the message's subject. |

### Persistence

- Memory
- File
- Replication (1 (no replication), 2, 3, ...)

### Ack

- Sync
- Async
- NAK (negative acknowledgment)
- in progress

### EOS (Exactly Once Semantics)

Publisher: Uses UNIQUE publish_id per message.
Subscriber: Uses double acknowledgment mechanism.

## Consumer

### Fast Push

Fast push consumer usually a good fit for replay
messages (doesn't require acknowledgment)

### Pull

### KV

Key-Value store,
Supports atomic CRUD operations and watch for changes (similar to ETCD).

### Mirrors

- Allow you to replicate a stream's messages into another stream without re-publishing them manually.
- They are read-only copies of another stream
- Replication in done asynchronously
- Messages can be deleted
- Only one stream available

## KV Store

Features:
- Atomic operations
1) Atomic operations
2) GET/SET/DEL/PURGE operations
3) Watch for changes
4) History
5) TTL (Time to Live)

## Object Store

Features:

- PUT/GET/DEL operations
- Watch for changes

## Mapping & Partitioning

- Mapping: Map a subject to a stream
- Partitioning: Partition a stream into multiple subjects

Example:

```text
time.us.east.>
```
## Profiling

### CPU

```bash
nats server request profile cpu [options]
```

| **Option**              | **Description**                                                              |
|-------------------------|------------------------------------------------------------------------------|
| *(no option)*           | Request a CPU profile from all servers in the system                         |
| `./profiles`            | Request a CPU profile from all servers and write to the `profiles` directory |
| `--timeout=10s`         | Request a CPU profile from all servers over a 10 second period               |
| `--name=servername1`    | Request a CPU profile from `servername1` only                                |
| `--tags=aws`            | Request a CPU profile from all servers tagged as `aws`                       |
| `--cluster=aws-useast2` | Request a CPU profile from all servers in the `aws-useast2` cluster only     |

### Memory

```bash
nats server request profile allocs [options]
```

| **Option**              | **Description**                                                                 |
|-------------------------|---------------------------------------------------------------------------------|
| *(no option)*           | Request a memory profile from all servers in the system                         |
| `./profiles`            | Request a memory profile from all servers and write to the `profiles` directory |
| `--name=servername1`    | Request a memory profile from `servername1` only                                |
| `--tags=aws`            | Request a memory profile from all servers tagged as `aws`                       |
| `--cluster=aws-useast2` | Request a memory profile from all servers in the `aws-useast2` cluster only     |
