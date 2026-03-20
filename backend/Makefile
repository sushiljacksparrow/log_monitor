.PHONY: auth log order payment query websocket clean \
        docker-build docker-up docker-down docker-logs \
        docker-auth docker-payment docker-log docker-order docker-query docker-websocket \
        dev dev-no-cache

# ─────────────────────────────────────
# LOCAL — run directly on your machine
# ─────────────────────────────────────
CMD_PATH = ./cmd
BIN_PATH = ./bin

auth:
	cd $(CMD_PATH)/auth-service && air

log:
	cd $(CMD_PATH)/log-processor && air

order:
	cd $(CMD_PATH)/order-service && air

payment:
	cd $(CMD_PATH)/payment-service && air

query:
	cd $(CMD_PATH)/query-service && air

websocket:
	cd $(CMD_PATH)/websocket-server && air

clean:
	rm -rf $(BIN_PATH)

# ─────────────────────────────────────
# DOCKER — run inside containers
# ─────────────────────────────────────

# build + start everything
dev:
	docker compose up --build

dev-no-cache:
	docker compose build --no-cache
	docker compose up

# start everything (no rebuild)
up:
	docker compose up

# stop everything
down:
	docker compose down

# stop everything + wipe volumes (nuclear option)
reset:
	docker compose down -v

# rebuild everything from scratch
build:
	docker compose build

# ─────────────────────────────────────
# DOCKER — start specific services
# ─────────────────────────────────────

docker-auth:
	docker compose up --build auth-service

docker-payment:
	docker compose up --build payment-service

docker-log:
	docker compose up --build log-processor
# 	docker compose up log-processor

docker-order:
	docker compose up --build order-service

docker-query:
	docker compose up --build query-service

docker-websocket:
	docker compose up --build websocket-server


# ─────────────────────────────────────
# DOCKER — stop specific services
# ─────────────────────────────────────

down-auth:
	docker compose stop auth-service

down-payment:
	docker compose stop payment-service

down-log:
	docker compose stop log-processor

down-order:
	docker compose stop order-service

down-query:
	docker compose stop query-service

down-websocket:
	docker compose stop websocket-server

# ─────────────────────────────────────
# DOCKER — logs
# ─────────────────────────────────────

logs:
	docker compose logs -f

logs-auth:
	docker compose logs -f auth-service

logs-payment:
	docker compose logs -f payment-service

logs-log:
	docker compose logs -f log-processor

logs-order:
	docker compose logs -f order-service

logs-query:
	docker compose logs -f query-service

logs-websocket:
	docker compose logs -f websocket-server

# ─────────────────────────────────────
# DOCKER — infra only (kafka + ES + kibana)
# ─────────────────────────────────────
infra:
	docker compose up zookeeper kafka elasticsearch kibana

# ─────────────────────────────────────
# PROTO → GEN
# ─────────────────────────────────────

PROTO_DIR = ./internal/grpc/proto
GEN_DIR = ./internal/grpc/gen

PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

proto:
	protoc \
		-I $(PROTO_DIR) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
