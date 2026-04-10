# Enterprise Service Bus

<p align="left">
  EN | <a href="README.ru.md">RU</a>
</p>

## Overview

Enterprise Service Bus platform for distributed systems integration based on Event-Driven Architecture (EDA). Provides reliable asynchronous message routing, transformation, and enrichment between decoupled systems using Apache Kafka as the event backbone. Implements Enterprise Integration Patterns (EIP) with guaranteed delivery, error handling, and observability.

**Language:** Go 1.26.1
**Module:** `github.com/async-human/esb`

### Infrastructure

| Component | Port | Description |
|-----------|------|-------------|
| PostgreSQL | 5432 | Relational database for metadata and configuration |
| Kafka (KRaft) | 9092 | Event streaming broker |
| Kafka UI | 8080 | Web UI for Kafka cluster management |
| Elasticsearch | 9200 | Log storage and full-text search |
| Kibana | 5601 | Log visualization (Elasticsearch UI) |
| Jaeger | 16686 | Distributed tracing UI and OTLP endpoint |
| Prometheus | 9090 | Metrics collection and storage |
| Grafana | 3000 | Dashboards and visualization |
| OpenTelemetry Collector | 4317 / 4318 | OTLP receiver (gRPC / HTTP), telemetry pipeline |

#### Services (local access)

| Service | Port | Swagger UI |
|---------|------|------------|
| inbound-connector | 8081 | http://localhost:8081/docs |

Observability stack: OpenTelemetry SDK → OTLP Collector → Jaeger (traces), Elasticsearch (logs), Prometheus (metrics).

### Deployment

The `deployment/` directory contains Docker Compose configurations for infrastructure and services. Each service has its own subdirectory with `.env.example`, `docker-compose.yml`, and `Dockerfile`.

Infrastructure core components are defined under `deployment/core/` with modular compose files per component. K8s and Terraform configurations are present under `deployment/core/k8s/` and `deployment/core/terraform/`.

### Configuration

Service configuration via `.env` files in `deployment/<service>/`. Global environment variables defined in `deployment/.env`. Network: bridge driver (`esb_network`).

## Development

### Prerequisites

- Go 1.26+
- [go-task](https://taskfile.dev/)
- Docker & Docker Compose

### Commands

```bash
task infra:up              # Start all infrastructure containers
task infra:down            # Stop and remove containers
task run SERVICE=<name>    # Run a service locally with env from deployment/<service>/.env
```

### API Generation

```bash
task api:generate SERVICE=<name>   # Generate Go code from OpenAPI spec (reads from shared/api/<name>/v1/)
task api:gen-inbound               # Shorthand: generate API for inbound-connector
```

Swagger UI for running services: `http://localhost:<port>/docs` (e.g. http://localhost:8081/docs for inbound-connector).
