APP_NAME=hermes
DEV_APP_NAME=hermes-dev

SRC_DIR=./cmd/hermes

.PHONY: build run all dev-build dev-run dev-watch dev-all clean

build:
	@echo "Building prod version..."
	@go build -ldflags "-X 'main.appName=$(APP_NAME)'" -o ./bin/$(APP_NAME) $(SRC_DIR)

run:
	@echo "Running prod version..."
	@APP_NAME=$(APP_NAME) ./bin/$(APP_NAME)

all: build run

dev-build:
	@echo "Building dev version..."
	@go build -ldflags "-X 'main.appName=$(APP_NAME)'" -o ./bin/$(DEV_APP_NAME) $(SRC_DIR)

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
