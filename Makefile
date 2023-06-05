#const
.PHONY: help
GO_VERSION ?= 1.17
GO_SWAGGER_TAG = v0.28.0
GOFUMPT_VERSION ?= v0.1.0
GOLANGCI_VERSION ?= v1.45.2


-include .env
-include .env.local

APPS ?= wallet stat


SWAGGER_YAML_FILE ?= api/swagger.yml
CLIENT_PKG_SUFFIX ?= client

DEFINITIONS_YAML := api/wallet-swagger.gen.yaml
SWAGGER_YAML_TEMPLATE := api/core-swagger.template.yml
SWAGGER_WITH_DEFINITIONS_YAML := api/core-swagger.yml

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


lint: ##  Linting code
	docker run --rm -v $(PWD):$(PWD) -w $(PWD) -u `id -u $(USER)` -e GOLANGCI_LINT_CACHE=/tmp/.cache -e GOCACHE=/tmp/.cache golangci/golangci-lint:$(GOLANGCI_VERSION) golangci-lint run -v --fix
	make lint-swagger


lint-all: ##
	docker run --rm -v $(PWD):$(PWD) -w $(PWD) -u `id -u $(USER)` -e GOLANGCI_LINT_CACHE=/tmp/.cache -e GOCACHE=/tmp/.cache golangci/golangci-lint:$(GOLANGCI_VERSION) golangci-lint run -v --fix
	make lint-swaggers

lint-ci: ##
	wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b .bin $(GOLANGCI_VERSION)
	./.bin/golangci-lint run -v --fix

lint-swagger: ##
	docker run --rm -v $(PWD):$(PWD) -w $(PWD) openapitools/openapi-generator-cli validate -i api/$(SERVICE_NAME)-swagger.yml --recommend

lint-swaggers: ##
	@$(foreach SERVICE_NAME, $(APPS), make lint-swagger SERVICE_NAME=$(SERVICE_NAME) &&) echo "==> $@ completed"

up: ##
	docker-compose -f docker-compose.yml pull
	docker-compose -f docker-compose.yml up

up-wallet: ##
	docker-compose -f docker-compose.yml up wallet

up-stat: ##
	docker-compose -f docker-compose.yml up stat

up-build: ##
	docker-compose up --build --remove-orphans --force-recreate

ps: ## Show running processes
	docker-compose ps

