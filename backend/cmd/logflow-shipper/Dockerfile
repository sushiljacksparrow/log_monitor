# BUILDER STAGE

FROM golang:1.25-alpine AS builder
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOFLAGS=-buildvcs=false

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

RUN go build -v -o logflow-shipper ./cmd/logflow-shipper/

# FINAL STAGE
FROM alpine:latest AS final
WORKDIR /app

COPY --from=builder /app/logflow-shipper .
CMD ["./logflow-shipper"]
