.ONESHELL: # Applies to every targets in the file!
.DEFAULT_GOAL := help

# Secure variables from .env
-include .env

# Project
# ................................................................................................ #
export PROJECT_NAME 			:= microgreens
export PROJECT_VERSION          :=
export PROJECT_ROOT             := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
export PROJECT_BIN              := $(PROJECT_ROOT)/bin
export PROJECT_DB_MIGRATION     := $(PROJECT_ROOT)/db_migration

# Postgres
# ................................................................................................ #
export DB_NAME                  ?= postgres
export POSTGRES_CONTAINER_NAME  := "$(PROJECT_NAME)-$(DB_NAME)"
export POSTGRES_HOST_PORT       := 5436
export SQL_MIGRATION            := $(PROJECT_DB_MIGRATION)/sql

# Text colors
# ................................................................................................ #
export RED                      := \033[0;31m
export BOLDRED                  := \033[1;31m
export GREEN                    := \033[0;32m
export YELLOW                   := \033[0;33m
export BLUE                     := \033[0;36m
export GRAY                     := \033[0;37m
export BOLD                     := \033[1m
export ENDCOLOR                 := \033[0m


##@ Postgres

.PHONY: postgres/run
postgres/run: ## Run container with postgres database
	@docker run \
	--name=$(POSTGRES_DB_NAME) \
	-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	--publish $(POSTGRES_HOST_PORT):5432 \
	-d \
	--rm \
	postgres

.PHONY: postgres/destroy
postgres/destroy: ## Stop and remove container with postgres database
	@docker ps -q --filter "name=$(POSTGRES_DB_NAME)" | grep -q . && docker stop -s 9 $(POSTGRES_DB_NAME) || true
	@docker container ls -a -q --filter "name=$(POSTGRES_DB_NAME)" | grep -q . && docker container rm -f $(POSTGRES_DB_NAME) || true

.PHONY: postgres/psql
postgres/psql: ## Attach to psql in postgres container
	@docker exec -it $(POSTGRES_DB_NAME) /bin/bash -c "psql -U postgres"

.PHONY: postgres/migrate/up
postgres/migrate/up: ## Start all postgres migration files with postfix "up"
	@migrate -path $(SQL_MIGRATION) -database "$(DB_NAME)://$(DB_NAME):${POSTGRES_PASSWORD}@localhost:$(POSTGRES_HOST_PORT)/$(DB_NAME)?sslmode=disable" up

.PHONY: postgres/migrate/down
postgres/migrate/down: ## Start all postgres migration files with postfix "down"
	@migrate -path $(SQL_MIGRATION) -database "$(DB_NAME)://$(DB_NAME):${POSTGRES_PASSWORD}@localhost:$(POSTGRES_HOST_PORT)/$(DB_NAME)?sslmode=disable" down

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: \"make <command>\"\n"} /^[a-zA-Z0-9_\/%-]+:.*?##/ { printf "  $(BLUE)%-29s$(ENDCOLOR) %s\n", $$1, $$2 } /^##@/ { printf "\n$(BOLD)%s$(ENDCOLOR)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""