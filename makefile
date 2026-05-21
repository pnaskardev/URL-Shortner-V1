.PHONY: auth core all

## VARIABLES
SERVICES_PATH=.

target: auth core

core:
	cd $(SERVICES_PATH)/core && air -c .air.toml

# auth:
# 	cd $(SERVICES_PATH)/auth && air -c .air.toml
