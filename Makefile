# ===============================
# Simple Go Project Makefile
# ===============================

# App name
APP_NAME = chat

# Go command
GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run

# Source directory
SRC_DIR = ./cmd

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
	$(GORUN) $(SRC_DIR)/main.go

# -------------------------------
# Build binary
.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_FILE) $(SRC_DIR)/main.go
	@echo "✅ Build complete: $(BIN_FILE)"

# -------------------------------
# Build & run binary
.PHONY: start
start: build
	$(BIN_FILE)