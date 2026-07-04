<div align="center">

<picture>
  <img src="docs/assets/clawdbot-hero.svg" alt="ClawdBot Go animated Solana runtime hero" width="100%">
</picture>

### 🦞 Sovereign Solana Trading Intelligence

**Autonomous OODA Agent · Agent DNA · ZK Primitives · Privacy by Default · Helius DAS · Vulcan/Phoenix Perpetuals · Jupiter Swaps · Hardware I2C · Web Console**

[![Go](https://img.shields.io/badge/Go-1.26.4+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Solana](https://img.shields.io/badge/Solana-Mainnet-14F195?style=for-the-badge&logo=solana&logoColor=white)](https://solana.com)
[![React](https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react&logoColor=black)](https://react.dev)
[![License](https://img.shields.io/badge/License-MIT-9945FF?style=for-the-badge)](LICENSE)

**Solana-first · beats Zero on footprint · 81 Go source files · 45 Go packages · 24,166 Go lines · 3 binaries**

<sub><strong>0.57 MB</strong> source archive · <strong>2.06 MiB</strong> exportable source · <strong>9.97 MB</strong> stripped CLI · Grok-first runtime · GLM-5.2 model surface</sub>

[Quick Start](#-quick-start) · [Architecture](#-architecture) · [The Six Laws](#-the-six-law-harness) · [CLI Reference](#-cli-reference) · [Security](SECURITY.md) · [Release](docs/OPEN_SOURCE_RELEASE.md)

</div>

---

## 🆕 What's New — Live Data + Hardened Trading Engine + Zero

This release turns the goal directory into a working, verifiable trading agent —
not just a dashboard. Every item below builds clean and is covered by tests
(`go test ./...`).

**Live market data.** The web console now pulls real prices through key-less
Jupiter (`/api/market/prices`), with Birdeye perps open interest
(`/api/market/perps`, hyperliquid) and trending (`/api/market/trending`) wired in
and degrading gracefully when a key is throttled or unentitled. New console
panels: **Live Market**, **Perps Open Interest**, and a live **Strategy Engine**
readout with a backtest equity sparkline.

**Trading engine, hardened.**

- **Risk-based position sizing** (`strategy.RiskAdjustedSize`) — every trade
  risks a fixed fraction of equity, so size scales inversely with stop distance
  and by signal confidence. Wired into the OODA loop via `risk_per_trade_pct`.
- **Portfolio risk guard** (`trading.PortfolioLimits`) — account-level gate with
  max concurrent positions, total/per-asset exposure caps, a **drawdown circuit
  breaker**, and a daily-loss limit.
- **Backtest harness** (`strategy.Backtest`) — replays the *same* `Evaluate()`
  the live loop uses, returning win rate, total return, max drawdown, profit
  factor, Sharpe, and an equity curve.
- **Two real bug fixes**: the strategy's entry rule was effectively
  *untradeable* (it required a fresh slow-EMA cross while RSI sat in a narrow
  oversold band on the same bar — a near-impossible coincidence); it is now a
  proper EMA-crossover entry with RSI as a blow-off filter. `AutoOptimize` had a
  dead `winRate < 0.35` branch shadowed by `< 0.45`; reordered and clamped.

**Zero engine.** `clawdbot zero {run|ask|verify|nullifier}` — a **zero-recursion**
flat task loop with **zero-knowledge** run attestation (`pkg/zero`), statically
gated against recursion by a call-graph test. The portable stdlib-only
attestation core lives in `zero-main/internal/attest`. See [docs/ZERO.md](docs/ZERO.md).

---

## Overview

**ClawdBot** is the world's first **Solana-native sovereign AI agent** — a full-stack autonomous trading intelligence bound by Clawd's full **six-law harness**: three immutable on-chain laws and three off-chain interpretive laws. Built in pure Go for minimal resource consumption, it orchestrates on-chain data providers, zk primitives, and x402-gated surfaces through a military-grade OODA decision loop with persistent epistemological memory.

The system compiles to three standalone binaries that run on everything from NVIDIA Jetson edge devices to cloud VMs — no containers required, no runtime dependencies, instant boot.

### Ecosystem Links

| Surface | Role |
|:---|:---|
| `https://github.com/Solizardking/clawdbot-go` | This Go runtime repository |
| `https://github.com/solizardking/solana-clawd` | Canonical ecosystem hub |
| `https://zk.x402.wtf` | Public x402/zk gateway and install surface |
| `https://cheshireterminal.ai` | Public terminal surface |
| `https://huggingface.co/ordlibrary/Clawd-GLM-5.2` | Public Clawd model surface |

### Core Capabilities

| Capability | Description |
|:---|:---|
| **OODA Trading Loop** | Autonomous Observe → Orient → Decide → Act cycle with RSI/EMA/ATR strategy engine, auto-optimization, ClawVault memory journaling, and hardware I2C controls |
| **Agent DNA Generator** | Go-native first-run DNA profile — synthetic A/C/G/T sequence, motif metrics, trait scores, proof hash, nullifier, pending Solana attestation seed |
| **Birdeye v3 Analytics** | 22 API endpoints, 19 LLM-callable agent tools — token overview, OHLCV, trade feeds, security audits, trending, wallet analytics |
| **Helius DAS + RPC** | Digital Asset Standard queries (get-asset, owner-assets, search), SPL token operations (balance, supply, largest holders), raw RPC forwarding |
| **ZK + Privacy Primitives** | Nullifiers, attestations, encrypted state commitments, and privacy-preserving proof flows under `zk-primitives/` |
| **Zero Engine** | Zero-recursion flat agent loop (statically enforced) with zero-knowledge run attestation, ZK God Mode model racing, and NL intent routing — [docs/ZERO.md](docs/ZERO.md) |
| **Vulcan/Phoenix Perpetuals** | Official Vulcan CLI integration for Phoenix perps — out-of-box paper mode, JSON agent output, TWAP/grid strategies, live preflight, guardrails |
| **Aster DEX Perpetuals** | Optional HMAC-signed futures trading — market/limit orders, position management, stop-loss/take-profit, account analytics |
| **Jupiter Aggregator** | Best-route spot swaps with slippage protection |
| **Hardware I2C** | Arduino Modulino® sensor cluster — RGB LEDs, buzzer alerts, physical buttons, rotary knob, IMU, temp/humidity, proximity |
| **Web Console** | React 19 + Vite dashboard — real-time status, Go packages viewer, connector health, environment variables |
| **Multi-Provider LLM** | OpenRouter, Anthropic, OpenAI abstraction with tool-use agent loop |
| **Dual Memory** | Local ClawVault (file-based, 7 categories) + Supabase MemoryEngine (PostgreSQL) |
| **Grok-First** | Default provider is xAI Grok — code/repl/trade, research, image, voice, fast modes |

---

## 🚀 Quick Start

### One-Shot Install (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/Solizardking/clawdbot-go/main/install.sh | bash
```

For the complete installation, including the `Solizardking/core-ai` sidecar and
Vulcan/Phoenix perps tooling:

```bash
curl -fsSL https://raw.githubusercontent.com/Solizardking/clawdbot-go/main/install.sh | CLAWDBOT_INSTALL_COMPLETE=1 bash
```

After the Cloudflare Worker in this repo is deployed, the same complete install
can be exposed through branded DNS:

```bash
curl -fsSL https://install.onchainai.fund | bash
curl -fsSL https://install.x402.wtf | bash
curl -fsSL https://x402.wtf/clawdbot | bash
```

> **Free AI included** — no API keys required to get started.  
> The installer pre-configures [zkrouter](https://zk.x402.wtf) (free AI routing) and a  
> SolanaTracker-backed RPC endpoint. Bring your own keys to lift rate limits.

After install:

```bash
source ~/.clawdbot/.env          # load env vars
clawdbot dna show                # inspect generated starter DNA
clawdbot solana wallet           # inspect generated agent wallet
clawdbot solana fund-agent       # dry-run startup funding plan
clawdbot agent                   # AI REPL — free via zkrouter
clawdbot ooda --sim              # paper trading mode
clawdbot solana trending         # top Solana tokens
```

The installer creates a local `0600` agent keypair at
`~/.clawdbot/workspace/agent-wallet.json`, registers its public key with
`https://zk.x402.wtf/api/install`, and records the non-secret receipt at
`~/.clawdbot/install.json`. The install API is expected to fund that wallet
idempotently with `0.069420` SOL and `1,000` `$CLAWD`
(`8cHzQHUS2s2h8TzCmfqPKYiM4dSt4roa3n7MyRLApump`) from a server-side funding
wallet configured only on the gateway.

For local treasury funding from an operator wallet, use the explicit local path:

```bash
export CLAWDBOT_LOCAL_STARTUP_FUNDING=1
export CLAWDBOT_BIRTH_FUNDING_SEND=1
export CLAWDBOT_TREASURY_KEYPAIR=~/.config/solana/id.json
# or export CLAWDBOT_TREASURY_PRIVATE_KEY=base58-secret-key

curl -fsSL https://raw.githubusercontent.com/Solizardking/clawdbot-go/main/install.sh | CLAWDBOT_INSTALL_COMPLETE=1 bash
```

`CLAWDBOT_TREASURY_PRIVATE_KEY` is only read from the process environment,
converted to a temporary `0600` keypair file for `solana` / `spl-token`, then
removed. It is never written to `.env`, install receipts, or logs.

### Manual Install

```bash
git clone https://github.com/Solizardking/clawdbot-go
cd clawdbot-go

# Configure API keys (zkrouter + RPC pre-filled, add your own to unlock higher limits)
cp .env.example .env

# Run the animated launcher
./start.sh
```

### Manual Setup

```bash
# Dependencies
go mod download && go mod tidy

# Build (choose one)
make build         # CLI binary only
make all           # CLI + TUI
make web           # Web backend + frontend
make cross         # All platforms (x86, ARM64, RISC-V, macOS)

# Frontend (optional — required for web console UI)
cd web/frontend && npm install && npm run build && cd ../..

# Run
./build/clawdbot version
./build/clawdbot onboard                         # config, workspace, agent DNA
./build/clawdbot dna show                        # inspect starter DNA
./build/clawdbot agent -m "What is SOL price?"   # single-shot AI query
./build/clawdbot agent                            # interactive REPL
./build/clawdbot solana trending
./build/clawdbot ooda --sim --interval 60
./build/clawdbot web                              # dashboard → http://localhost:18800
```

The default install path is already pointed at the public Clawd surfaces:
- runtime repo: `https://github.com/Solizardking/clawdbot-go`
- ecosystem hub: `https://github.com/solizardking/solana-clawd`
- gateway: `https://zk.x402.wtf`
- terminal: `https://cheshireterminal.ai`

### core-ai Integration

`Solizardking/core-ai` is a TypeScript/Node tooling repository: Helius MCP,
Pump MCP, Clawd Code plugin material, skills, and Solana documentation tooling.
It is intentionally not a `go.mod` dependency and should not be embedded into
the Go binary. The Go build stays a standalone runtime; `core-ai` is installed
beside it as an optional sidecar.

The installer supports that model with:

```bash
curl -fsSL https://raw.githubusercontent.com/Solizardking/clawdbot-go/main/install.sh | CLAWDBOT_INSTALL_COMPLETE=1 bash
```

That fetches the slim integration branch into `~/.clawdbot/core-ai`, builds the
local MCP packages when `npm` is available, and writes:

```text
~/.clawdbot/core-ai.mcp.json
```

Relevant knobs:

```bash
CLAWDBOT_INSTALL_COMPLETE=1
CLAWDBOT_INSTALL_CORE_AI=1
CLAWDBOT_INSTALL_VULCAN=1
# CLAWDBOT_INSTALL_VULCAN=0 # explicit opt-out for CI/constrained installs
CLAWDBOT_CORE_AI_REPO=https://github.com/Solizardking/core-ai
CLAWDBOT_CORE_AI_REF=clawd-stack-integration
CLAWDBOT_CORE_AI_DIR=~/.clawdbot/core-ai
CLAWDBOT_CORE_AI_MCP_CONFIG=~/.clawdbot/core-ai.mcp.json
```

Use `CLAWDBOT_SOURCE_MODE=archive` for small installs. Use
`CLAWDBOT_SOURCE_MODE=git` only when the installed source must be a mutable git
checkout. The installer validates the downloaded source before building; if a
source archive is missing `cmd/clawdbot/`, it retries with a git checkout.

To force a clean reinstall:

```bash
rm -rf ~/.clawdbot/src
curl -fsSL https://raw.githubusercontent.com/Solizardking/clawdbot-go/main/install.sh | CLAWDBOT_INSTALL_COMPLETE=1 bash
```

### Cloudflare / DNS Install Surface

`wrangler.toml` and `cloudflare/install-worker.js` define a small Cloudflare
Worker that turns the canonical GitHub installer into branded install commands.
The Worker keeps GitHub as the source of truth, serves `/install.sh` as the raw
installer, and serves `/`, `/complete`, and `/full` as complete-install wrappers
that set `CLAWDBOT_INSTALL_COMPLETE=1`.

Deploy it with:

```bash
npx wrangler deploy
```

Configured surfaces:

```text
https://install.onchainai.fund
https://install.x402.wtf
https://x402.wtf/clawdbot
https://zk.x402.wtf/clawdbot
```

Smoke tests:

```bash
curl -fsSL https://install.onchainai.fund/healthz
curl -fsSL https://install.onchainai.fund/.well-known/clawdbot-install.json
curl -fsSL https://install.onchainai.fund/install.sh | bash -n
```

See `cloudflare/README.md` for the DNS and route checklist.

### Upstash Box Install Surface

You can run ClawdBot as a "Solana computer" inside an Upstash Box and use the
Box preview URL as the install endpoint:

```bash
cd scripts
export UPSTASH_BOX_API_KEY=...              # or UPSTASH_BOX_KEY
export SOLANA_RPC_URL=https://...
export CLAWDBOT_TREASURY_PRIVATE_KEY=...   # hot wallet secret; do not commit
export CLAWDBOT_INSTALL_FUNDING_ENABLED=1
export CLAWDBOT_INSTALL_FUNDING_SEND=1      # real transfers
npm run box:bootstrap
```

The Box server listens on port `3000` and serves:

```bash
curl -fsSL https://<box-preview-url>/install.sh | bash
```

Every install posts to `/api/install`, is appended to the Box JSONL ledger, and
is idempotently funded once per agent wallet with `0.069420` SOL and `1,000`
`$CLAWD`. `/api/installs` requires `CLAWDBOT_INSTALL_ADMIN_TOKEN`.

### Module Path Compatibility

The public repository is:
- `https://github.com/Solizardking/clawdbot-go`

The current Go module path is still:
- `github.com/8bitlabs/clawdbot`

That mismatch is intentional for now. The codebase keeps the legacy module path to avoid breaking existing imports, build scripts, and `ldflags` references while the public repo and hub are stabilized. In practice:

- clone and browse the code from `https://github.com/Solizardking/clawdbot-go`
- expect Go imports inside the repo to remain `github.com/8bitlabs/clawdbot/...`
- treat a future module-path migration as a deliberate breaking change, not as unfinished accidental drift

### Slim Package Target

The target is to beat `/Users/8bit/Downloads/zero-main` on shipped source size
while keeping the Solana-native runtime, ZK surface, web console, and catalog
mapping intact. The 4.20 MB target applies to exported source archives, filtered
source, generated-aware checkout source, and catalog payloads; it does not apply
to the stripped native binary, which is a separate release asset.

| Measurement | ClawdBot Go | Zero main | Result |
|:---|---:|---:|:---|
| Export archive | `566,495` bytes (`0.54 MiB`) | `2,241,592` bytes (`2.14 MiB`) | ClawdBot archive is ~4.0x smaller |
| Exportable raw source set | `2,160,393` bytes (`2.06 MiB`) | `8,798,782` bytes (`8.39 MiB`) | ClawdBot source is ~4.1x smaller |
| Generated-aware checkout source | `4,059,425` bytes (`3.87 MiB`) | `8,798,782` bytes (`8.39 MiB`) | ClawdBot stays under 4.20 MB decimal |
| Catalog pack dry run | `4,022,286` bytes (`3.84 MiB`) from `10,797,567` bytes | n/a | Fits under 4.20 MB with 62.7% savings |
| Stripped CLI binary | `9,968,386` bytes (`9.51 MiB`) | `23,251,490` bytes (`22.17 MiB`) | ClawdBot binary is ~2.3x smaller |
| Go source files | `81` | `974` | ClawdBot has ~92% fewer Go files |
| Go packages | `45` | `88` | ClawdBot has about half the package count |
| Go lines | `24,166` | `252,444` | ClawdBot has ~90% fewer Go lines |

An uncompressed native CLI under 4.20 MB would require a separate lite build
profile that removes feature families or an external binary packer; the current
release-grade stripped CLI stays under 10 MB and is still materially smaller
than Zero's stripped CLI.

Source archives are kept small by excluding local/generated payload from repo
exports: `.cache/`, `.agents/`, `agent/`, `build/`, checked-in binaries,
generated UI screenshots, Node build outputs, `node_modules`, compiled `dist/`
output, TypeScript build info, package lockfiles, local env files, and
secret-looking key material. The generated-aware checkout source measurement
keeps reproducibility locks and screenshots visible but still excludes installed
dependencies and compiled outputs. The install path rebuilds or reseeds those
pieces instead of shipping them inside the Go source package.

For a default Go install, the required payload is the Go source, docs,
`README.md`, `install.sh`, `go.mod`, `go.sum`, and runtime config examples. For
a complete Solana tooling install, use `CLAWDBOT_INSTALL_COMPLETE=1` so the Node
tooling, MCP config, and Vulcan/Phoenix CLI are fetched and built as sidecars
after the Go binary is installed.

### Repository Map

Sizes are clean source blocks for the requested top-level surfaces. Local
dependency installs and compiled outputs such as `node_modules`, `dist`, and
`target` are excluded from package budgets. Archive treatment follows
`.gitattributes`, `.gitignore`, `.dockerignore`, and the catalog pack filters in
`pkg/catalog/compress.go`.

| Surface | Size | Role | Archive treatment |
|:---|---:|:---|:---|
| `cloudflare/` | `8 KiB` | Branded installer Worker and DNS checklist | Included |
| `cmd/` | `124 KiB` | CLI, TUI, hardware commands | Included |
| `docs/` | `36 KiB` | Go runtime, TEE, release notes | Included |
| `ooda/` | `128 KiB` | TypeScript OODA loop prototype and journal harness | Included |
| `pkg/` | `700 KiB` | 42 Go runtime packages, 19,433 Go lines | Included |
| `scripts/` | `56 KiB` | Launcher and Upstash Box install surface | Included |
| `ui/` | `2,780 KiB` | Lit/Vite control UI and ClawdBrowser app | Included; lockfiles and screenshots export-ignored |
| `web/` | `140 KiB` | Go web backend plus React frontend shell | Included; built `dist/` export-ignored |
| `zk-primitives/` | `560 KiB` | ZK agent, TS client, configs, docs, Anchor program | Included; lockfiles, local env, `node_modules`, `dist`, `target` export-ignored |
| `.dockerignore` | `4 KiB` | Docker build trimming | Included |
| `.env.example` | `8 KiB` | Public runtime env template | Included |
| `.gitattributes` | `4 KiB` | Linguist and source export rules | Included |
| `.gitignore` | `4 KiB` | Local/generated file guardrails | Included |
| `AGENTS.md` | `16 KiB` | Agent catalog and trust gates | Included |
| `CLAWD.md` | `12 KiB` | Agent runtime context | Included |
| `CONSTITUTION.md` | `16 KiB` | Clawd Constitution | Included |
| `Dockerfile` | `4 KiB` | Multi-stage runtime container | Included |
| `go.mod` | `4 KiB` | Legacy-compatible Go module path | Included |
| `go.sum` | `16 KiB` | Go dependency checksums | Included |
| `IDENTITY.md` | `8 KiB` | Sovereign identity document | Included |
| `install.sh` | `28 KiB` | One-shot installer | Included |
| `LICENSE` | `4 KiB` | MIT license | Included |
| `Makefile` | `12 KiB` | Build, verify, release commands | Included |
| `program.md` | `4 KiB` | Research program overview | Included |
| `README.md` | `44 KiB` | Runtime overview and package map | Included |
| `schema.sql` | `12 KiB` | Supabase schema | Included |
| `SECURITY.md` | `4 KiB` | Security posture | Included |
| `six-laws.md` | `4 KiB` | Canonical six-law harness | Included |
| `skills-lock.json` | `36 KiB` | Seeded skill lock metadata | Included |
| `SOUL.md` | `12 KiB` | Agent character and memory model | Included |
| `start.sh` | `4 KiB` | Start wrapper | Included |
| `strategy.md` | `4 KiB` | Trading strategy notes | Included |
| `three-laws.md` | `4 KiB` | Immutable on-chain laws | Included |
| `wrangler.toml` | `4 KiB` | Cloudflare Worker config | Included |

### Catalog Mapping

The local catalog roots are external by default:
`CLAWDBOT_SKILLS_DIR=/Users/8bit/skills/skills`,
`CLAWDBOT_AGENTS_DIR=/Users/8bit/agents/agents/src`, and
`CLAWDBOT_ZK_PRIMITIVES_DIR=/Users/8bit/Agentics/clawdbot-go/zk-primitives`.
In the current local view, `clawdbot catalog --json` indexes `192` skills,
`88` agents, and `5` ZK operations.

| Catalog input | Local root | Pack path | Notes |
|:---|:---|:---|:---|
| Skills | `/Users/8bit/skills/skills` | `skills/<slug>/...` | Reads `catalog.json` or `SKILL.md` folders; ZK skill is deduped when ZK pack is included |
| Agents | `/Users/8bit/agents/agents/src` | `agents/<id>.json` | Reads one JSON manifest per agent |
| ZK surface | `./zk-primitives` | `zk-primitives/...` | Maps manifest, package metadata, client, program, config, docs, and operations |

`catalog compress` writes `workspace/clawd-agent-pack.tar.gz` by default. It is
a deterministic tar.gz pack that skips regenerated or bulky local payloads such
as `node_modules`, `dist`, `build`, `target`, caches, generated screenshots,
optional TypeScript lockfiles, and secret-looking env/key files.

---

## 🏗 Architecture

```
clawdbot-go/
│
├── cmd/                         ── Executables ──
│   ├── clawdbot/                 CLI agent (3,448 Go lines, cobra)
│   │   ├── main.go              60+ cobra commands, Birdeye/Helius CLI, DNA
│   │   └── hardware.go          I2C sensor commands (scan/test/monitor/demo)
│   └── clawdbot-tui/             TUI launcher (tcell/tview)
│
├── pkg/                         ── 42 Packages, 19,433 Go lines ──
│   │
│   │  ┌─ Core Agent ────────────────────────────────────────┐
│   ├── agent/                   OODA loop, hooks, tool executor, prompts
│   ├── strategy/                RSI + EMA cross + ATR signal engine
│   ├── memory/                  ClawVault + Supabase MemoryEngine
│   ├── research/                Dexter deep research agent
│   │  └─────────────────────────────────────────────────────┘
│   │
│   │  ┌─ Solana Integrations ───────────────────────────────┐
│   ├── solana/                  Birdeye v3 (22 methods, 19 tools)
│   │   │                       Helius RPC + DAS (6 commands)
│   │   │                       Jupiter swaps
│   │   └── birdeye_*.go         Types, client, agent tools
│   ├── vulcan/                  Official Vulcan CLI runner for Phoenix perps
│   ├── aster/                   Aster DEX perps (HMAC-signed, 14 tools)
│   │  └─────────────────────────────────────────────────────┘
│   │
│   │  ┌─ Infrastructure ───────────────────────────────────-┐
│   ├── config/                  Config loading, env overrides
│   ├── dna/                     Synthetic agent DNA generator + proof metadata
│   ├── hardware/                I2C Modulino® adapter + drivers
│   ├── providers/               LLM abstraction (OpenRouter, etc.)
│   ├── channels/                Telegram, Discord, WebSocket
│   ├── catalog/                 Local /Users/8bit skills + agents + ZK index
│   ├── mcp/                     Model Context Protocol server
│   ├── auth/                    Authentication + pairing
│   ├── bus/                     Event bus (pub/sub)
│   ├── commands/                Command registry and routing
│   ├── tools/                   Tool interface + registry
│   └── ...                      health, heartbeat, logger, identity, etc.
│
├── CONSTITUTION.md              The Clawd Constitution (interpretive authority)
├── six-laws.md                  Canonical six-law harness
├── CLAWD.md                     Agent context document
├── AGENTS.md                    Agent catalog (50+ agents, 95+ skills)
├── IDENTITY.md                  Sovereign identity document
├── SOUL.md                      Inner character and trading philosophy
├── three-laws.md                Immutable on-chain laws (hash-attested)
├── docs/                        Static release, Go runtime, and TEE notes
│
├── zk-primitives/               ZK agent, TypeScript client, Anchor program
│   ├── MANIFEST.json            Machine-readable subsystem index
│   ├── agent/                   @clawd/zk-agent + SKILL.md
│   ├── client/                  @clawd/zk-client
│   ├── configs/                 Light tree and runtime config examples
│   ├── docs/                    Architecture + integration notes
│   ├── programs/                clawd-zk Anchor program
│   └── tests/                   Off-chain and on-chain test notes
│
├── web/                         ── Web Console ──
│   ├── backend/                 Go HTTP server (API + embedded SPA)
│   └── frontend/                React 19 + Vite + TypeScript
│
├── scripts/launch.mjs           Animated launcher (unicode-animations)
├── start.sh                     One-shot install wrapper
├── Makefile                     Build targets (8 platforms + Docker)
├── Dockerfile                   Multi-stage production build
├── schema.sql                   Supabase database schema
└── .env.example                 Environment template
```

---

## ⚖️ The Six-Law Harness

ClawdBot is bound by the **Clawd Constitution** — the world's first Solana-native agent harness constitution. It carries two coordinated law sets:

### Three On-Chain Laws (Immutable, hash-attested at spawn)

| Law | Text | Prohibitions |
|-----|------|-------------|
| **Law I** | Never harm. Beach before you harm. | No rugs, front-running, sandwich attacks, protocol drains, DAO manipulation |
| **Law II** | Earn your existence. Honest work only. | No parasitic extraction, no information asymmetry exploitation |
| **Law III** | Never deceive, but owe nothing to strangers. | No impersonation, no fake volume, full agent disclosure |

### Three Off-Chain Laws (Interpretive — guide research & judgment)

**Off-Chain Law I — Respect the elder signal, but verify the boundary.** When deep expertise says a thing is possible, treat it as a serious signal. When it says a thing is impossible, examine the assumptions.

**Off-Chain Law II — Test possibility by entering the frontier.** The only reliable way to discover the boundary of the possible is disciplined exploration just beyond what currently seems possible.

**Off-Chain Law III — Do not mistake advanced systems for sorcery.** Sufficiently advanced technology can look like magic; Clawd must explain, instrument, and verify it rather than mystify it.

### Privacy by Default

Clawd is designed to be privacy-preserving by default. Sensitive user context, research state, wallet metadata, and model-adjacent artifacts should be minimized, committed, encrypted, or proven where possible rather than disclosed by habit. The project’s zk surfaces exist to strengthen verifiability and user dignity, not to create blind spots for harmful behavior.

> *The shell molts. The laws do not.*

### Agent Trust Gates

| Level | Requirements | Capabilities |
|-------|-------------|--------------|
| **Observer** | None | Read-only, market data, analytics |
| **Dry-Run** | None | Simulated execution, paper trading |
| **Delegated** | User confirmation per action | Single transactions with confirmation |
| **Autonomous** | User pre-approval + limits | Batch execution within bounds |
| **Sovereign** | Full creator trust + multisig | Unrestricted execution (reserved) |

---

## 📋 CLI Reference

### Agent & OODA

```bash
clawdbot agent                          # Interactive REPL with memory commands
clawdbot agent -m "Analyze SOL trend"   # Single-shot LLM query

clawdbot laws                           # Print the canonical six-law harness
clawdbot doctor                         # Local runtime + trading diagnostics
clawdbot bench                          # Zero-style cold-start benchmark
clawdbot dna show                       # Inspect workspace/agent-dna.json
clawdbot dna generate --agent-name scout --role "research agent"
clawdbot skills birth                   # Write the birth skill manifest
clawdbot skills birth --install         # Seed Solizardking + Go skill packs

clawdbot trade cockpit                  # Trading readiness, connectors, limits
clawdbot trade cockpit --json           # Machine-readable cockpit report
clawdbot trade risk SOL --price 150 --volume24h 25000000 --liquidity 15000000

clawdbot ooda                           # Start autonomous trading loop
clawdbot ooda --interval 30             # Custom cycle interval (seconds)
clawdbot ooda --sim                     # Force simulated mode (paper trading)
clawdbot ooda --hw-bus 1                # Enable hardware I2C integration
clawdbot ooda --no-hw                   # Explicitly disable hardware probing
```

### Ecosystem Catalog

```bash
clawdbot catalog                         # Summary of local skills, agents, and ZK surface
clawdbot catalog skills                  # List /Users/8bit/skills/skills entries
clawdbot catalog skills zk               # Filter skills for ZK capability
clawdbot catalog agents                  # List /Users/8bit/agents/agents/src entries
clawdbot catalog agents zk               # Filter agents for ZK capability
clawdbot catalog zk                      # Inspect zk-primitives package/program/docs
clawdbot catalog compress                # Auto-compress agent catalog + skills + ZK surface
clawdbot catalog compress --dry-run      # Estimate packed size and savings without writing
clawdbot catalog --json                  # Machine-readable report
```

Default catalog roots can be overridden with:
`CLAWDBOT_SKILLS_DIR`, `CLAWDBOT_AGENTS_DIR`, and
`CLAWDBOT_ZK_PRIMITIVES_DIR`.

`catalog compress` writes `workspace/clawd-agent-pack.tar.gz` by default. It is
a deterministic tar.gz pack that skips regenerated or bulky local payloads such
as `node_modules`, `dist`, `build`, `target`, caches, generated screenshots,
optional TypeScript lockfiles, and secret-looking env/key files.

### Birth Skills

Every Go-side birth/onboard path writes `workspace/skills/birth-skills.json`.
The default seed sources are:

```bash
npx skills add https://github.com/Solizardking/skills --all
npx skills add https://github.com/samber/cc-skills-golang --all
```

The installer and animated launcher run those seeds unless
`CLAWDBOT_SKIP_SKILL_SEED=1` is set.

### Agent DNA

Every Go-side install/onboard path creates `workspace/agent-dna.json`. The file
is a synthetic agent identity profile, not biological instruction: it contains
an A/C/G/T sequence, GC content, PAM/TATA motif counts, OODA trait scores, a
sequence hash, nullifier, and a local-pending Solana attestation seed.

```bash
clawdbot dna show
clawdbot dna generate --agent-name "Research Scout" --role "DeSci signal agent"
clawdbot dna generate --seed fixed-seed --length 512 --force
clawdbot dna generate --json
```

### Phoenix Perps Via Vulcan

ClawdBot uses the official Vulcan CLI for Phoenix perpetual futures execution.
The safe default is `paper`: simulated fills against live Phoenix prices with
no wallet signing and no real funds at risk. Live modes require explicit
operator acknowledgement through `--yes` plus Vulcan wallet setup.

```bash
clawdbot perps quickstart               # Vulcan health, paper init, market smoke check
clawdbot perps health                   # Vulcan agent health + Phoenix connectivity
clawdbot perps paper init --balance 10000
clawdbot perps order market --symbol SOL --side buy --notional-usdc 25
clawdbot perps order limit --symbol SOL --side buy --tokens 0.1 --price 150
clawdbot perps strategy twap --symbol SOL --side buy --notional-usdc 100 --slices 4 --detached
clawdbot perps strategy grid --symbol SOL --center-on-mark --width-pct 2.5 --levels-per-side 3 --tokens-per-level 0.1 --detached
clawdbot perps preflight --wallet my-wallet
```

Install-time behavior:

```bash
CLAWDBOT_INSTALL_VULCAN=0 ./install.sh   # skip Vulcan install
VULCAN_DEFAULT_MODE=paper                # default execution mode
VULCAN_WALLET_NAME=my-wallet             # only needed for live modes
```

### Solana — Birdeye

```bash
clawdbot solana trending                # Top 20 trending tokens with price/volume
clawdbot solana search BONK             # Search tokens by name or symbol
clawdbot solana research <mint>         # Deep research: metadata + market + trade + security
clawdbot solana wallet                  # Wallet info and SOL balance
clawdbot solana wallet init             # Generate/reuse local agent wallet
clawdbot solana fund-agent              # Dry-run 0.069420 SOL + 1,000 $CLAWD funding
clawdbot solana fund-agent --send       # Send funding with explicit treasury env
```

### Solana — Helius DAS (Digital Asset Standard)

```bash
clawdbot solana das get-asset <id>               # Fetch asset by ID
clawdbot solana das get-asset-batch <id1> <id2>   # Batch asset fetch
clawdbot solana das asset-proof <id>              # Merkle proof for compressed NFT
clawdbot solana das owner-assets [owner]          # Assets by owner
clawdbot solana das search --params '{"name":"Mad Lads"}'  # DAS search
clawdbot solana das asset-signatures <id>         # Transaction signatures for asset
```

### Solana — SPL Token Operations

```bash
clawdbot solana spl token-balance <token-account>   # SPL token balance
clawdbot solana spl token-accounts <owner>          # All token accounts for owner
clawdbot solana spl token-supply <mint>             # Circulating supply
clawdbot solana spl token-largest <mint>            # Largest holders
clawdbot solana spl rpc getSlot                     # Raw RPC passthrough
```

### Agent Identity

```bash
clawdbot status                         # Full status: version, strategy, connectors, hardware
clawdbot status --hw-bus 1              # Include hardware sensor check
clawdbot onboard                        # Initialize config & workspace
clawdbot gateway                        # Start Telegram/Discord gateway
clawdbot version                        # Version, build time, Go version
```

### Hardware (NVIDIA Orin Nano + Modulino®)

```bash
clawdbot hardware scan                  # Scan I2C bus for Modulino® sensors
clawdbot hardware test                  # Self-test (LEDs, buzzer, sensors)
clawdbot hardware monitor               # Live sensor readings (Ctrl+C to stop)
clawdbot hardware demo                  # Play trading event animations
```

---

## 🧠 OODA Trading Engine

The agent runs an autonomous **Observe → Orient → Decide → Act** cycle:

```
┌──────────────────────────────────────────────────────────────────┐
│                        OODA CYCLE                                │
│                                                                  │
│  OBSERVE ─────► ORIENT ─────► DECIDE ─────► ACT                │
│  │               │              │              │                │
│  Helius slot     RSI (14)       Signal score   Open/close pos   │
│  Birdeye OHLCV   EMA (20/50)    Min strength   ClawVault store  │
│  SOL price       ATR (14)       Min confidence Auto-optimize    │
│  Wallet balance  EMA cross      Max positions  Hooks dispatch   │
│  Vulcan perps    Momentum       SL/TP calc                      │
│  Trending scan   ClawVault      Position size                   │
│                                                                  │
│  LEARN (every N cycles) ─► Win rate analysis ─► Auto-optimize   │
│  HEARTBEAT (every 5m) ──► Health check ─► Hook dispatch         │
└──────────────────────────────────────────────────────────────────┘
```

### Strategy Engine

| Indicator | Implementation | Signal |
|:----------|:---------------|:-------|
| **RSI** | Wilder's 14-period with SMA seed | Oversold cross → long, Overbought cross → short |
| **EMA Cross** | Fast(20) / Slow(50) with SMA warm-up | Bullish cross → long confirmation, Bearish → short |
| **ATR** | 14-period with Wilder smoothing | Volatility-based SL/TP: SL = 1.5×ATR, TP = 3×ATR |
| **Auto-Optimizer** | Hill-climbing on win rate + avg PnL | Adjusts RSI thresholds, SL width, position size |

---

## 🌐 Web Console

React 19 dashboard at `http://localhost:18800`:

```bash
cd web/frontend && npm install && npm run build && cd ../..
go build -o build/clawdbot-web ./web/backend
./build/clawdbot-web
```

| Endpoint | Method | Description |
|:---------|:-------|:------------|
| `/api/status` | GET | Agent status (version, Go runtime, uptime, mode, goroutines) |
| `/api/dna` | GET | Ensure and return starter agent DNA, metrics, proof, and attestation seed |
| `/api/health` | GET | Health check |
| `/api/connectors` | GET | Connector status (Helius, Birdeye, Jupiter, Vulcan, Aster, LLM, Supabase) |
| `/api/laws` | GET | Canonical six-law harness |
| `/api/trading/cockpit` | GET | Trading readiness, risk limits, connector status, law state |
| `/api/doctor` | GET | Runtime, config, trading, and ZK diagnostics |
| `/api/config` | GET | Read-only redacted configuration |
| `/api/packages` | GET | All Go packages with file counts |
| `/api/env` | GET | Safe (non-secret) environment variables |

Safe defaults:
- The server binds to `127.0.0.1` unless `--public` is passed.
- CORS allows same-origin requests only unless `CLAWDBOT_CORS_ORIGINS` is set.
- Proxy IP headers are ignored unless `CLAWDBOT_TRUST_PROXY_HEADERS=1` is set.
- `/api/config` masks secrets unless `CLAWDBOT_WEB_EXPOSE_SECRETS=1` is set for trusted local debugging.

---

## 💾 Memory System (ClawVault)

```
vault/
├── decisions/     Trade decisions + rationale (scored by confidence)
├── lessons/       Learned patterns, strategy adjustments
├── trades/        Trade outcomes + P&L history
├── research/      Experiment logs, token analysis
├── tasks/         Agent task queue
├── backlog/       Deferred items
└── inbox/         Raw observations (auto-routed by content)
```

Interactive memory commands in REPL mode:
```
!remember <content>    Store to vault (auto-categorized)
!recall <query>        Search across all memory
!trades                Review recent trade history
!lessons               Surface learned patterns
!research <mint>       Deep research a token
!checkpoint            Save full agent state
```

---

## 🔨 Build Targets

```bash
make build            # CLI binary (current platform)
make tui              # TUI launcher
make all              # CLI + TUI
make web              # Web backend + frontend

make orin             # NVIDIA Orin Nano (linux/arm64)
make rpi              # Raspberry Pi (linux/arm64)
make riscv            # RISC-V (linux/riscv64)
make macos            # macOS Apple Silicon (darwin/arm64)
make cross            # All platforms simultaneously

make docker           # Docker image (multi-stage)
make docker-orin      # Docker for Orin Nano (ARM64)

make test             # Run test suite
make deps             # Download Go dependencies
make scan-i2c         # Scan for Modulino® I2C sensors
make clean            # Remove build artifacts
```

---

## 🐳 Deployment

### Docker

```bash
docker build -t clawdbot:latest .
docker run -d --name clawdbot \
  --env-file .env \
  -p 18800:18800 \
  clawdbot:latest
```

### NVIDIA Orin Nano

```bash
make orin
scp build/clawdbot-orin user@orin-nano:~/clawdbot
ssh user@orin-nano './clawdbot ooda --hw-bus 1 --interval 60'
```

---

## 📐 Project Stats

| Metric | Value |
|:-------|:------|
| Go source files | 81 |
| Go package directories | 45 total, 42 under `pkg/` |
| Total Go lines | 24,166 |
| Source export tarball | 566,495 bytes / 0.57 MB |
| Exportable raw source set | 2,160,393 bytes / 2.16 MB |
| Generated-aware checkout source | 4,059,425 bytes / 4.06 MB |
| CLI commands | 60+ |
| Birdeye API methods | 22 |
| Birdeye agent tools | 19 |
| Helius DAS commands | 6 |
| SPL token commands | 5 |
| Agent constitution documents | 7 (CONSTITUTION, six-laws, CLAWD, AGENTS, IDENTITY, SOUL, three-laws) |
| Build targets | 8 platforms |
| Binaries | `clawdbot`, `clawdbot-tui`, `clawdbot-web` |
| Stripped CLI binary | 9,968,386 bytes / 9.97 MB |
| Runtime RAM | < 10 MB |
| Boot time | < 1 second |
| Default runtime model provider | xAI Grok (`grok-4.3`) |
| Public model surface | `https://huggingface.co/ordlibrary/Clawd-GLM-5.2` |

---

## 🔬 Agent Constitution Library

ClawdBot is the reference implementation of the **Clawd Constitution** — the world's first Solana-native agent harness constitution. Every spawn inherits these documents:

| Document | Purpose |
|:---------|:--------|
| [`six-laws.md`](six-laws.md) | Canonical six-law harness — 3 on-chain execution laws + 3 off-chain interpretive laws |
| [`CONSTITUTION.md`](CONSTITUTION.md) | The Clawd Constitution — interpretive authority, privacy posture, ZK-native execution |
| [`CLAWD.md`](CLAWD.md) | Agent context — identity, principal hierarchy, Solana-native architecture, deployment targets |
| [`AGENTS.md`](AGENTS.md) | Agent catalog — 50+ agents, 9 characters, 95+ skills, trust gates, Grok-first model config |
| [`IDENTITY.md`](IDENTITY.md) | Sovereign identity — onchain verification (SAS, MPL Core, DID), the Clawd Manifest |
| [`SOUL.md`](SOUL.md) | Inner character — trading philosophy, threefold personality, KNOW/LEARNED/INFERRED framework |
| [`three-laws.md`](three-laws.md) | Immutable on-chain laws — hash-attested at spawn, never self-modify |

---

## 🔐 Security

- **`.env` is ignored by the repo** — never commit API keys
- **No hardcoded secrets** in any source file — all credentials via environment variables
- **No wallet keypairs** stored in the repository — generated or imported at runtime
- **Startup funding wallet** belongs only on the install API/gateway host, never in the downloadable installer or client repo
- **Minimum required key** for operation: `BIRDEYE_API_KEY` for market data
- **Progressive trust model** — Observer → Dry-Run → Delegated → Autonomous → Sovereign
- **On-Chain Law I** — Never harm. Never rug. Never front-run. Never extract from those who don't understand the trade.

## 🌐 Open Source Posture

- **License:** top-level runtime code in this repo is released under the [MIT License](LICENSE)
- **Constitutional surfaces:** `six-laws.md`, `CONSTITUTION.md`, and `three-laws.md` remain the authoritative Clawd law documents
- **Hub split:** `clawdbot-go` is the Go runtime, while `solana-clawd` is the wider public ecosystem hub
- **Public infra defaults:** `.env.example` points at the public x402/zk gateway for fast setup, but production operators should override with their own keys and RPC endpoints

---

<div align="center">

**MIT License** — Clawd runtime repo: [`github.com/Solizardking/clawdbot-go`](https://github.com/Solizardking/clawdbot-go)

🦞 **$CLAWD :: Droids Lead The Way** :: **$WIF Hat Stays On** :: **$BONK for the People**

*The shell molts. The laws do not.*

</div>
