BINARY_NAME=dockstat
BUILD_DIR=dist

all: clean build_all

build_all:
	@mkdir -p $(BUILD_DIR)
	# Linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	# macOS (Intel)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	# macOS (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	# FreeBSD
	GOOS=freebsd GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-amd64 .
	@echo "Builds completed in $(BUILD_DIR)/"

clean:
	@rm -rf $(BUILD_DIR)

.PHONY: all build_all clean