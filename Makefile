# ── ClawdBot Go :: Makefile ─────────────────────────────────────────────
# Build targets for x86_64, ARM64 (NVIDIA Orin Nano), and Arduino bridge
#
# Usage:
#   make build          Build for current platform
#   make orin           Cross-compile for NVIDIA Orin Nano (linux/arm64)
#   make tui            Build TUI launcher
#   make all            Build all targets
#   make docker         Build Docker image
#   make clean          Remove build artifacts
#   make install        Install to /usr/local/bin
#   make test           Run tests
#   make scan-i2c       Scan I2C bus for Modulino® sensors
# ──────────────────────────────────────────────────────────────────────

VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILDTIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GOVERSION := $(shell go version | cut -d' ' -f3)

MODULE    := github.com/8bitlabs/clawdbot
# The Go import path stays stable for v1 compatibility even though the public
# runtime repo is https://github.com/Solizardking/clawdbot-go and the wider hub
# is https://github.com/solizardking/solana-clawd.
PKG_VER   := $(MODULE)/pkg/config

LDFLAGS   := -s -w \
  -X $(PKG_VER).Version=$(VERSION) \
  -X $(PKG_VER).GitCommit=$(COMMIT) \
  -X $(PKG_VER).BuildTime=$(BUILDTIME) \
  -X $(PKG_VER).GoVersion=$(GOVERSION)

# Shared build settings
GO        := go
GOBUILD   := $(GO) build -trimpath -ldflags "$(LDFLAGS)"
GOTEST    := $(GO) test -v -race
GOFILES   := $(shell find cmd pkg web/backend -type f -name '*.go')

# Output directories
BUILD_DIR := ./build
BIN_CLI   := $(BUILD_DIR)/clawdbot
BIN_TUI   := $(BUILD_DIR)/clawdbot-tui

.PHONY: all build orin rpi riscv macos cross tui web docker docker-orin clean install test lint verify audit release-check deps scan-i2c help

# ── Default ───────────────────────────────────────────────────────────

all: build tui

# ── Build for current platform ────────────────────────────────────────

build:
	@echo "🦞 Building ClawdBot CLI..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BIN_CLI) ./cmd/clawdbot
	@echo "✓ $(BIN_CLI) built ($(shell file $(BIN_CLI) | cut -d: -f2))"
	@ls -lh $(BIN_CLI)

tui:
	@echo "🦞 Building ClawdBot TUI Launcher..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BIN_TUI) ./cmd/clawdbot-tui
	@echo "✓ $(BIN_TUI) built"
	@ls -lh $(BIN_TUI)

web:
	@echo "🦞 Building ClawdBot Web Console..."
	$(MAKE) -C web all

# ── NVIDIA Orin Nano (Linux ARM64) ────────────────────────────────────
# The Orin Nano runs Ubuntu 22.04 aarch64 (Jetson Linux / JetPack 6.x)
# CGO enabled for I2C syscalls (hardware/modulino.go)

orin:
	@echo "🦞 Cross-compiling for NVIDIA Orin Nano (linux/arm64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		$(GOBUILD) -o $(BUILD_DIR)/clawdbot-orin ./cmd/clawdbot
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		$(GOBUILD) -o $(BUILD_DIR)/clawdbot-tui-orin ./cmd/clawdbot-tui
	@echo "✓ Orin Nano binaries:"
	@ls -lh $(BUILD_DIR)/clawdbot-orin $(BUILD_DIR)/clawdbot-tui-orin
	@echo ""
	@echo "📦 Deploy to Orin Nano:"
	@echo "  scp $(BUILD_DIR)/clawdbot-orin user@orin-nano:~/clawdbot"
	@echo "  scp $(BUILD_DIR)/clawdbot-tui-orin user@orin-nano:~/clawdbot-tui"

# ── Raspberry Pi / Generic ARM ────────────────────────────────────────

rpi:
	@echo "🦞 Cross-compiling for Raspberry Pi (linux/arm64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		$(GOBUILD) -o $(BUILD_DIR)/clawdbot-rpi ./cmd/clawdbot
	@echo "✓ $(BUILD_DIR)/clawdbot-rpi built"

# ── RISC-V ────────────────────────────────────────────────────────────

riscv:
	@echo "🦞 Cross-compiling for RISC-V (linux/riscv64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 \
		$(GOBUILD) -o $(BUILD_DIR)/clawdbot-riscv ./cmd/clawdbot
	@echo "✓ $(BUILD_DIR)/clawdbot-riscv built"

# ── macOS (Apple Silicon) ─────────────────────────────────────────────

macos:
	@echo "🦞 Building for macOS (darwin/arm64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 \
		$(GOBUILD) -o $(BUILD_DIR)/clawdbot-macos ./cmd/clawdbot
	@echo "✓ $(BUILD_DIR)/clawdbot-macos built"

# ── All platforms ─────────────────────────────────────────────────────

