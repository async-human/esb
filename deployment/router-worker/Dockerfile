FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем platform и сервис
COPY platform /app/platform
COPY services/router-worker /app/services/router-worker

# Создаём workspace и добавляем модули
RUN go work init && go work use ./platform ./services/router-worker

# Копируем .env файл
COPY deployment/router-worker/.env /app/.env

WORKDIR /app/services/router-worker
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/router-worker ./cmd/app/

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/router-worker ./router-worker
COPY --from=builder /app/.env ./deployment/router-worker/.env

EXPOSE 8080

CMD ["./router-worker"]
