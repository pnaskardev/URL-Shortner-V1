.PHONY: analytics-dashboard-service redirector-service shortener-service core build proto-gen proto-lint proto-breaking dev infra-up infra-down infra-clean infra-logs

## VARIABLES
SERVICES_PATH=.
COMPOSE = docker compose

target: analytics-dashboard-service redirector-service shortener-service core

## DEV — one command: infra up (healthy) then all services hot-reloading
dev: infra-up
	$(MAKE) -j4 target

## INFRA — Postgres, Redis, RabbitMQ in Docker
infra-up:                 ## boot infra, block until healthy
	$(COMPOSE) up -d --wait
infra-down:               ## stop containers, keep data
	$(COMPOSE) down
infra-clean:              ## stop + wipe volumes (fresh DBs next up)
	$(COMPOSE) down -v
infra-logs:
	$(COMPOSE) logs -f

## Compile all services against fresh contracts, without launching dev servers.
build: proto-gen
	cd $(SERVICES_PATH)/core && go build ./...
	cd $(SERVICES_PATH)/shortener-service && go build ./...
	cd $(SERVICES_PATH)/redirector-service && go build ./...
	cd $(SERVICES_PATH)/analytics-dashboard-service && go build ./...

core: proto-gen
	cd $(SERVICES_PATH)/core && air -c .air.toml

shortener-service: proto-gen
	cd $(SERVICES_PATH)/shortener-service && air -c .air.toml

redirector-service: proto-gen
	cd $(SERVICES_PATH)/redirector-service && air -c .air.toml

analytics-dashboard-service: proto-gen
	cd $(SERVICES_PATH)/analytics-dashboard-service && air -c .air.toml

## PROTOBUF / SCHEMA CONTRACTS
## Requires: buf (go install github.com/bufbuild/buf/cmd/buf@latest)
##       and protoc-gen-go (go install google.golang.org/protobuf/cmd/protoc-gen-go@latest)
proto-gen:
	cd $(SERVICES_PATH)/contracts && buf generate

proto-lint:
	cd $(SERVICES_PATH)/contracts && buf lint

proto-breaking:
	cd $(SERVICES_PATH)/contracts && buf breaking --against ".git#ref=origin/main,subdir=contracts"
