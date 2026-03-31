# BUILDER STAGE

FROM golang:1.25-alpine AS builder
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOFLAGS=-buildvcs=false

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go build -v -o auth-service ./cmd/auth-service/


# FINAL STAGE
FROM alpine:latest AS final
WORKDIR /app

COPY --from=builder /app/auth-service .
EXPOSE 8080
CMD ["./auth-service"]
