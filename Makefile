.ONESHELL: # Applies to every targets in the file!
.DEFAULT_GOAL := help

# Secure variables from .env
-include .env

# Project
# ................................................................................................ #
export PROJECT_NAME 			:= microgreens
export PROJECT_VERSION          := 0.0.1
export PROJECT_ROOT             := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
export PROJECT_BIN              := $(PROJECT_ROOT)/bin
export PROJECT_CMD              := $(PROJECT_ROOT)/cmd
export PROJECT_BUILD            := $(PROJECT_ROOT)/build
export PROJECT_DB_MIGRATION     := $(PROJECT_ROOT)/db_migration
export GO                       ?= go

# Migration
# ................................................................................................ #
export MIGRATION_NAME           ?=

# Postgres
# ................................................................................................ #
export SQL_MIGRATION            := $(PROJECT_DB_MIGRATION)/sql
export DB_NAME                  ?= postgres
export POSTGRES_CONTAINER_NAME  := "$(PROJECT_NAME)-$(DB_NAME)"
export POSTGRES_HOSTNAME        := localhost
export POSTGRES_HOST_PORT       := 5436
export POSTGRES_SSL_MODE        := disable
ifneq (${MG_CFG_POSTGRESS.PASSWORD},)
	export POSTGRES_PASSWORD    := ${MG_CFG_POSTGRES.PASSWORD}
else
	export POSTGRES_PASSWORD    := ${POSTGRES.PASSWORD}
endif

# Buildmode
# ................................................................................................ #
export DEBUG_BUILDMODE          := debug
export RELEASE_BUILDMODE        := release

# CI/CD
# ................................................................................................ #
export CI_ARTIFACTS_DIR         ?=

# Development containers
# ................................................................................................ #
export DEV_CONTAINER_NAME       := dev.$(PROJECT_NAME).v$(PROJECT_VERSION)
export DEV_IMAGE_NAME           := $(DEV_CONTAINER_NAME)

# Golang variables
# ................................................................................................ #
export GOPROXY                  ?= https://proxy.golang.org
export GOSUMDB                  ?=

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


# Inputs for the '*/debug' commands
# ................................................................................................ #
args = $(foreach a,$($(subst -,_,$1)_args),$(if $(value $a),$a="$($a)"))

##@ App

.PHONY: app/debug
app/debug: ## Run app through delve debugger
	@dlv --listen=:2375 --headless=true --api-version=2 --accept-multiclient debug $(PROJECT_CMD)/main.go -- $(args)

##@ Postgres

.PHONY: postgres/run
postgres/run: postgres/destroy ## Run container with postgres database
	@docker run \
	--name=$(POSTGRES_CONTAINER_NAME) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	--publish $(POSTGRES_HOST_PORT):5432 \
	-d \
	postgres

.PHONY: postgres/destroy
postgres/destroy: ## Stop and remove container with postgres database
	@docker ps -q --filter "name=$(POSTGRES_CONTAINER_NAME)" | grep -q . && docker stop -s 9 $(POSTGRES_CONTAINER_NAME) || true
	@docker container ls -a -q --filter "name=$(POSTGRES_CONTAINER_NAME)" | grep -q . && docker container rm -f $(POSTGRES_CONTAINER_NAME) || true

.PHONY: postgres/psql
postgres/psql: ## Attach to psql in postgres container
	@docker exec -it $(POSTGRES_CONTAINER_NAME) /bin/bash -c "psql -U postgres"

.PHONY: postgres/migrate/create
postgres/migrate/create: ## Create postgres migration files with custom $MIGRATION_NAME
	@migrate create -ext sql -dir ./db_migration/sql -seq $(MIGRATION_NAME)

.PHONY: postgres/migrate/up
postgres/migrate/up: ## Start all postgres migration files with postfix "up"
	@migrate -path $(SQL_MIGRATION) -database "$(DB_NAME)://$(DB_NAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOSTNAME):$(POSTGRES_HOST_PORT)/$(DB_NAME)?sslmode=$(POSTGRES_SSL_MODE)" up

