FROM golang:1.20 AS builder

WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

# Добавляем переменную среды окружения для подключения к PostgreSQL
ENV DB_CONNECTION_STRING="postgres://god_user:ASDqwe123___23Dfd@164.92.179.123:5050/postgres?sslmode=disable"
# для прода
# ENV DATABASE_URL="postgres://srv-captain--postgres-db?port=5432&dbname=postgres&user=god_user&password=ASDqwe123___23Dfd"

ENTRYPOINT ["/app"]
