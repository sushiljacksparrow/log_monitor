# LogFlow

## Project Overview
LogFlow is a real-time distributed log aggregation and search platform. It collects logs from microservices via Kafka, stores them in Elasticsearch, and provides live streaming + historical search through a React web UI.

## Architecture
- **Frontend:** React/TypeScript (Vite), served via Nginx in Docker
- **Backend:** Go microservices communicating via Kafka and gRPC
- **Data flow:** Services -> Kafka topics -> Log Processor -> Elasticsearch -> Query Service (gRPC) -> API Gateway (REST/WebSocket) -> Frontend

## Key Directories
- `frontend/src/` — React app (components, hooks, lib)
- `backend/cmd/` — Service entrypoints (api-gateway, auth-service, order-service, payment-service, log-processor, query-service)
- `backend/internal/` — Shared packages (config, kafka, elasticsearch, websocket, api-gateway, constants, utils)

## Running the Project
```
docker compose up --build -d
```
- Frontend: http://localhost:80
- API Gateway: http://localhost:8000
- Elasticsearch: http://localhost:9200
- Kibana: http://localhost:5601
- Grafana: http://localhost:3000

## Product Roadmap
See [ROADMAP.md](ROADMAP.md) for the full 2-year product plan. When asked to build features, check the roadmap for the next unchecked item in the current quarter. Mark items as complete after shipping.

## Code Conventions
- Backend: Go, standard library style, no frameworks except Gin for HTTP and sarama for Kafka
- Frontend: TypeScript, functional React with hooks, no state management library
- Proto: gRPC service definitions in `backend/internal/query-service/proto/`
- Config: Environment variables loaded via godotenv, with fallback to OS env vars in containers
