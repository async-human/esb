# Enterprise Service Bus

<p align="left">
  EN | <a href="README.ru.md">RU</a>
</p>

## Project Overview

Enterprise Service Bus platform for distributed systems integration based on Event-Driven Architecture (EDA). Provides reliable asynchronous message routing, transformation, and enrichment between decoupled systems using Apache Kafka as the event backbone. Implements Enterprise Integration Patterns (EIP) with guaranteed delivery, error handling, and observability.

## Infrastructure

The `deployment/` directory contains infrastructure components managed via Docker Compose and Taskfile automation.

### Components

| Component | Port | Description |
|-----------|------|-------------|
| PostgreSQL | 5432 | Relational database for metadata and configuration |
| Kafka (KRaft) | 9092 | Event streaming broker |
| Kafka UI | 8080 | Web UI for Kafka cluster management |
| Elasticsearch | 9200 | Log storage and full-text search |
| Jaeger | 16686 | Distributed tracing UI and OTLP endpoint |
| Prometheus | 9090 | Metrics collection and storage |
| Grafana | 3000 | Dashboards and visualization |

### Management Commands

Infrastructure lifecycle is managed through [go-task](https://taskfile.dev/):

```bash
task infra:up        # Start all infrastructure containers
task infra:down      # Stop and remove containers
task infra:clean     # Remove generated .env files
task infra:env       # Generate component-specific .env files from deployment/.env
```

### Configuration

Global environment variables are defined in `deployment/.env`. Network uses `172.25.0.0/16` subnet with bridge driver (`esb_network`).