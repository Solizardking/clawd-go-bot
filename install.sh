#!/usr/bin/env bash
# ╔══════════════════════════════════════════════════════════════════════════════╗
# ║  ClawdBot — One-Shot Installer                                               ║
# ║  curl -fsSL https://raw.githubusercontent.com/Solizardking/clawdbot-go/main/install.sh | bash
# ╚══════════════════════════════════════════════════════════════════════════════╝

set -euo pipefail

REPO="https://github.com/Solizardking/clawdbot-go"
HUB_REPO="https://github.com/solizardking/solana-clawd"
TERMINAL_URL="https://cheshireterminal.ai"
INSTALL_API="https://zk.x402.wtf/api/install"
ZKROUTER_BASE="https://clawdrouter-zk.fly.dev/v1"
RPC_URL="https://zk.x402.wtf/api/solana/rpc-public"
INSTALL_DIR="${CLAWDBOT_INSTALL_DIR:-$HOME/.clawdbot}"
BIN_DIR="${CLAWDBOT_BIN_DIR:-$HOME/.local/bin}"

# ── Colours ──────────────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; CYAN='\033[0;36m'
YELLOW='\033[1;33m'; BOLD='\033[1m'; RESET='\033[0m'

info()    { echo -e "${CYAN}  ▶${RESET} $*"; }
success() { echo -e "${GREEN}  ✓${RESET} $*"; }
warn()    { echo -e "${YELLOW}  ⚠${RESET} $*"; }
die()     { echo -e "${RED}  ✗ ERROR:${RESET} $*" >&2; exit 1; }

# ── Banner ────────────────────────────────────────────────────────────────────
echo -e "${CYAN}"
cat << 'EOF'
    ██████╗██╗      █████╗ ██╗    ██╗██████╗ ██████╗  ██████╗ ████████╗
   ██╔════╝██║     ██╔══██╗██║    ██║██╔══██╗██╔══██╗██╔═══██╗╚══██╔══╝
   ██║     ██║     ███████║██║ █╗ ██║██║  ██║██████╔╝██║   ██║   ██║
   ██║     ██║     ██╔══██║██║███╗██║██║  ██║██╔══██╗██║   ██║   ██║
   ╚██████╗███████╗██║  ██║╚███╔███╔╝██████╔╝██████╔╝╚██████╔╝   ██║
    ╚═════╝╚══════╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚═════╝ ╚═════╝  ╚═════╝    ╚═╝
EOF
echo -e "${RESET}"
echo -e "${BOLD}  🦞 Sovereign Solana Trading Intelligence — Installer${RESET}"
echo -e "  Free AI via zkrouter · SolanaTracker RPC included"
echo

# ── Detect OS / Arch ──────────────────────────────────────────────────────────
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)        warn "Unsupported arch: $ARCH — will build from source" ;;
esac

info "Platform: ${OS}/${ARCH}"

# ── Check dependencies ────────────────────────────────────────────────────────
check_cmd() {
  command -v "$1" >/dev/null 2>&1
}

if ! check_cmd go; then
  echo
  warn "Go not found. Installing Go 1.22..."
  if [[ "$OS" == "darwin" ]]; then
    if check_cmd brew; then
      brew install go
    else
      die "Install Go from https://go.dev/dl/ then re-run this script"
    fi
  elif [[ "$OS" == "linux" ]]; then
    GOTAR="go1.22.5.linux-${ARCH}.tar.gz"
    curl -fsSL "https://go.dev/dl/${GOTAR}" -o "/tmp/${GOTAR}"
    sudo tar -C /usr/local -xzf "/tmp/${GOTAR}"
    export PATH="$PATH:/usr/local/go/bin"
  else
    die "Install Go from https://go.dev/dl/ then re-run this script"
  fi
fi

GO_VERSION="$(go version 2>&1 | awk '{print $3}')"
success "Go: ${GO_VERSION}"

# ── Check git ─────────────────────────────────────────────────────────────────
check_cmd git || die "git is required. Install it and re-run."

# ── Clone or update repo ──────────────────────────────────────────────────────
REPO_DIR="$INSTALL_DIR/src"
mkdir -p "$INSTALL_DIR"

if [[ -d "$REPO_DIR/.git" ]]; then
  info "Updating existing repo..."
  git -C "$REPO_DIR" pull --ff-only --quiet
else
  info "Cloning clawdbot-go..."
  git clone --depth=1 --quiet "$REPO" "$REPO_DIR"
fi
success "Source ready at $REPO_DIR"

# ── Build ──────────────────────────────────────────────────────────────────────
info "Building clawdbot binary..."
mkdir -p "$INSTALL_DIR/bin"
cd "$REPO_DIR"
go mod download -x 2>/dev/null | tail -3 || true
go build -ldflags="-s -w" -o "$INSTALL_DIR/bin/clawdbot" ./cmd/clawdbot/
success "Binary built: $INSTALL_DIR/bin/clawdbot"

# ── Install to PATH ────────────────────────────────────────────────────────────
mkdir -p "$BIN_DIR"
cp "$INSTALL_DIR/bin/clawdbot" "$BIN_DIR/clawdbot"
success "Installed to $BIN_DIR/clawdbot"

