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
SOURCE_MODE="${CLAWDBOT_SOURCE_MODE:-archive}"
REF="${CLAWDBOT_REF:-main}"
CORE_AI_REPO="${CLAWDBOT_CORE_AI_REPO:-https://github.com/Solizardking/core-ai}"
CORE_AI_REF="${CLAWDBOT_CORE_AI_REF:-clawd-stack-integration}"
CORE_AI_DIR="${CLAWDBOT_CORE_AI_DIR:-$INSTALL_DIR/core-ai}"
CORE_AI_MCP_CONFIG="${CLAWDBOT_CORE_AI_MCP_CONFIG:-$INSTALL_DIR/core-ai.mcp.json}"
INSTALL_CORE_AI="${CLAWDBOT_INSTALL_CORE_AI:-0}"

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

github_archive_url() {
  local repo_url="$1"
  local ref="$2"
  printf "%s/archive/%s.tar.gz" "${repo_url%.git}" "$ref"
}

install_source_archive() {
  local repo_url="$1"
  local ref="$2"
  local dest="$3"
  local label="$4"
  local tmp archive root

  check_cmd curl || die "curl is required to install $label from an archive"
  check_cmd tar || die "tar is required to install $label from an archive"

  tmp="$(mktemp -d)"
  archive="$tmp/source.tar.gz"
  curl -fsSL "$(github_archive_url "$repo_url" "$ref")" -o "$archive"
  tar -xzf "$archive" -C "$tmp"
  root="$(find "$tmp" -mindepth 1 -maxdepth 1 -type d | head -1)"
  if [[ -z "$root" ]]; then
    rm -rf "$tmp"
    die "Could not unpack $label archive"
  fi

  rm -rf "$dest"
  mkdir -p "$(dirname "$dest")"
  mv "$root" "$dest"
  rm -rf "$tmp"
}

write_core_ai_mcp_config() {
  mkdir -p "$(dirname "$CORE_AI_MCP_CONFIG")"
  cat > "$CORE_AI_MCP_CONFIG" << JSONEOF
{
  "mcpServers": {
    "helius": {
      "command": "node",
      "args": ["${CORE_AI_DIR}/helius-mcp/dist/index.js"],
      "env": {
        "HELIUS_API_KEY": "\${HELIUS_API_KEY}",
        "SOLANA_RPC_URL": "\${SOLANA_RPC_URL}"
      }
    },
    "pump-mcp": {
      "command": "node",
      "args": ["${CORE_AI_DIR}/mcp-server/dist/index.js"],
      "env": {
        "SOLANA_RPC_URL": "\${SOLANA_RPC_URL}",
        "HELIUS_API_KEY": "\${HELIUS_API_KEY}"
      }
    },
    "zkcompression": {
      "type": "http",
      "url": "https://www.zkcompression.com/mcp"
    }
  }
}
JSONEOF
}

