# BUILDER
FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN go install github.com/air-verse/air@latest 

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go build -o log-processor ./cmd/log-processor/

# FINAL
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/log-processor .
RUN chmod +x log-processor

CMD ["./log-processor"]