include .env

.PHONY: build

build:
	sam build

local-api:
	sam local start-api

local-invoke:
	sam local invoke

dev-cloud:
	sam  sync --stack-name $(STACK_NAME) --profile $(PROFILE)

dev-cloud-watch:
	sam  sync --stack-name $(STACK_NAME) --watch --profile $(PROFILE)

tail-logs:
	sam logs -n HelloWorldFunction --profile pulpfree --stack-name mail-api --tail