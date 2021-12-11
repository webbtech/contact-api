include .env

// found 'watcher' at https://github.com/canthefason/go-watcher

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
	sam logs -n ContactAPIFunction --profile $(PROFILE) --stack-name $(STACK_NAME) --tail

tail-logs-trace:
	sam logs -n ContactAPIFunction --profile $(PROFILE) --stack-name $(STACK_NAME) --tail --include-traces