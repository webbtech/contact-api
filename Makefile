include .env

# found 'watcher' at https://github.com/canthefason/go-watcher
# that wasn't working as expected, so found 'fswatch' at: https://github.com/emcrisostomo/fswatch

.PHONY: build

build:
	sam build

local-api:
	sam local start-api --profile $(PROFILE)

local-invoke:
	sam local invoke --profile $(PROFILE)

dev-cloud:
	sam  sync --stack-name $(STACK_NAME) --profile $(PROFILE)

dev-cloud-watch:
	sam  sync --stack-name $(STACK_NAME) --watch --profile $(PROFILE)

tail-logs:
	sam logs -n ContactAPIFunction --profile $(PROFILE) --stack-name $(STACK_NAME) --tail

tail-logs-trace:
	sam logs -n ContactAPIFunction --profile $(PROFILE) --stack-name $(STACK_NAME) --tail --include-traces

validate:
	sam validate
	
watch:
	fswatch -o ./src | xargs -n1 -I{} sam build