# Add to PATH if needed
if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
  SHELL_RC=""
  if [[ -f "$HOME/.zshrc" ]]; then SHELL_RC="$HOME/.zshrc"
  elif [[ -f "$HOME/.bashrc" ]]; then SHELL_RC="$HOME/.bashrc"
  elif [[ -f "$HOME/.profile" ]]; then SHELL_RC="$HOME/.profile"
  fi
  if [[ -n "$SHELL_RC" ]]; then
    echo "export PATH=\"\$PATH:$BIN_DIR\"" >> "$SHELL_RC"
    warn "Added $BIN_DIR to PATH in $SHELL_RC — restart your shell or run: export PATH=\"\$PATH:$BIN_DIR\""
  fi
fi

# ── Register install & get ID ─────────────────────────────────────────────────
info "Registering install with clawdrouter..."
INSTALL_ID=""
if check_cmd curl; then
  INSTALL_PAYLOAD="{\"os\":\"${OS}\",\"arch\":\"${ARCH}\",\"version\":\"$(cd "$REPO_DIR" && git rev-parse --short HEAD 2>/dev/null || echo 'unknown')\"}"
  INSTALL_RESP="$(curl -sf -X POST "$INSTALL_API" \
    -H "Content-Type: application/json" \
    -d "$INSTALL_PAYLOAD" 2>/dev/null || echo '{}')"
  INSTALL_ID="$(echo "$INSTALL_RESP" | grep -o '"installId":"[^"]*"' | cut -d'"' -f4 || echo '')"
fi

if [[ -z "$INSTALL_ID" ]]; then
  INSTALL_ID="local-$(date +%s)"
  warn "Could not reach install API — using local ID"
fi

success "Install ID: ${INSTALL_ID}"

# ── Write .env ────────────────────────────────────────────────────────────────
ENV_FILE="$INSTALL_DIR/.env"
if [[ ! -f "$ENV_FILE" ]]; then
  info "Writing default .env..."
  cat > "$ENV_FILE" << ENVEOF
# ════════════════════════════════════════════════════════════════════
# ClawdBot Environment — generated by installer
# Edit this file to add your own API keys
# ════════════════════════════════════════════════════════════════════

# ── Install identity ──────────────────────────────────────────────
CLAWDBOT_INSTALL_ID=${INSTALL_ID}

# ── Free AI via zkrouter (no key needed) ──────────────────────────
# Powered by clawdrouter-zk.fly.dev — $CLAWD sovereign AI routing
ZKROUTER_BASE_URL=${ZKROUTER_BASE}
ZKROUTER_API_KEY=clawdbot-free

# ── Solana RPC (SolanaTracker-backed, no key needed) ─────────────
SOLANA_RPC_URL=${RPC_URL}
HELIUS_RPC_URL=${RPC_URL}

# ── Optional: bring your own keys for higher limits ──────────────
# OPENROUTER_API_KEY=sk-or-...
# HELIUS_API_KEY=your-helius-key
# BIRDEYE_API_KEY=your-birdeye-key

# ── Wallet (required for live trading) ───────────────────────────
# WALLET_PRIVATE_KEY=your-base58-private-key

# ── Optional channels ─────────────────────────────────────────────
# TELEGRAM_BOT_TOKEN=your-telegram-token
ENVEOF
  success ".env written to $ENV_FILE"
else
  warn ".env already exists at $ENV_FILE — not overwriting"
fi

# Symlink config into home
if [[ ! -L "$HOME/.clawdbot" && "$INSTALL_DIR" != "$HOME/.clawdbot" ]]; then
  ln -sfn "$INSTALL_DIR" "$HOME/.clawdbot"
fi

# ── Done ──────────────────────────────────────────────────────────────────────
echo
echo -e "${GREEN}${BOLD}  ══════════════════════════════════════════${RESET}"
echo -e "${GREEN}${BOLD}  🦞 ClawdBot installed successfully!${RESET}"
echo -e "${GREEN}${BOLD}  ══════════════════════════════════════════${RESET}"
echo
echo -e "  ${BOLD}Get started:${RESET}"
echo -e "  ${CYAN}source ${ENV_FILE}${RESET}          # load env vars"
echo -e "  ${CYAN}clawdbot version${RESET}             # verify install"
echo -e "  ${CYAN}clawdbot agent${RESET}               # start AI REPL (free via zkrouter)"
echo -e "  ${CYAN}clawdbot ooda --sim${RESET}          # paper trading mode"
echo -e "  ${CYAN}clawdbot solana trending${RESET}     # top Solana tokens"
echo
echo -e "  ${BOLD}Edit your config:${RESET}  ${CYAN}nano ${ENV_FILE}${RESET}"
echo -e "  ${BOLD}Runtime repo:${RESET}      ${CYAN}${REPO}${RESET}"
echo -e "  ${BOLD}Ecosystem hub:${RESET}    ${CYAN}${HUB_REPO}${RESET}"
echo -e "  ${BOLD}Terminal:${RESET}         ${CYAN}${TERMINAL_URL}${RESET}"
echo -e "  ${BOLD}Free AI portal:${RESET}    ${CYAN}https://zk.x402.wtf${RESET}"
echo
echo -e "  🦞 \$CLAWD :: Droids Lead The Way"
echo
