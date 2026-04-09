FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем platform и сервис
COPY platform /app/platform
COPY services/outbound-connector /app/services/outbound-connector

# Создаём workspace и добавляем модули
RUN go work init && go work use ./platform ./services/outbound-connector

# Копируем .env файл
COPY deployment/outbound-connector/.env /app/.env

WORKDIR /app/services/outbound-connector
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/outbound-connector ./cmd/app/

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/outbound-connector ./outbound-connector
COPY --from=builder /app/.env ./deployment/outbound-connector/.env

EXPOSE 8080

CMD ["./outbound-connector"]
