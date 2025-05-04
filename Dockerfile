FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl tar
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY .env .
COPY db/migrations ./migrations
RUN chmod +x ./migrate

EXPOSE 3000
CMD ["sh", "-c", "echo 'Current directory:'; pwd; echo 'Migrations directory contents:'; ls -la /app/migrations; echo 'Running migrations...'; /app/migrate -path /app/migrations -database \"$DB_SOURCE\" -verbose up && echo 'Starting application...' && /app/main"]
