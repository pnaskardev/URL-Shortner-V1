.PHONY: shortener-service core all

## VARIABLES
SERVICES_PATH=.

target: shortener-service core

core:
	cd $(SERVICES_PATH)/core && air -c .air.toml

shortener-service:
	cd $(SERVICES_PATH)/shortener-service && air -c .air.toml