cross: build orin rpi riscv macos
	@echo ""
	@echo "🦞 All cross-compilation complete:"
	@ls -lh $(BUILD_DIR)/

# ── Docker ────────────────────────────────────────────────────────────

docker:
	@echo "🐳 Building Docker image..."
	docker build -t clawdbot:$(VERSION) -t clawdbot:latest .
	@echo "✓ Docker image built: clawdbot:$(VERSION)"

docker-orin:
	@echo "🐳 Building Docker image for Orin Nano (linux/arm64)..."
	docker buildx build --platform linux/arm64 \
		-t clawdbot:$(VERSION)-orin \
		-t clawdbot:latest-orin .
	@echo "✓ Docker image built: clawdbot:$(VERSION)-orin"

# ── Install ───────────────────────────────────────────────────────────

install: build tui
	@echo "📦 Installing to /usr/local/bin..."
	install -m 755 $(BIN_CLI) /usr/local/bin/clawdbot
	install -m 755 $(BIN_TUI) /usr/local/bin/clawdbot-tui
	@echo "✓ Installed clawdbot and clawdbot-tui"

# ── Test ──────────────────────────────────────────────────────────────

test:
	@echo "🧪 Running tests..."
	$(GOTEST) ./...

lint:
	@echo "🔍 Running linter..."
	golangci-lint run ./...

verify:
	@echo "🔎 Verifying Go formatting..."
	@files="$$(gofmt -l $(GOFILES))"; \
	if [ -n "$$files" ]; then \
		echo "gofmt needed:"; \
		echo "$$files"; \
		exit 1; \
	fi
	@echo "🔎 Running go vet..."
	$(GO) vet ./...
	@echo "🧪 Running race tests..."
	$(GOTEST) ./...
	@echo "🦞 Building release entrypoints..."
	$(GOBUILD) -o /tmp/clawdbot-verify ./cmd/clawdbot
	$(GOBUILD) -o /tmp/clawdbot-tui-verify ./cmd/clawdbot-tui

audit: verify
	@echo "🔒 Running vulnerability scan when govulncheck is installed..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not installed; install with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi

release-check: verify
	@echo "📦 Checking tracked generated artifacts..."
	@bad="$$(git ls-files -- .cache build dist clawdbot ':(glob)**/.next/**' ':(glob)**/target/**' ':(glob)**/*.tsbuildinfo')"; \
	if [ -n "$$bad" ]; then \
		echo "generated artifacts are tracked and must be removed from git:"; \
		echo "$$bad"; \
		exit 1; \
	fi

# ── Dependencies ──────────────────────────────────────────────────────

deps:
	@echo "📦 Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# ── Hardware ──────────────────────────────────────────────────────────

scan-i2c:
	@echo "🔍 Scanning I2C bus for Modulino® sensors..."
	@i2cdetect -y 1 2>/dev/null || echo "i2cdetect not available — install i2c-tools"
	@echo ""
	@echo "Expected Modulino® addresses:"
	@echo "  0x29 — Distance (VL53L4CD)"
	@echo "  0x3C — Buzzer (PKLCS1212E)"
	@echo "  0x44 — Thermo (HS3003)"
	@echo "  0x6A — Movement (LSM6DSOX)"
	@echo "  0x6C — Pixels (LC8822)"
	@echo "  0x76 — Knob (PEC11J)"
	@echo "  0x7C — Buttons (3x push)"

# ── Clean ─────────────────────────────────────────────────────────────

clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf .cache
	@echo "✓ Clean"

# ── Help ──────────────────────────────────────────────────────────────

help:
	@echo "ClawdBot Go — Makefile targets:"
	@echo ""
	@echo "  build       Build for current platform"
	@echo "  tui         Build TUI launcher"
	@echo "  web         Build web frontend + backend"
	@echo "  all         Build CLI + TUI"
	@echo "  orin        Cross-compile for NVIDIA Orin Nano (linux/arm64)"
	@echo "  rpi         Cross-compile for Raspberry Pi (linux/arm64)"
	@echo "  riscv       Cross-compile for RISC-V (linux/riscv64)"
	@echo "  macos       Build for macOS Apple Silicon"
	@echo "  cross       All cross-compilation targets"
	@echo "  docker      Build Docker image"
	@echo "  docker-orin Build Docker for Orin Nano"
	@echo "  install     Install to /usr/local/bin"
	@echo "  test        Run tests"
	@echo "  lint        Run linter"
	@echo "  verify      Format check + vet + race tests + builds"
	@echo "  audit       verify + govulncheck when installed"
	@echo "  release-check verify + tracked artifact guard"
	@echo "  deps        Download dependencies"
	@echo "  scan-i2c    Scan for Modulino sensors"
	@echo "  clean       Remove build artifacts"
	@echo ""
	@echo "  Version: $(VERSION) | Commit: $(COMMIT)"
