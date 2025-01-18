# Makefile for building TUI app

# Name of the output executable
APP_NAME = music-efx

# Path to the main Go file
MAIN_FILE = cmd/tui.go

# Output directory for build artifacts
OUTPUT_DIR = output

# Default target, builds for the current platform
all: build

# Build the project for the current platform
build:
	@echo "Building the application..."
	@mkdir -p $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(APP_NAME) $(MAIN_FILE)

# Cross-compile for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME) $(MAIN_FILE)

# Cross-compile for Windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME).exe $(MAIN_FILE)

# Build for all platforms (Linux, Windows, macOS)
build-all: build-linux build-windows

# Clean the build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(OUTPUT_DIR)

# Run the application (current platform)
run:
	@echo "Running the application..."
	./$(OUTPUT_DIR)/$(APP_NAME)

.PHONY: all build build-linux build-windows build-all clean run