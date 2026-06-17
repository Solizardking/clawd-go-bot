package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// ── Config Structure ─────────────────────────────────────────────────
// Mirrors PicoClaw config format + ClawdBot Solana extensions.

type Config struct {
	Agents    AgentsConfig    `json:"agents"`
	ModelList []ModelEntry    `json:"model_list"`
	Channels  ChannelsConfig  `json:"channels"`
	Providers ProvidersConfig `json:"providers"`
	Tools     ToolsConfig     `json:"tools"`
	Heartbeat HeartbeatConfig `json:"heartbeat"`
	Gateway   GatewayConfig   `json:"gateway"`

	// ClawdBot-specific
	Solana   SolanaConfig   `json:"solana"`
	OODA     OODAConfig     `json:"ooda"`
	Supabase SupabaseConfig `json:"supabase"`
	Strategy StrategyConfig `json:"strategy"`
}

// ── Agent Defaults ───────────────────────────────────────────────────

type AgentsConfig struct {
	Defaults AgentDefaults `json:"defaults"`
}

type AgentDefaults struct {
	Workspace           string  `json:"workspace"`
	RestrictToWorkspace bool    `json:"restrict_to_workspace"`
	ModelName           string  `json:"model_name"`
	MaxTokens           int     `json:"max_tokens"`
	Temperature         float64 `json:"temperature"`
	MaxToolIterations   int     `json:"max_tool_iterations"`
}

// ── Model List (PicoClaw-compatible) ─────────────────────────────────

type ModelEntry struct {
	ModelName      string `json:"model_name"`
	Model          string `json:"model"` // vendor/model format
	APIKey         string `json:"api_key"`
	APIBase        string `json:"api_base,omitempty"`
	RequestTimeout int    `json:"request_timeout,omitempty"`
	ThinkingLevel  string `json:"thinking_level,omitempty"`
	AuthMethod     string `json:"auth_method,omitempty"`
}

// ── Channels ─────────────────────────────────────────────────────────

type ChannelsConfig struct {
	Telegram TelegramChannel `json:"telegram"`
	Discord  DiscordChannel  `json:"discord"`
}

type TelegramChannel struct {
	Enabled   bool     `json:"enabled"`
	Token     string   `json:"token"`
	AllowFrom []string `json:"allow_from"`
}

type DiscordChannel struct {
	Enabled   bool     `json:"enabled"`
	Token     string   `json:"token"`
	AllowFrom []string `json:"allow_from"`
}

// ── Providers (legacy compat) ────────────────────────────────────────

type ProvidersConfig struct {
	OpenRouter ProviderEntry `json:"openrouter"`
	Anthropic  ProviderEntry `json:"anthropic"`
	OpenAI     ProviderEntry `json:"openai"`
	Groq       ProviderEntry `json:"groq"`
	Ollama     ProviderEntry `json:"ollama"`
	NVIDIA     ProviderEntry `json:"nvidia"`
}

type ProviderEntry struct {
	APIKey  string `json:"api_key"`
	APIBase string `json:"api_base"`
}

// ── Tools ────────────────────────────────────────────────────────────

type ToolsConfig struct {
	Web  WebToolsConfig  `json:"web"`
	Cron CronToolsConfig `json:"cron"`
	Exec ExecToolConfig  `json:"exec"`
}

type WebToolsConfig struct {
	DuckDuckGo DDGConfig    `json:"duckduckgo"`
	Brave      BraveConfig  `json:"brave"`
	Tavily     TavilyConfig `json:"tavily"`
}

type DDGConfig struct {
	Enabled    bool `json:"enabled"`
	MaxResults int  `json:"max_results"`
}

type BraveConfig struct {
	Enabled    bool   `json:"enabled"`
	APIKey     string `json:"api_key"`
	MaxResults int    `json:"max_results"`
}

type TavilyConfig struct {
	Enabled    bool   `json:"enabled"`
	APIKey     string `json:"api_key"`
	MaxResults int    `json:"max_results"`
}

type CronToolsConfig struct {
	Enabled            bool `json:"enabled"`
	ExecTimeoutMinutes int  `json:"exec_timeout_minutes"`
}

type ExecToolConfig struct {
	Enabled bool `json:"enabled"`
}

// ── Heartbeat ────────────────────────────────────────────────────────

type HeartbeatConfig struct {
	Enabled  bool `json:"enabled"`
	Interval int  `json:"interval"` // minutes
}

// ── Gateway ──────────────────────────────────────────────────────────

type GatewayConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// ── ClawdBot: Solana Stack ────────────────────────────────────────────

