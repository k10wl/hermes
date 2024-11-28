APP_NAME=hermes
DEV_APP_NAME=hermes-dev
DEV_PORT=8124
SRC_DIR=.

.PHONY: build run install all dev-build dev-run dev-all watch clean test-web pre-test-app test-app test

build:
	@echo "Building prod version..."
	@go build -o ./bin/$(APP_NAME) $(SRC_DIR)

run:
	@echo "Running prod version..."
	@APP_NAME=$(APP_NAME) ./bin/$(APP_NAME)

install:
	@echo "Running binary installation..."
	@go install $(SRC_DIR)

all: build run

dev-build:
	@echo "Building dev version..."
	@go build -ldflags "-X github.com/k10wl/hermes/internal/settings.appName=$(DEV_APP_NAME) -X github.com/k10wl/hermes/internal/settings.DefaultPort=$(DEV_PORT)" -o ./bin/$(DEV_APP_NAME) $(SRC_DIR)

dev-run:
	@echo "Running dev version..."
	@APP_NAME=$(DEV_APP_NAME) ./bin/$(DEV_APP_NAME)

dev-all: dev-build dev-run

watch:
	@echo "Running dev version in watch mode..."
	@APP_NAME=$(DEV_APP_NAME) air -- serve --port $(DEV_PORT)

clean:
	@echo "Cleaning up..."
	@rm -rf ./bin ./tmp

MAX_TEST_DURATION=5s

test-web:
	@echo ">> Testing web..."
	@(cd internal/web && npm run test)
	@echo ">> Web passed test"

pre-test-app: 
	@echo ">> Testing helper functions..."
	@go test ./internal/test_helpers/... -v -timeout $(MAX_TEST_DURATION)
	@echo ">> Helpers passed test"

test-app:
	@echo ">> Testing app..."
	@export HERMES_TEST_HELPERS_SKIP=true; export HERMES_DB_DNS=:memory:; go test -ldflags "-X github.com/k10wl/hermes/internal/settings.appName=$(DEV_APP_NAME) -X github.com/k10wl/hermes/internal/settings.DefaultPort=$(DEV_PORT)" ./... -v -timeout $(MAX_TEST_DURATION)
	@echo ">> App passed test"

test: pre-test-app test-app test-web
	@echo "Tests passed"