clean: ##
	docker-compose down
	rm -rf .data/*

test: ##
	docker-compose -f docker-compose.unit.yml pull
	docker-compose -f docker-compose.unit.yml up -d
# OSX does not support time quantifiers -eg sleep 20s
#	sleep 20
	docker-compose -f docker-compose.unit.yml exec -T test sh -c 'go test --count=1 ./... -v'
	docker-compose -f docker-compose.unit.yml down


gen-server-native: ##
ifndef SERVICE_NAME
	$(error SERVICE_NAME is not set)
endif
	echo "service=${SERVICE_NAME}"
	rm -rf $(PWD)/internal/app/$(SERVICE_NAME)/server/*
	mkdir -p $(PWD)/internal/app/$(SERVICE_NAME)/server
	#generate
	swagger generate server \
		-f $(PWD)/api/$(SERVICE_NAME)-swagger.yml \
		-A $(SERVICE_NAME) \
		--template-dir $(PWD)/tools/swagger-templates \
		--exclude-main \
		-P models.Principal \
		-t $(PWD)/internal/app/$(SERVICE_NAME)/server
	@docker run --rm --platform linux/amd64 -v $(PWD):$(PWD) -w $(PWD) plavreshin/goimports:latest -e -w internal/app/$(SERVICE_NAME)/server

gen-server: ##
ifndef SERVICE_NAME
	$(error SERVICE_NAME is not set)
endif
	echo "service=${SERVICE_NAME}"
	rm -rf $(PWD)/internal/app/$(SERVICE_NAME)/server/*
	mkdir -p $(PWD)/internal/app/$(SERVICE_NAME)/server
	#generate
	docker run --rm -it -u `id -u $(USER)` \
		-v $(PWD):$(PWD) \
		-w $(PWD) quay.io/goswagger/swagger:$(GO_SWAGGER_TAG) generate server \
		-f $(PWD)/api/$(SERVICE_NAME)-swagger.yml \
		-A $(SERVICE_NAME) \
		--template-dir $(PWD)/tools/swagger-templates \
		--exclude-main \
		-P models.Principal \
		-t $(PWD)/internal/app/$(SERVICE_NAME)/server
	@docker run --rm  -v $(PWD):$(PWD) -w $(PWD) plavreshin/goimports:latest -e -w internal/app/$(SERVICE_NAME)/server


gen-servers:
	@$(foreach SERVICE_NAME, $(APPS), make gen-server SERVICE_NAME=$(SERVICE_NAME) &&) echo "==> $@ completed"

goimports: ##
	@docker run --rm  -v $(PWD):$(PWD) -w $(PWD) plavreshin/goimports:latest -e -w internal/app/$(SERVICE_NAME)/server

goimports-all: ##
	@docker run --rm  -v $(PWD):$(PWD) -w $(PWD) plavreshin/goimports:latest -e -w internal/

gen-clients: ##
	@$(foreach SERVICE_NAME, $(APPS), make gen-client SERVICE_NAME=$(SERVICE_NAME) &&) echo "==> $@ completed"

gen-client: ##
ifndef SERVICE_NAME
	$(error SERVICE_NAME is not set)
endif
	echo "service=${SERVICE_NAME}"
	rm -rf $(PWD)/pkg/client/$(SERVICE_NAME)*
	docker run --rm -it -u `id -u $(USER)` \
		-v $(PWD):$(PWD) \
		-w $(PWD) quay.io/goswagger/swagger:$(GO_SWAGGER_TAG) generate client \
		-f $(PWD)/api/$(SERVICE_NAME)-swagger.yml \
		-A $(SERVICE_NAME) \
		--template-dir $(PWD)/tools/swagger-templates \
		-c $(SERVICE_NAME)$(CLIENT_PKG_SUFFIX) \
		-m $(SERVICE_NAME)$(MODELS_PKG_SUFFIX) \
		-P models.Principal \
		-t $(PWD)/pkg/client

gen-client-native: ##
ifndef SERVICE_NAME
	$(error SERVICE_NAME is not set)
endif
	echo "service=${SERVICE_NAME}"
	rm -rf $(PWD)/pkg/client/$(SERVICE_NAME)*
	swagger generate client \
		-f $(PWD)/api/$(SERVICE_NAME)-swagger.yml \
		-A $(SERVICE_NAME) \
		--template-dir $(PWD)/tools/swagger-templates \
		-c $(SERVICE_NAME)$(CLIENT_PKG_SUFFIX) \
		-m $(SERVICE_NAME)$(MODELS_PKG_SUFFIX) \
		-P models.Principal \
		-t $(PWD)/pkg/client


gen-handler-declarations: ##
	@$(foreach SERVICE_NAME, $(APPS), make gen-handler-declaration SERVICE_NAME=$(SERVICE_NAME) &&) echo "==> $@ completed"

gen-handler-declaration: ##
ifndef SERVICE_NAME
	$(error SERVICE_NAME is not set)
endif
	SERVICE_NAME=$(SERVICE_NAME) go run scripts/generate-handler-declarations/main.go

gen-handler-implementations: ##
	@$(foreach SERVICE_NAME, $(APPS), make gen-handler-implementation SERVICE_NAME=$(SERVICE_NAME) &&) echo "==> $@ completed"

gen-handler-implementation: ##
ifndef SERVICE_NAME
	$(error SERVICE_NAME is not set)
endif
	SERVICE_NAME=$(SERVICE_NAME) go run scripts/generate-handler-implementations/main.go

logs: ##
	@docker-compose logs -f

gen-app-messages:
	go run scripts/generate-app-messages/main.go




