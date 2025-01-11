# Variables
APP_NAME = chat-app
BUILD_DIR = bin
SRC_DIR = ./cmd

.PHONY: all
all: build
	@echo "Detected OS: $(OS)"

.PHONY: build
build: 
	@echo "Building the app..."
	go build -o $(BUILD_DIR)/${APP_NAME} $(SRC_DIR)

.PHONY: run
run: build
	@echo "Running the app..."
	$(BUILD_DIR)/$(APP_NAME)

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rmdir /S /Q $(BUILD_DIR) || echo "Nothing to clean..."