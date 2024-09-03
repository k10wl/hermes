APP_NAME=hermes
DEV_APP_NAME=hermes-dev
DEV_PORT=8124

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

dev-watch:
	@echo "Running dev version in watch mode..."
	@APP_NAME=$(DEV_APP_NAME) air

dev-all: dev-build dev-run

clean:
	@echo "Cleaning up..."
	@rm -f ./bin/$(APP_NAME) ./bin/$(DEV_APP_NAME)
	@echo "Done"

test:
	@echo "Running tests..."
	@go test -ldflags "-X github.com/k10wl/hermes/internal/settings.appName=$(DEV_APP_NAME) -X github.com/k10wl/hermes/internal/settings.DefaultPort=$(DEV_PORT) -X github.com/k10wl/hermes/internal/settings.DefaultDatabaseName=':memory:'" ./... -v
	@echo "Done"