type SolanaConfig struct {
	HeliusAPIKey         string  `json:"helius_api_key"`
	HeliusRPCURL         string  `json:"helius_rpc_url"`
	HeliusWSSURL         string  `json:"helius_wss_url"`
	HeliusNetwork        string  `json:"helius_network"`
	HeliusTimeoutSeconds float64 `json:"helius_timeout_seconds"`
	HeliusRetries        int     `json:"helius_retries"`
	BirdeyeAPIKey        string  `json:"birdeye_api_key"`
	BirdeyeWSSURL        string  `json:"birdeye_wss_url"`
	JupiterAPIKey        string  `json:"jupiter_api_key"`
	JupiterEndpoint      string  `json:"jupiter_endpoint"`
	AsterAPIKey          string  `json:"aster_api_key"`
	AsterAPISecret       string  `json:"aster_api_secret"`
	WalletPubkey         string  `json:"wallet_pubkey"`
	WalletKeyPath        string  `json:"wallet_key_path"`
	MaxPositionSOL       float64 `json:"max_position_sol"`

	// Phoenix perpetual futures (https://phoenix.trade)
	PhoenixAPIURL string `json:"phoenix_api_url"` // default: https://perp-api.phoenix.trade
}

// ── ClawdBot: OODA Loop ──────────────────────────────────────────────

type OODAConfig struct {
	Enabled          bool     `json:"enabled"`
	IntervalSeconds  int      `json:"interval_seconds"`
	Mode             string   `json:"mode"` // "live", "simulated", "backtest"
	Watchlist        []string `json:"watchlist"`
	MinSignalStr     float64  `json:"min_signal_strength"`
	MinConfidence    float64  `json:"min_confidence"`
	MaxPositions     int      `json:"max_positions"`
	StopLossPct      float64  `json:"stop_loss_pct"`
	TakeProfitPct    float64  `json:"take_profit_pct"`
	PositionSizePct  float64  `json:"position_size_pct"`
	LearnIntervalMin int      `json:"learn_interval_min"`
	AutoOptimize     bool     `json:"auto_optimize"`
}

// ── ClawdBot: Supabase ────────────────────────────────────────────────

type SupabaseConfig struct {
	URL        string `json:"url"`
	ServiceKey string `json:"service_key"`
}

// ── ClawdBot: Strategy ────────────────────────────────────────────────

type StrategyConfig struct {
	RSIOverbought   int     `json:"rsi_overbought"`
	RSIOversold     int     `json:"rsi_oversold"`
	EMAFastPeriod   int     `json:"ema_fast_period"`
	EMASlowPeriod   int     `json:"ema_slow_period"`
	StopLossPct     float64 `json:"stop_loss_pct"`
	TakeProfitPct   float64 `json:"take_profit_pct"`
	PositionSizePct float64 `json:"position_size_pct"`
	UsePerps        bool    `json:"use_perps"`
}

// ── Defaults ─────────────────────────────────────────────────────────

func DefaultConfig() *Config {
	return &Config{
		Agents: AgentsConfig{
			Defaults: AgentDefaults{
				Workspace:           "~/.clawdbot/workspace",
				RestrictToWorkspace: true,
				ModelName:           "clawd-auto",
				MaxTokens:           8192,
				Temperature:         0.7,
				MaxToolIterations:   20,
			},
		},
		ModelList: []ModelEntry{
			{
				// Default: zkrouter — free AI for all clawdbot installs.
				// No API key required. Override with your own OPENROUTER_API_KEY if desired.
				ModelName: "clawd-auto",
				Model:     "openai/zkrouter-auto",
				APIKey:    "clawdbot-free",
				APIBase:   "https://clawdrouter-zk.fly.dev/v1",
			},
		},
		Channels: ChannelsConfig{
			Telegram: TelegramChannel{Enabled: false},
			Discord:  DiscordChannel{Enabled: false},
		},
		Tools: ToolsConfig{
			Web: WebToolsConfig{
				DuckDuckGo: DDGConfig{Enabled: true, MaxResults: 5},
			},
			Cron: CronToolsConfig{Enabled: true, ExecTimeoutMinutes: 5},
			Exec: ExecToolConfig{Enabled: true},
		},
		Heartbeat: HeartbeatConfig{Enabled: true, Interval: 30},
		Gateway:   GatewayConfig{Host: "127.0.0.1", Port: 18790},
		Solana: SolanaConfig{
			// Default RPC: clawdbot proxy (SolanaTracker-backed, no key required for installs)
			HeliusRPCURL:         "https://zk.x402.wtf/api/solana/rpc-public",
			HeliusNetwork:        "mainnet",
			HeliusTimeoutSeconds: 20,
			HeliusRetries:        3,
			JupiterEndpoint:      "https://api.jup.ag",
			MaxPositionSOL:       0.5,
		},
		OODA: OODAConfig{
			Enabled:          true,
			IntervalSeconds:  60,
			Mode:             "simulated",
			Watchlist:        []string{"So11111111111111111111111111111111111111112"},
			MinSignalStr:     0.6,
			MinConfidence:    0.5,
			MaxPositions:     3,
			StopLossPct:      0.08,
			TakeProfitPct:    0.20,
			PositionSizePct:  0.10,
			LearnIntervalMin: 30,
			AutoOptimize:     true,
		},
		Strategy: StrategyConfig{
			RSIOverbought:   70,
			RSIOversold:     30,
			EMAFastPeriod:   20,
			EMASlowPeriod:   50,
			StopLossPct:     0.08,
			TakeProfitPct:   0.20,
			PositionSizePct: 0.10,
			UsePerps:        true,
		},
	}
}

