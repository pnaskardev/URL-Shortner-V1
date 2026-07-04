.PHONY: shortener-service core all proto-gen proto-lint proto-breaking

## VARIABLES
SERVICES_PATH=.

target: shortener-service core

core:
	cd $(SERVICES_PATH)/core && air -c .air.toml

shortener-service:
	cd $(SERVICES_PATH)/shortener-service && air -c .air.toml

## PROTOBUF / SCHEMA CONTRACTS
## Requires: buf (go install github.com/bufbuild/buf/cmd/buf@latest)
##       and protoc-gen-go (go install google.golang.org/protobuf/cmd/protoc-gen-go@latest)
proto-gen:
	cd $(SERVICES_PATH)/contracts && buf generate

proto-lint:
	cd $(SERVICES_PATH)/contracts && buf lint

proto-breaking:
	cd $(SERVICES_PATH)/contracts && buf breaking --against ".git#ref=origin/main,subdir=contracts"
