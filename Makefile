# ===============================
# Simple Go Project Makefile
# ===============================

# App name
APP_NAME = chat

# Go command
GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run

# Binary output
BIN_DIR = ./bin
BIN_FILE = $(BIN_DIR)/$(APP_NAME)

# -------------------------------
# Default target: build
.PHONY: all
all: build

# -------------------------------
# Run locally without building binary
.PHONY: run
run:
	$(GORUN) main.go

# -------------------------------
# Build binary
.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_FILE) /main.go
	@echo "✅ Build complete: $(BIN_FILE)"

# -------------------------------
# Build & run binary
.PHONY: start
start: build
	$(BIN_FILE)