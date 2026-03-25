.PHONY: auth log order payment query gateway clean \
        dev dev-no-cache up down reset build infra proto logs \
        docker-auth docker-payment docker-log docker-order docker-query docker-gateway \
        down-auth down-payment down-log down-order down-query down-gateway \
        logs-auth logs-payment logs-log logs-order logs-query logs-gateway
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

gateway:
	cd $(CMD_PATH)/api-gateway && air

clean:
	rm -rf $(BIN_PATH)
	
# ─────────────────────────────────────
# DOCKER — run inside containers
# ─────────────────────────────────────

dev:
	docker compose up --build

dev-no-cache:
	docker compose build --no-cache
	docker compose up

up:
	docker compose up

down:
	docker compose down

reset:
	docker compose down -v

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

docker-order:
	docker compose up --build order-service

docker-query:
	docker compose up --build query-service

docker-gateway:
	docker compose up --build api-gateway

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

down-gateway:
	docker compose stop api-gateway

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

logs-gateway:
	docker compose logs -f api-gateway

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
