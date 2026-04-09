FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем platform и сервис
COPY platform /app/platform
COPY services/management-api /app/services/management-api

# Создаём workspace и добавляем модули
RUN go work init && go work use ./platform ./services/management-api

# Копируем .env файл
COPY deployment/management-api/.env /app/.env

WORKDIR /app/services/management-api
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/management-api ./cmd/app/

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/management-api ./management-api
COPY --from=builder /app/.env ./deployment/management-api/.env

EXPOSE 8080

CMD ["./management-api"]
