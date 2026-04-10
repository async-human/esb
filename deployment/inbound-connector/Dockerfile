FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем корневой go.mod/go.sum
COPY go.mod go.sum /app/

# Копируем shared API (generated code)
COPY pkg /app/pkg

# Копируем platform и сервис
COPY platform /app/platform
COPY services/inbound-connector /app/services/inbound-connector

# Копируем .env файл
COPY deployment/inbound-connector/.env /app/.env

# Создаём workspace: корневой модуль нужен для pkg/ (generated API code)
RUN go work init && go work use . ./platform ./services/inbound-connector

WORKDIR /app/services/inbound-connector
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/inbound-connector ./cmd/app/

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/inbound-connector ./inbound-connector
COPY --from=builder /app/.env ./deployment/inbound-connector/.env

EXPOSE 8080

CMD ["./inbound-connector"]
