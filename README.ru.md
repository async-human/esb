# Enterprise Service Bus

<p align="left">
  <a href="README.md">EN</a> | RU
</p>

## Описание

Платформа Enterprise Service Bus для организации взаимодействия распределённых корпоративных систем на основе событийной архитектуры (EDA). Обеспечивает надёжную асинхронную маршрутизацию, трансформацию и обогащение сообщений между системами с использованием Apache Kafka в качестве шины событий. Реализует паттерны интеграции корпоративных приложений (EIP), гарантируя доставку, обработку ошибок и наблюдаемость.

**Язык:** Go 1.26.1
**Модуль:** `github.com/async-human/esb`


### Инфраструктура

| Компонент | Порт | Описание |
|-----------|------|----------|
| PostgreSQL | 5432 | Реляционная БД для метаданных и конфигурации |
| Kafka (KRaft) | 9092 | Брокер потоковой передачи событий |
| Kafka UI | 8080 | Веб-интерфейс управления кластером Kafka |
| Elasticsearch | 9200 | Хранилище логов и полнотекстовый поиск |
| Kibana | — | Визуализация логов (UI для Elasticsearch) |
| Jaeger | 16686 | Распределённая трассировка (UI + OTLP) |
| Prometheus | 9090 | Сбор и хранение метрик |
| Grafana | 3000 | Визуализация и дашборды |
| OpenTelemetry Collector | — | Приём OTLP, конвейер телеметрии |

Стек наблюдаемости: OpenTelemetry SDK → OTLP Collector → Jaeger (трейсы), Elasticsearch (логи), Prometheus (метрики).

### Развёртывание

Директория `deployment/` содержит конфигурации Docker Compose для инфраструктуры и сервисов. Каждый сервис имеет поддиректорию с `.env.example`, `docker-compose.yml` и `Dockerfile`.

Базовые компоненты инфраструктуры определены в `deployment/core/` с модульными compose-файлами по компонентам. Конфигурации Kubernetes и Terraform находятся в `deployment/core/k8s/` и `deployment/core/terraform/`.

### Конфигурация

Конфигурация сервисов через `.env` файлы в `deployment/<service>/`. Глобальные переменные окружения определены в `deployment/.env`. Сеть: bridge-драйвер (`esb_network`).

## Разработка

### Требования

- Go 1.26+
- [go-task](https://taskfile.dev/)
- Docker и Docker Compose

### Команды

```bash
task infra:up              # Запустить все контейнеры инфраструктуры
task infra:down            # Остановить и удалить контейнеры
task run SERVICE=<имя>     # Запустить сервис локально с env из deployment/<service>/.env
```
