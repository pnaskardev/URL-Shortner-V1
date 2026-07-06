.PHONY: redirector-service shortener-service core all build proto-gen proto-lint proto-breaking

## VARIABLES
SERVICES_PATH=.

target: redirector-service shortener-service core

## Compile all services against fresh contracts, without launching dev servers.
build: proto-gen
	cd $(SERVICES_PATH)/core && go build ./...
	cd $(SERVICES_PATH)/shortener-service && go build ./...
	cd $(SERVICES_PATH)/redirector-service && go build ./...

core: proto-gen
	cd $(SERVICES_PATH)/core && air -c .air.toml

shortener-service: proto-gen
	cd $(SERVICES_PATH)/shortener-service && air -c .air.toml

redirector-service: proto-gen
	cd $(SERVICES_PATH)/redirector-service && air -c .air.toml


## PROTOBUF / SCHEMA CONTRACTS
## Requires: buf (go install github.com/bufbuild/buf/cmd/buf@latest)
##       and protoc-gen-go (go install google.golang.org/protobuf/cmd/protoc-gen-go@latest)
proto-gen:
	cd $(SERVICES_PATH)/contracts && buf generate

proto-lint:
	cd $(SERVICES_PATH)/contracts && buf lint

proto-breaking:
	cd $(SERVICES_PATH)/contracts && buf breaking --against ".git#ref=origin/main,subdir=contracts"