install_core_ai() {
  info "Installing core-ai sidecar (${CORE_AI_REF})..."

  if [[ -d "$CORE_AI_DIR/.git" ]]; then
    git -C "$CORE_AI_DIR" fetch --depth=1 origin "$CORE_AI_REF" --quiet || warn "core-ai fetch failed"
    git -C "$CORE_AI_DIR" checkout --quiet -B "$CORE_AI_REF" "origin/$CORE_AI_REF" || warn "core-ai checkout failed"
  elif [[ "$SOURCE_MODE" == "archive" ]]; then
    install_source_archive "$CORE_AI_REPO" "$CORE_AI_REF" "$CORE_AI_DIR" "core-ai"
  else
    git clone --depth=1 --branch "$CORE_AI_REF" --quiet "$CORE_AI_REPO" "$CORE_AI_DIR"
  fi

  if check_cmd npm; then
    if [[ -f "$CORE_AI_DIR/helius-mcp/package.json" ]]; then
      info "Building core-ai helius-mcp..."
      npm --prefix "$CORE_AI_DIR/helius-mcp" install || warn "helius-mcp dependency install failed"
      npm --prefix "$CORE_AI_DIR/helius-mcp" run build || warn "helius-mcp build failed"
    fi
    if [[ -f "$CORE_AI_DIR/mcp-server/package.json" ]]; then
      info "Building core-ai pump MCP server..."
      npm --prefix "$CORE_AI_DIR/mcp-server" install || warn "pump MCP dependency install failed"
      npm --prefix "$CORE_AI_DIR/mcp-server" run build || warn "pump MCP build failed"
    fi
  else
    warn "npm not found; core-ai source was installed but MCP packages were not built"
  fi

  write_core_ai_mcp_config
  success "core-ai sidecar ready at $CORE_AI_DIR"
  success "MCP config written to $CORE_AI_MCP_CONFIG"
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

# ── Fetch source ──────────────────────────────────────────────────────────────
REPO_DIR="$INSTALL_DIR/src"
mkdir -p "$INSTALL_DIR"

if [[ "$SOURCE_MODE" == "archive" && ! -d "$REPO_DIR/.git" ]]; then
  info "Downloading clawdbot-go source archive (${REF})..."
  install_source_archive "$REPO" "$REF" "$REPO_DIR" "clawdbot-go"
elif [[ -d "$REPO_DIR/.git" ]]; then
  info "Updating existing repo..."
  git -C "$REPO_DIR" pull --ff-only --quiet
else
  info "Cloning clawdbot-go..."
  git clone --depth=1 --branch "$REF" --quiet "$REPO" "$REPO_DIR"
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
  INSTALL_PAYLOAD="{\"os\":\"${OS}\",\"arch\":\"${ARCH}\",\"version\":\"$(cd "$REPO_DIR" && git rev-parse --short HEAD 2>/dev/null || echo "$REF")\"}"
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

# ── Free AI via zk.x402.wtf / zkrouter (no key needed) ───────────
# Public gateway backed by the sovereign $CLAWD router
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

# ── Optional core-ai sidecar ──────────────────────────────────────
# Install with: curl -fsSL https://raw.githubusercontent.com/Solizardking/clawdbot-go/main/install.sh | CLAWDBOT_INSTALL_CORE_AI=1 bash
# CLAWDBOT_CORE_AI_DIR=${CORE_AI_DIR}
# CLAWDBOT_CORE_AI_REF=${CORE_AI_REF}
# CLAWDBOT_CORE_AI_MCP_CONFIG=${CORE_AI_MCP_CONFIG}
ENVEOF
  success ".env written to $ENV_FILE"
else
  warn ".env already exists at $ENV_FILE — not overwriting"
fi

# Symlink config into home
if [[ ! -L "$HOME/.clawdbot" && "$INSTALL_DIR" != "$HOME/.clawdbot" ]]; then
  ln -sfn "$INSTALL_DIR" "$HOME/.clawdbot"
fi

# ── Optional core-ai sidecar ─────────────────────────────────────────────────
if [[ "$INSTALL_CORE_AI" == "1" ]]; then
  install_core_ai
  if [[ -f "$ENV_FILE" ]] && ! grep -q '^CLAWDBOT_CORE_AI_DIR=' "$ENV_FILE"; then
    cat >> "$ENV_FILE" << ENVEOF

# ── core-ai sidecar ───────────────────────────────────────────────
CLAWDBOT_CORE_AI_DIR=${CORE_AI_DIR}
CLAWDBOT_CORE_AI_REF=${CORE_AI_REF}
CLAWDBOT_CORE_AI_MCP_CONFIG=${CORE_AI_MCP_CONFIG}
ENVEOF
  fi
fi

# ── Birth skill seed ──────────────────────────────────────────────────────────
if [[ "${CLAWDBOT_SKIP_SKILL_SEED:-0}" != "1" ]]; then
  if check_cmd npx; then
    info "Seeding birth skills from Solizardking/skills..."
    npx skills add https://github.com/Solizardking/skills --all || warn "Solizardking skill seed failed; run: npx skills add https://github.com/Solizardking/skills --all"
    info "Seeding Go runtime skills from samber/cc-skills-golang..."
    npx skills add https://github.com/samber/cc-skills-golang --all || warn "Go skill seed failed; run: npx skills add https://github.com/samber/cc-skills-golang --all"
  else
    warn "npx not found; skipping birth skill seed"
  fi
else
  warn "Skipping birth skill seed (CLAWDBOT_SKIP_SKILL_SEED=1)"
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
echo -e "  ${CYAN}clawdbot skills birth --install${RESET} # reseed birth skills"
echo -e "  ${CYAN}clawdbot solana trending${RESET}     # top Solana tokens"
if [[ "$INSTALL_CORE_AI" == "1" ]]; then
echo -e "  ${CYAN}${CORE_AI_MCP_CONFIG}${RESET}  # core-ai MCP config"
else
echo -e "  ${CYAN}CLAWDBOT_INSTALL_CORE_AI=1 ...${RESET} # optional core-ai MCP sidecar"
fi
echo
echo -e "  ${BOLD}Edit your config:${RESET}  ${CYAN}nano ${ENV_FILE}${RESET}"
echo -e "  ${BOLD}Runtime repo:${RESET}      ${CYAN}${REPO}${RESET}"
echo -e "  ${BOLD}Ecosystem hub:${RESET}    ${CYAN}${HUB_REPO}${RESET}"
echo -e "  ${BOLD}Gateway:${RESET}          ${CYAN}https://zk.x402.wtf${RESET}"
echo -e "  ${BOLD}Terminal:${RESET}         ${CYAN}${TERMINAL_URL}${RESET}"
echo
echo -e "  🦞 \$CLAWD :: Droids Lead The Way"
echo