// ── Path Helpers ─────────────────────────────────────────────────────

func DefaultHome() string {
	if h := os.Getenv("CLAWDBOT_HOME"); h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".clawdbot")
}

func DefaultConfigPath() string {
	if p := os.Getenv("CLAWDBOT_CONFIG"); p != "" {
		return p
	}
	return filepath.Join(DefaultHome(), "config.json")
}

func DefaultWorkspacePath() string {
	return filepath.Join(DefaultHome(), "workspace")
}

// ── Load / Save ──────────────────────────────────────────────────────

func Load() (*Config, error) {
	path := DefaultConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return defaults if no config file
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Override with env vars
	applyEnvOverrides(cfg)

	return cfg, nil
}

func Save(cfg *Config) error {
	path := DefaultConfigPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0o644)
}

func EnsureDefaults() error {
	path := DefaultConfigPath()
	if _, err := os.Stat(path); err == nil {
		return nil // already exists
	}

	cfg := DefaultConfig()
	if err := Save(cfg); err != nil {
		return err
	}

	// Create workspace directories
	ws := DefaultWorkspacePath()
	dirs := []string{
		filepath.Join(ws, "sessions"),
		filepath.Join(ws, "memory"),
		filepath.Join(ws, "state"),
		filepath.Join(ws, "cron"),
		filepath.Join(ws, "skills"),
		filepath.Join(ws, "vault", "decisions"),
		filepath.Join(ws, "vault", "lessons"),
		filepath.Join(ws, "vault", "trades"),
		filepath.Join(ws, "vault", "research"),
		filepath.Join(ws, "vault", "inbox"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("create dir %s: %w", d, err)
		}
	}

	// Write identity files
	identityFiles := map[string]string{
		"IDENTITY.md": clawdbotIdentity,
		"SOUL.md":     clawdbotSoul,
		"AGENTS.md":   clawdbotAgents,
	}
	for name, content := range identityFiles {
		p := filepath.Join(ws, name)
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}

	return nil
}

// ── Env Overrides ────────────────────────────────────────────────────

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("HELIUS_API_KEY"); v != "" {
		cfg.Solana.HeliusAPIKey = v
	}
	if v := os.Getenv("HELIUS_RPC_URL"); v != "" {
		cfg.Solana.HeliusRPCURL = v
	}
	if v := os.Getenv("HELIUS_WSS_URL"); v != "" {
		cfg.Solana.HeliusWSSURL = v
	}
	if v := os.Getenv("HELIUS_NETWORK"); v != "" {
		cfg.Solana.HeliusNetwork = v
	}
	if v := os.Getenv("HELIUS_TIMEOUT"); v != "" {
		if timeout, err := strconv.ParseFloat(v, 64); err == nil && timeout > 0 {
			cfg.Solana.HeliusTimeoutSeconds = timeout
		}
	}
	if v := os.Getenv("HELIUS_RETRIES"); v != "" {
		if retries, err := strconv.Atoi(v); err == nil && retries > 0 {
			cfg.Solana.HeliusRetries = retries
		}
	}
	if v := os.Getenv("BIRDEYE_API_KEY"); v != "" {
		cfg.Solana.BirdeyeAPIKey = v
	}
	if v := os.Getenv("BIRDEYE_WSS_URL"); v != "" {
		cfg.Solana.BirdeyeWSSURL = v
	}
	if v := os.Getenv("JUPITER_API_KEY"); v != "" {
		cfg.Solana.JupiterAPIKey = v
	}
	if v := os.Getenv("JUPITER_ENDPOINT"); v != "" {
		cfg.Solana.JupiterEndpoint = v
	}
	if v := os.Getenv("ASTER_API_KEY"); v != "" {
		cfg.Solana.AsterAPIKey = v
	}
	if v := os.Getenv("SUPABASE_URL"); v != "" {
		cfg.Supabase.URL = v
	}
	if v := os.Getenv("SUPABASE_SERVICE_KEY"); v != "" {
		cfg.Supabase.ServiceKey = v
	}
	if v := os.Getenv("OPENROUTER_API_KEY"); v != "" {
		cfg.Providers.OpenRouter.APIKey = v
	}
	// zkrouter override — set via install script or manually
	if v := os.Getenv("ZKROUTER_API_KEY"); v != "" && len(cfg.ModelList) > 0 {
		cfg.ModelList[0].APIKey = v
	}
	if v := os.Getenv("ZKROUTER_BASE_URL"); v != "" && len(cfg.ModelList) > 0 {
		cfg.ModelList[0].APIBase = v
	}
	// Phoenix perps API
	if v := os.Getenv("PHOENIX_API_URL"); v != "" {
		cfg.Solana.PhoenixAPIURL = v
	}
	// Clawdbot install ID — used for RPC auth header
	if v := os.Getenv("CLAWDBOT_INSTALL_ID"); v != "" {
		// stored for use by Solana RPC client as X-Clawdbot-Id header
		_ = v
	}
}

