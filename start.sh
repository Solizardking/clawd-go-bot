#!/bin/bash
# ─────────────────────────────────────────────────────────────────────
# ClawdBot Go :: One-Shot Start Script
# Installs dependencies, compiles everything, runs animated launcher
# ─────────────────────────────────────────────────────────────────────
set -e

GREEN='\033[1;38;2;20;241;149m'
PURPLE='\033[1;38;2;153;69;255m'
TEAL='\033[1;38;2;0;212;255m'
RED='\033[1;38;2;255;64;96m'
DIM='\033[38;2;85;102;128m'
RESET='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$SCRIPT_DIR"

cd "$ROOT"

echo ""
echo -e "${GREEN}    🦞 ClawdBot Go — One-Shot Start${RESET}"
echo -e "${DIM}    ────────────────────────────────${RESET}"
echo ""

# ── Check Node.js ─────────────────────────────────────────────────
if ! command -v node &>/dev/null; then
  echo -e "${RED}  ✗ Node.js not found. Install it: https://nodejs.org${RESET}"
  exit 1
fi

# ── Check Go ──────────────────────────────────────────────────────
if ! command -v go &>/dev/null; then
  echo -e "${RED}  ✗ Go not found. Install it: https://go.dev/dl/${RESET}"
  exit 1
fi

LAUNCHER_DIR="$SCRIPT_DIR/scripts"
if [ ! -d "$LAUNCHER_DIR/node_modules" ]; then
  echo -e "  ${TEAL}⏳${RESET} Installing launcher dependencies..."
  cd "$LAUNCHER_DIR"
  npm install --no-audit --no-fund --silent 2>/dev/null
  cd "$ROOT"
  echo -e "  ${GREEN}✔${RESET} Launcher ready"
fi

# ── Create build dir ─────────────────────────────────────────────
mkdir -p build

# ── Load .env into environment ───────────────────────────────────
if [ -f ".env" ]; then
  set -a
  source .env 2>/dev/null || true
  set +a
  echo -e "  ${GREEN}✔${RESET} Loaded .env"
fi

echo ""

# ── Run the animated launcher ────────────────────────────────────
exec node "$SCRIPT_DIR/scripts/launch.mjs"
