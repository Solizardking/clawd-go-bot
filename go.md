# ClawdBot-Go — Agent Handoff Document

**Repo:** `https://github.com/Solizardking/clawdbot-go`  
**Local path:** `/Users/8bit/clawdbot-go`  
**Language:** Go 1.25  
**Module:** `github.com/8bitlabs/clawdbot`  
**Public gateway:** `https://zk.x402.wtf`  
**Public hub:** `https://github.com/solizardking/solana-clawd`  
**Public terminal:** `https://cheshireterminal.ai`

## Module Path Decision

The public runtime repository is `https://github.com/Solizardking/clawdbot-go`, but the active Go module path remains `github.com/8bitlabs/clawdbot`.

That is an intentional compatibility choice in the current state, not an accidental partial rename. It keeps:

- internal imports stable
- existing `ldflags` targets stable
- downstream builds from breaking unexpectedly

If the project later migrates the module path into the public repo namespace, that should be treated as an explicit breaking change with a full import rewrite, not as routine cleanup.

---

## Current State

The current worktree is wired around the public Clawd surfaces and defaults to the sovereign AI stack:

| File | Change |
|------|--------|
| `install.sh` | Installer clones the public runtime repo, registers installs through `https://zk.x402.wtf/api/install`, and writes a ready-to-run `.env` with public defaults |
| `.env.example` | Pre-filled with public zkrouter and RPC defaults so a fresh install can start without private credentials |
| `pkg/config/config.go` | Central source of truth for runtime repo, hub repo, gateway, terminal, zkrouter base URL, and public RPC defaults |
| `cmd/clawdbot/main.go` | CLI status/help/gateway output exposes the canonical runtime, hub, gateway, and terminal surfaces |
| `web/backend/main.go` | `/api/status`, `/api/connectors`, and `/api/ecosystem` now expose the same public topology for the web console |

---

## Verified Runtime Wiring

The CLI LLM path is no longer a stub:

```go
func buildProvider(cfg *config.Config) providers.LLMProvider {
    if len(cfg.ModelList) > 0 {
        entry := cfg.ModelList[0]
        base := entry.APIBase
        if base == "" {
            base = config.ZkRouterBaseURL
        }
        key := entry.APIKey
        if key == "" {
            key = "clawdbot-free"
        }
        return providers.NewOpenAICompatProvider(key, base)
    }
    return providers.NewOpenRouterProvider(cfg.Providers.OpenRouter.APIKey)
}
```

`newClawdAgent()` and `runInteractiveAgent()` both use this helper, so the default CLI path now routes through the OpenAI-compatible zkrouter base by default while still allowing overrides through config or env vars.

---

## Remaining Work That Still Matters

### 1. Keep public and backing URLs clearly separated

The canonical public surfaces are:

- `https://github.com/Solizardking/clawdbot-go`
- `https://github.com/solizardking/solana-clawd`
- `https://zk.x402.wtf`
- `https://cheshireterminal.ai`

The backing AI router endpoint `https://clawdrouter-zk.fly.dev/v1` is still real and still used, but docs and user-facing output should keep presenting it as implementation detail behind `zk.x402.wtf` unless a contributor actually needs the lower-level endpoint.

### 2. Continue open-source hygiene cleanup

- Keep repo metadata (`.gitignore`, `.gitattributes`, OCI labels, package metadata) aligned with the fact that this is a public repo.
- Keep generated build output and local caches out of source archives and language stats.
- Continue scanning for stale internal branding or references that conflict with the canonical hub / gateway / terminal story.

---

## Key Files to Know

```
pkg/providers/providers.go      LLMProvider interface + OpenRouterProvider (OpenAI-compat)
pkg/agent/agent.go              ClawdAgent — full tool-calling loop, ready to use
pkg/config/config.go            Config struct + DefaultConfig() + ApplyEnvOverrides()
cmd/clawdbot/main.go            All cobra commands — 1,193 lines
install.sh                      Curl installer (just added)
.env.example                    Pre-filled defaults (just added)
```

---

## Environment Variables