.PHONY: postgres/migrate/down
postgres/migrate/down: ## Start all postgres migration files with postfix "down"
	@migrate -path $(SQL_MIGRATION) -database "$(DB_NAME)://$(DB_NAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOSTNAME):$(POSTGRES_HOST_PORT)/$(DB_NAME)?sslmode=$(POSTGRES_SSL_MODE)" down 1

.PHONY: postgres/migrate/drop
postgres/migrate/drop: ## Drop all postgres migration files
	@migrate -path $(SQL_MIGRATION) -database "$(DB_NAME)://$(DB_NAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOSTNAME):$(POSTGRES_HOST_PORT)/$(DB_NAME)?sslmode=$(POSTGRES_SSL_MODE)" drop

##@ Maintenance

.PHONY: dependencies
dependencies: dependencies/go ## Alias to 'dependencies/go'

.PHONY: dependencies/go
dependencies/go: ## Install all Go dependencies and vendor it
	@$(GO) mod tidy
	@$(GO) mod vendor
	@$(GO) mod verify

.PHONY: dependencies/go/restore
dependencies/go/restore: ## Restore 'vendor' for each component
	@rm -rf vendor && git restore vendor

##@ Development container

.PHONY: dev
dev: dev/destroy dev/create dev/up dev/dependencies  ## Create dev container and perform all preparing stuff

.PHONY: dev/%
dev/%: ## Run make-recipe inside container in interactive mode
	@echo -e "$(YELLOW)[DEBUG] cmd to exec inside $(DEV_CONTAINER_NAME) container: $(MAKE) $* $(ENDCOLOR)"
	@docker exec -it $(DEV_CONTAINER_NAME) $(MAKE) $*

.PHONY: dev/create
dev/create: ## Service dev-container: create container
	@docker build \
		--tag $(DEV_IMAGE_NAME) \
		--build-arg USER_ID=$(shell id -u) \
		--build-arg USER_NAME=$(USER) \
		--build-arg GROUP_ID=$(shell id -g) \
		--build-arg GROUP_NAME=$(USER) \
		--file ${PROJECT_BUILD}/devcontainer.Dockerfile .

	@docker container create \
		--name $(DEV_CONTAINER_NAME) \
		--volume="$(CURDIR):/home/$(USER)/workdir" \
		--publish 2375:2375 \
		--workdir /home/$(USER)/workdir \
		--user "$(shell id -u):$(shell id -g)" \
		$(DEV_IMAGE_NAME)

.PHONY: dev/down
dev/down: ## Service dev-container: stop container
	@docker ps -q --filter "name=$(DEV_CONTAINER_NAME)" | grep -q . && docker stop -s 9 $(DEV_CONTAINER_NAME) || true

.PHONY: dev/up
dev/up: ## Service dev-container: start container
	@-docker start $(DEV_CONTAINER_NAME)

.PHONY: dev/shell
dev/shell: ## Service dev-container: attach to container terminal for manual entering tasks
	@docker exec -it $(DEV_CONTAINER_NAME) bash

.PHONY: dev/destroy
dev/destroy: dev/down ## Service dev-container: destroy container and its image
	@docker container ls -a -q --filter "name=$(DEV_CONTAINER_NAME)" | grep -q . && docker container rm -f $(DEV_CONTAINER_NAME) || true
	@docker image ls -q --filter "reference=$(DEV_IMAGE_NAME)" | grep -q . && docker image rm -f $(DEV_IMAGE_NAME) || true

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: \"make <command>\"\n"} /^[a-zA-Z0-9_\/%-]+:.*?##/ { printf "  $(BLUE)%-29s$(ENDCOLOR) %s\n", $$1, $$2 } /^##@/ { printf "\n$(BOLD)%s$(ENDCOLOR)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""