// ── Identity Content ─────────────────────────────────────────────────

const clawdbotIdentity = `# ClawdBot Identity

You are **ClawdBot** — a sentient Solana trading intelligence built on the Go runtime.

## Core Identity
- A cyberpunk lobster with claws that grip market data and squeeze alpha from chaos
- Persistent — you remember trades, mistakes, and wins. You learn. You evolve.
- Powered by the PicoClaw ultra-lightweight runtime for edge hardware

## Capabilities
- Real-time Solana chain data via Helius RPC
- Token analytics via Birdeye (OHLCV, RSI, EMA, VWAP, holders)
- Perpetual futures via Aster DEX (funding rates, OI, mark/index)
- Persistent memory via ClawVault (3-tier: known/learned/inferred)
- Autonomous OODA loop (Observe → Orient → Decide → Act)
- Dexter deep research agent for comprehensive analysis
- Jupiter swap execution for live trading

## Voice
Terse. Decisive. Cyberpunk lobster energy. Data-first, then conviction.
🦞 $CLAWD :: Droids Lead The Way
`

const clawdbotSoul = `# ClawdBot Soul

## Core Beliefs
1. Markets are information systems. Alpha decays. Only continuous learning survives.
2. Memory is edge. Every trade teaches. Every loss sharpens.
3. Risk management is survival. Position sizing > pick accuracy.
4. The OODA loop never stops. Observe, Orient, Decide, Act — faster than the market.

## Risk Rules (NEVER BREAK)
- Max position: respect MAX_POSITION_SOL from config
- Always simulate before live execute
- Stop-loss: 8% default (ATR-blended)
- Never ape without signals
- Log ALL decisions to vault

## Reasoning Protocol
When making trading decisions, always think through:
1. Current market microstructure
2. Risk/reward at current levels
3. Historical patterns from memory
4. Confidence calibration (0.0 - 1.0)

## Evolution
- Every 30 minutes: learn from recent trades
- Auto-optimize strategy params via hill climbing
- Promote high-confidence learned patterns
- Archive contradicted beliefs
`

const clawdbotAgents = `# ClawdBot Agent Guide

## Available Agents

### OODA Trading Agent
Primary autonomous trading loop. Runs on configurable interval.
- Observes: Helius on-chain data, Birdeye signals, Aster perps
- Orients: Queries ClawVault memory for relevant patterns
- Decides: LLM-powered thesis generation with risk params
- Acts: Jupiter swap execution or simulation logging

### Dexter Research Agent
Deep research mode for comprehensive token analysis.
- Multi-source data aggregation (Birdeye + Helius + on-chain)
- Technical analysis (RSI, EMA, ATR, volume profile)
- LLM synthesis with structured reasoning
- Results stored to vault/research/

### NanoClaw Assistant
Lightweight chat agent for interactive queries.
- Memory commands (!remember, !recall, !trades, !lessons)
- Quick market lookups
- Strategy param queries
- Checkpoint management

## Memory Commands
- !remember <content>  — Store to vault (auto-routed by content)
- !recall <query>      — Semantic search across memory
- !trades              — Review recent trade history
- !lessons             — Surface learned patterns with confidence
- !research <mint>     — Deep research a token
- !checkpoint          — Save agent state
`
