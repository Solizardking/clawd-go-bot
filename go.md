# ClawdBot-Go — Agent Handoff Document

**Repo:** `https://github.com/Solizardking/clawdbot-go`  
**Local path:** `/Users/8bit/clawdbot-go`  
**Language:** Go 1.25  
**Module:** `github.com/8bitlabs/clawdbot`  
**Public hub:** `https://github.com/solizardking/solana-clawd`  
**Public terminal:** `https://cheshireterminal.ai`

---

## What Was Just Done (last commit: `b54ccc2`)

Four files were added/modified to wire clawdbot into the sovereign AI stack:

| File | Change |
|------|--------|
| `install.sh` | Curl installer — clones repo, builds binary, writes `.env`, POSTs install telemetry to `https://zk.x402.wtf/api/install` to get an `installId` |
| `.env.example` | Pre-filled with free zkrouter AI + SolanaTracker RPC so users need zero API keys |
| `pkg/config/config.go` | Default `ModelList[0]` now points at `clawdrouter-zk.fly.dev/v1` with model `openai/zkrouter-auto`; default `HeliusRPCURL` points at `https://zk.x402.wtf/api/solana/rpc-public` |
| `README.md` | Quick-start section updated with curl install command |

---

## Critical Gap — LLM Is Not Actually Wired

The biggest pending work: **`runInteractiveAgent` is a stub.** The REPL loop at `cmd/clawdbot/main.go:1137` does not call any LLM. It just prints placeholder text:

```go
// cmd/clawdbot/main.go:1137
func runInteractiveAgent(cfg *config.Config) error {
    // ...
    default:
        fmt.Printf("[CLAWDBOT] Processing with %s...\n", cfg.Agents.Defaults.ModelName)
        fmt.Printf("(LLM integration pending — connect your API keys in config)\n")
    }
```

Similarly at `cmd/clawdbot/main.go:~120`:

```go
// TODO: Wire to LLM provider from config
_ = cfg
```

The `ClawdAgent` struct and the `providers.LLMProvider` interface are fully built (`pkg/agent/agent.go`, `pkg/providers/providers.go`) — they just aren't called from the CLI yet.

---

## What Needs to Be Done

### 1. Add a zkrouter-compatible provider constructor

`pkg/providers/providers.go` has `NewOpenRouterProvider(apiKey string)` which hardcodes `baseURL = "https://openrouter.ai/api/v1"`. Add an overload that accepts a custom base URL:

```go
// Add this function to pkg/providers/providers.go
func NewOpenAICompatProvider(apiKey, baseURL string) *OpenRouterProvider {
    return &OpenRouterProvider{
        apiKey:  apiKey,
        baseURL: baseURL,  // e.g. "https://clawdrouter-zk.fly.dev/v1"
        client:  &http.Client{Timeout: 120 * time.Second},
    }
}
```

The `OpenRouterProvider` already uses the OpenAI chat completions format (`/chat/completions`), so it works with zkrouter out of the box — just needs the configurable base URL.

### 2. Wire provider from config in `runInteractiveAgent`

Replace the stub in `cmd/clawdbot/main.go:runInteractiveAgent`:

```go
func runInteractiveAgent(cfg *config.Config) error {
    // Build provider from config.ModelList[0] (defaults to zkrouter)
    var provider providers.LLMProvider
    if len(cfg.ModelList) > 0 {
        entry := cfg.ModelList[0]
        base := entry.APIBase
        if base == "" {
            base = "https://clawdrouter-zk.fly.dev/v1"
        }
        provider = providers.NewOpenAICompatProvider(entry.APIKey, base)
    } else {
        provider = providers.NewOpenRouterProvider(cfg.Providers.OpenRouter.APIKey)
    }

    agent, err := agent.NewClawdAgent(agent.AgentConfig{
        Model:         cfg.ModelList[0].Model,  // "openai/zkrouter-auto"
        Provider:      provider,
        MaxIterations: cfg.Agents.Defaults.MaxToolIterations,
        MaxTokens:     cfg.Agents.Defaults.MaxTokens,
        Temperature:   cfg.Agents.Defaults.Temperature,
    })
    if err != nil {
        return err
    }

    // Then run the REPL loop using agent.Run(ctx, userInput)
    ...
}
```

Also fix the single-shot `-m` message path at `cmd/clawdbot/main.go:~120` — same provider construction, then call `agent.Run(ctx, message)`.

### 3. Wire the OODA loop provider

`clawdbot ooda` likely has the same stub. Find `runOODALoop` or equivalent and apply the same provider wiring.

### 4. Update README install URL

The README references:
```
git clone https://github.com/Solizardking/clawdbot-go.git
```
This was partially updated but double-check — all references should use:
```
https://github.com/Solizardking/clawdbot-go
```

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
| `OPENROUTER_API_KEY` | — | Override to use OpenRouter instead |
| `HELIUS_API_KEY` | — | Override to use Helius directly |

---

## Infrastructure Context

The zkrouter stack is deployed at:

| Endpoint | What |
|----------|------|
| `https://clawdrouter-zk.fly.dev/v1` | OpenAI-compatible AI router (fly.dev) |
| `https://zk.x402.wtf` | Next.js web frontend (Vercel) |
| `https://zk.x402.wtf/api/solana/rpc-public` | Public SolanaTracker RPC proxy (no key needed) |
| `https://zk.x402.wtf/api/install` | Install registration → Neon DB tracking |

The install API at `POST /api/install` returns:
```json
{
  "ok": true,
  "installId": "cb_...",
  "zkrouterKey": "clawdbot-<32-char-hex>",
  "zkrouterBase": "https://clawdrouter-zk.fly.dev/v1",
  "rpcUrl": "https://zk.x402.wtf/api/solana/rpc-public"
}
```

The installer already calls this and writes the returned key to `~/.clawdbot/.env`.

---

## Build & Test

```bash
cd /Users/8bit/clawdbot-go    # or clone from GitHub

go mod download
go build ./...                 # must compile cleanly

# Run agent (should hit zkrouter after your wiring changes)
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
