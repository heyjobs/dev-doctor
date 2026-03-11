.PHONY: build install uninstall clean test

# Build variables
BINARY_NAME=dev-doctor
BUILD_DIR=.
INSTALL_DIR=$(HOME)/bin

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/dev-doctor

# Install the binary to ~/bin (or /usr/local/bin with sudo)
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Installed $(BINARY_NAME) to $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo ""
	@echo "Make sure $(INSTALL_DIR) is in your PATH."
	@echo "If not, add this to your ~/.zshrc or ~/.bashrc:"
	@echo "  export PATH=\"\$$HOME/bin:\$$PATH\""

# Install to /usr/local/bin (requires sudo)
install-global: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Installed $(BINARY_NAME) to /usr/local/bin/$(BINARY_NAME)"

# Uninstall from ~/bin
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Uninstalled $(BINARY_NAME)"

# Uninstall from /usr/local/bin
uninstall-global:
	@echo "Uninstalling $(BINARY_NAME) from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Uninstalled $(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "✓ Cleaned"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Show help
help:
	@echo "dev-doctor Makefile commands:"
	@echo ""
	@echo "  make build           - Build the binary"
	@echo "  make install         - Install to ~/bin (user-level)"
	@echo "  make install-global  - Install to /usr/local/bin (system-level, requires sudo)"
	@echo "  make uninstall       - Remove from ~/bin"
	@echo "  make uninstall-global- Remove from /usr/local/bin"
	@echo "  make clean           - Remove build artifacts"
	@echo "  make test            - Run tests"
	@echo "  make help            - Show this help message"