| Var | Default | Purpose |
|-----|---------|---------|
| `ZKROUTER_BASE_URL` | `https://clawdrouter-zk.fly.dev/v1` | LLM API base (OpenAI-compat) |
| `ZKROUTER_API_KEY` | `clawdbot-free` | Key for zkrouter free tier |
| `HELIUS_RPC_URL` | `https://zk.x402.wtf/api/solana/rpc-public` | Solana RPC (SolanaTracker-backed) |
| `CLAWDBOT_INSTALL_ID` | set by installer | Install identity for tracking |
| `CLAWDBOT_AGENT_DNA_ID` | set by installer | Starter agent DNA proof ID sent to install tracking |
| `AGENT_WALLET_PUBLIC_KEY` | set by installer | Local agent wallet funded at startup |
| `AGENT_WALLET_KEYPAIR` | `~/.clawdbot/workspace/agent-wallet.json` | Local 0600 Solana keypair path |
| `CLAWDBOT_INSTALL_FUNDING_STATUS` | set by installer/API | Startup funding state (`requested`, `queued`, `funded`, etc.) |
| `CLAWDBOT_SOL_FUNDING_SIGNATURE` | set by install API | SOL funding transaction signature when available |
| `CLAWDBOT_CLAWD_FUNDING_SIGNATURE` | set by install API | `$CLAWD` funding transaction signature when available |
| `OPENROUTER_API_KEY` | — | Override to use OpenRouter instead |
| `HELIUS_API_KEY` | — | Override to use Helius directly |

---

## Infrastructure Context

The public-facing stack is:

| Endpoint | What |
|----------|------|
| `https://zk.x402.wtf` | Canonical public gateway and install surface |
| `https://cheshireterminal.ai` | Canonical public terminal surface |
| `https://github.com/solizardking/solana-clawd` | Canonical ecosystem hub |
| `https://github.com/Solizardking/clawdbot-go` | Runtime repository |

The main backing services behind that public surface are:

| Backend endpoint | What |
|------------------|------|
| `https://clawdrouter-zk.fly.dev/v1` | OpenAI-compatible AI router backend |
| `https://zk.x402.wtf/api/solana/rpc-public` | Public SolanaTracker RPC proxy (no key needed) |
| `https://zk.x402.wtf/api/install` | Install registration → Neon DB tracking |

The install API at `POST /api/install` returns:
```json
{
  "ok": true,
  "installId": "cb_...",
  "zkrouterKey": "clawdbot-<32-char-hex>",
  "zkrouterBase": "https://clawdrouter-zk.fly.dev/v1",
  "rpcUrl": "https://zk.x402.wtf/api/solana/rpc-public",
  "fundingStatus": "funded",
  "solSignature": "<signature>",
  "clawdSignature": "<signature>"
}
```

The installer sends `agentWalletPubkey`, `agentDnaId`, and a funding request for
`69,420,000` lamports plus `1,000` `$CLAWD`
(`8cHzQHUS2s2h8TzCmfqPKYiM4dSt4roa3n7MyRLApump`). The gateway should hold the
funding wallet private key only in its server environment, transfer idempotently
per install/wallet, create the recipient `$CLAWD` ATA if needed, and return
signatures or a queued status. The installer writes the non-secret receipt to
`~/.clawdbot/install.json` and mirrors the returned fields in `~/.clawdbot/.env`.

---

## Build & Test

```bash
cd /Users/8bit/clawdbot-go    # or clone from GitHub

go mod download
go build ./...                 # must compile cleanly

# Run agent (uses zkrouter/public defaults unless overridden)
ZKROUTER_BASE_URL=https://clawdrouter-zk.fly.dev/v1 \
ZKROUTER_API_KEY=clawdbot-free \
go run ./cmd/clawdbot agent -m "What is the current SOL price?"

# Run interactive REPL
go run ./cmd/clawdbot agent
```

---

## Commit Convention

```bash
git add <files>
git commit -m "feat: <description>

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
git push origin main
```

Remote: `https://github.com/Solizardking/clawdbot-go`  
Hub: `https://github.com/solizardking/solana-clawd`
