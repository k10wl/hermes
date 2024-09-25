APP_NAME=hermes
DEV_APP_NAME=hermes-dev
DEV_PORT=8127

SRC_DIR=.

.PHONY: build run all dev-build dev-run dev-watch dev-all clean

build:
	@echo "Building prod version..."
	@go build -o ./bin/$(APP_NAME) $(SRC_DIR)
	@echo "Done"

run:
	@echo "Running prod version..."
	@APP_NAME=$(APP_NAME) ./bin/$(APP_NAME)

install:
	@echo "Running binary instalation..."
	@go install $(SRC_DIR)
	@echo "Done"

all: build run

dev-build:
	@echo "Building dev version..."
	@go build -ldflags "-X github.com/k10wl/hermes/internal/settings.appName=$(DEV_APP_NAME) -X github.com/k10wl/hermes/internal/settings.DefaultPort=$(DEV_PORT)" -o ./bin/$(DEV_APP_NAME) $(SRC_DIR)
	@echo "Done"

dev-run:
	@echo "Running dev version..."
	@APP_NAME=$(DEV_APP_NAME) ./bin/$(DEV_APP_NAME)

dev-all: dev-build dev-run

serve-watch:
	@echo "Running dev version in watch mode..."
	@APP_NAME=$(DEV_APP_NAME) air -- serve --port 8124


clean:
	@echo "Cleaning up..."
	@rm -f ./bin/$(APP_NAME) ./bin/$(DEV_APP_NAME)
	@echo "Done"

pre-test: 
	@echo ">> Testing helper functions..."
	@go test  ./internal/test_helpers/... -v
	@echo ">> Finished helpers testing"
test-app:
	@echo ">> Testing app..."
	@export HERMES_TEST_HELPERS_SKIP=true; export HERMES_DB_DNS=:memory:; go test -ldflags "-X github.com/k10wl/hermes/internal/settings.appName=$(DEV_APP_NAME) -X github.com/k10wl/hermes/internal/settings.DefaultPort=$(DEV_PORT)" ./... -v
	@echo ">> Finished app testing"
test:
	@echo "> Starting testing"
	@make pre-test && make test-app
	@echo "> Done"
