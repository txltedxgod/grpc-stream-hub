# grpc-stream-hub

> High-throughput **bidirectional gRPC streaming hub** with Protocol Buffers schema, metadata auth interceptors, and fanout pub/sub broadcasting in **Go**.

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![gRPC](https://img.shields.io/badge/Protocols-gRPC%20%7C%20Protobuf-244C5A?style=flat-square&logo=grpc)](https://grpc.io)
[![CI](https://img.shields.io/badge/CI-Passing-238636?style=flat-square&logo=githubactions)](https://github.com/txltedxgod/grpc-stream-hub/actions)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)](https://docker.com)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)

`#grpc` `#protobuf` `#bidirectional-streaming` `#pubsub` `#golang` `#microservices` `#high-throughput`

---

## 🏛️ Bidirectional Streaming Architecture

```mermaid
sequenceDiagram
    autonumber
    participant ClientA as Client Alice (gRPC)
    participant Interceptor as Auth Stream Interceptor
    participant Hub as Central Hub Broadcast Engine
    participant ClientB as Client Bob (gRPC)

    ClientA->>Interceptor: Open Bidirectional Stream (Header: x-api-token)
    Interceptor->>Hub: Register ClientA channel for Room: general
    ClientB->>Interceptor: Open Bidirectional Stream (Header: x-api-token)
    Interceptor->>Hub: Register ClientB channel for Room: general

    ClientA->>Hub: Send MessageRequest { room: "general", content: "Hello!" }
    Hub->>Hub: Fanout to all active room subscribers
    Hub-->>ClientA: Echo MessageResponse { sender: "Alice", content: "Hello!" }
    Hub-->>ClientB: Stream MessageResponse { sender: "Alice", content: "Hello!" }
```

---

## Features

- **Full Bidirectional Streaming:** High-efficiency HTTP/2 multiplexed streams via standard gRPC.
- **Topic / Room Isolation:** Central lock-striped hub manages dynamic subscriptions without blocking publishers.
- **Stream Interceptors:** Metadata token validation and request tracing on stream initialization.
- **Slow Consumer Protection:** Non-blocking broadcast drops slow buffers to prevent head-of-line blocking.

## Quick Start

```bash
# 1. Run gRPC Hub Server
go run server/main.go server/hub.go server/interceptor.go -port=50051

# 2. Run Interactive Client
go run client/main.go -server=localhost:50051 -room=general -user=alice
```
