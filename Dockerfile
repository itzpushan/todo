FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build main server binary
RUN go build -o server cmd/server/main.go

# Build seed binary
RUN go build -o seed cmd/seed/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/seed /app/seed
COPY .env .env

EXPOSE 8000

CMD sh -c "/app/seed && echo '✅ Seeded successfully' && /app/server"
