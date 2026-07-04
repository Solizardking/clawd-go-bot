// ClawdBot Web Console — web-based dashboard and agent control.
// ClawdBot Go web launcher — serves embedded frontend,
// provides API for config management and gateway control.
//
// Usage:
//   go build -o clawdbot-web ./web/backend/
//   ./clawdbot-web [config.json]
//   ./clawdbot-web -public config.json

package main

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/8bitlabs/clawdbot/pkg/birthfund"
	"github.com/8bitlabs/clawdbot/pkg/config"
	dnaPkg "github.com/8bitlabs/clawdbot/pkg/dna"
	"github.com/8bitlabs/clawdbot/pkg/doctor"
	"github.com/8bitlabs/clawdbot/pkg/laws"
	"github.com/8bitlabs/clawdbot/pkg/strategy"
	"github.com/8bitlabs/clawdbot/pkg/trading"
	"github.com/8bitlabs/clawdbot/pkg/wallet"
)

const banner = `
  ╔══════════════════════════════════════════════╗
  ║       🦞 ClawdBot OS — Web Console           ║
  ║  Sovereign Solana Trading Intelligence       ║
  ╚══════════════════════════════════════════════╝`

func main() {
	port := flag.String("port", "18800", "Port to listen on")
	public := flag.Bool("public", false, "Listen on all interfaces (0.0.0.0) instead of localhost only")
	noBrowser := flag.Bool("no-browser", false, "Do not auto-open browser on startup")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "ClawdBot Web Console — Dashboard and agent control\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [config.json]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	configPath := defaultConfigPath()
	if flag.NArg() > 0 {
		configPath = flag.Arg(0)
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		log.Fatalf("Config path error: %v", err)
	}

	portNum, err := strconv.Atoi(*port)
	if err != nil || portNum < 1 || portNum > 65535 {
		log.Fatalf("Invalid port: %s", *port)
	}

	var addr string
	if *public {
		addr = "0.0.0.0:" + *port
	} else {
		addr = "127.0.0.1:" + *port
	}

	// Determine project root (directory containing go.mod)
	projectRoot := findProjectRoot(absPath)

	// API routes
	mux := http.NewServeMux()

	// API: Status
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":       "running",
			"version":      "1.0.0",
			"agent":        "ClawdBot Go",
			"config":       absPath,
			"uptime":       time.Since(startTime).String(),
			"mode":         os.Getenv("AGENT_MODE"),
			"go_version":   runtime.Version(),
			"go_os":        runtime.GOOS,
			"go_arch":      runtime.GOARCH,
			"num_cpu":      runtime.NumCPU(),
			"goroutines":   runtime.NumGoroutine(),
			"dna_path":     dnaPkg.DefaultPath(config.DefaultWorkspacePath()),
			"public_links": ecosystemLinks(),
		})
	})

	mux.HandleFunc("/api/dna", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := dnaPkg.DefaultPath(config.DefaultWorkspacePath())
		value, created, err := dnaPkg.EnsureFile(path, dnaPkg.Options{
			AgentName: "ClawdBot",
			Role:      "sovereign Solana trading intelligence",
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"path":    path,
			"created": created,
			"dna":     value,
		})
	})

	// API: Config read
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if webEnvBool("CLAWDBOT_WEB_EXPOSE_SECRETS") {
			data, err := os.ReadFile(absPath)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Write(data)
			return
		}
		cfg, err := loadRuntimeConfig(absPath)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(redactedConfig(cfg))
	})

	// API: Health
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"agent":  "clawdbot-go",
		})
	})

	mux.HandleFunc("/api/install", installAPIHandler())
	mux.HandleFunc("/api/installs", installsAPIHandler())

	// API: Connectors status
	mux.HandleFunc("/api/connectors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		connectors := []map[string]any{
			{"name": "x402 Gateway", "status": urlStatus(os.Getenv("ZKROUTER_BASE_URL"), config.ZkRouterBaseURL), "type": "gateway"},
			{"name": "Clawd Terminal", "status": "public", "type": "terminal"},
			{"name": "Helius", "status": envStatus("HELIUS_API_KEY"), "type": "rpc"},
			{"name": "Birdeye", "status": envStatus("BIRDEYE_API_KEY"), "type": "analytics"},
			{"name": "Jupiter", "status": envStatus("JUPITER_API_KEY"), "type": "swap"},
			{"name": "Aster", "status": envStatus("ASTER_API_KEY"), "type": "perps"},
			{"name": "Vulcan", "status": binaryStatus("vulcan"), "type": "perps_cli"},
			{"name": "OpenRouter", "status": envStatus("OPENROUTER_API_KEY"), "type": "llm"},
			{"name": "Supabase", "status": envStatus("SUPABASE_URL"), "type": "database"},
		}
		json.NewEncoder(w).Encode(connectors)
	})

	mux.HandleFunc("/api/ecosystem", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"runtime_repo": config.RuntimeRepoURL,
			"hub_repo":     config.HubRepoURL,
			"gateway":      config.GatewayURL,
			"terminal":     config.TerminalURL,
		})
	})

	mux.HandleFunc("/api/laws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(laws.Six)
	})

	mux.HandleFunc("/api/trading/cockpit", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cfg, err := loadRuntimeConfig(absPath)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(trading.BuildCockpitReport(cfg, time.Now()))
	})

	// API: Live strategy signal — runs the real strategy engine over a series of
	// closes/highs/lows (POST) or a deterministic demo series (GET), so the web
	// console exercises the same Evaluate() the OODA loop uses.
	mux.HandleFunc("/api/trading/signal", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cfg, err := loadRuntimeConfig(absPath)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		params := strategyParamsFromConfig(cfg)
		closes, highs, lows := demoSeries()
		if r.Method == http.MethodPost {
			var body struct {
				Closes []float64 `json:"closes"`
				Highs  []float64 `json:"highs"`
				Lows   []float64 `json:"lows"`
			}
			if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&body); err != nil {
				http.Error(w, "invalid signal payload", http.StatusBadRequest)
				return
			}
			if len(body.Closes) > 0 {
				closes = body.Closes
				highs = body.Highs
				lows = body.Lows
				if len(highs) != len(closes) || len(lows) != len(closes) {
					highs, lows = deriveHighsLows(closes)
				}
			}
		}
		sig := strategy.Evaluate(closes, highs, lows, params)
		json.NewEncoder(w).Encode(map[string]any{
			"params":  params,
			"signal":  sig,
			"samples": len(closes),
		})
	})

	// API: Backtest — replays the strategy over a demo (GET) or posted series and
	// returns win rate, drawdown, Sharpe, and the equity curve.
	mux.HandleFunc("/api/trading/backtest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cfg, err := loadRuntimeConfig(absPath)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		params := strategyParamsFromConfig(cfg)
		bars := demoBars(400)
		if r.Method == http.MethodPost {
			var body struct {
				Bars []strategy.Bar `json:"bars"`
			}
			if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<20)).Decode(&body); err != nil {
				http.Error(w, "invalid backtest payload", http.StatusBadRequest)
				return
			}
			if len(body.Bars) > 0 {
				bars = body.Bars
			}
		}
		result := strategy.Backtest(bars, params, params.EMASlowPeriod+5)
		json.NewEncoder(w).Encode(map[string]any{
			"params": params,
			"bars":   len(bars),
			"result": result,
		})
	})

	// API: Portfolio guard — evaluates the account-level risk gate against the
	// current config limits and a posted (or zero) exposure snapshot.
	mux.HandleFunc("/api/trading/portfolio", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cfg, err := loadRuntimeConfig(absPath)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		limits := portfolioLimitsFromConfig(cfg)
		var body struct {
			Asset    string               `json:"asset"`
			SizeSOL  float64              `json:"sizeSol"`
			Exposure trading.OpenExposure `json:"exposure"`
		}
		if r.Method == http.MethodPost {
			_ = json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&body)
		}
		if body.SizeSOL <= 0 {
			body.SizeSOL = cfg.Solana.MaxPositionSOL
		}
		result := limits.CheckEntry(body.Asset, body.SizeSOL, body.Exposure)
		json.NewEncoder(w).Encode(map[string]any{
			"limits":   limits,
			"exposure": body.Exposure,
			"candidate": map[string]any{
				"asset":   body.Asset,
				"sizeSol": body.SizeSOL,
			},
			"guard": result,
		})
	})

	mux.HandleFunc("/api/doctor", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cfg, err := loadRuntimeConfig(absPath)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(doctor.Run(doctor.Options{
			Config:        cfg,
			ConfigPath:    absPath,
			WorkspacePath: config.DefaultWorkspacePath(),
			ProjectRoot:   projectRoot,
		}))
	})

	// API: Packages — list all Go packages in the project
	mux.HandleFunc("/api/packages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		pkgs := listGoPackages(projectRoot)
		json.NewEncoder(w).Encode(pkgs)
	})

	// API: Environment variables (safe, non-secret subset)
	mux.HandleFunc("/api/env", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		env := map[string]string{
			"AGENT_MODE": os.Getenv("AGENT_MODE"),
			"HOSTNAME":   os.Getenv("HOSTNAME"),
			"PWD":        os.Getenv("PWD"),
			"SHELL":      os.Getenv("SHELL"),
		}
		json.NewEncoder(w).Encode(env)
	})

	// Serve embedded frontend (or static files)
	frontendDir := filepath.Join(filepath.Dir(absPath), "web", "frontend", "dist")
	if _, err := os.Stat(frontendDir); err == nil {
		mux.Handle("/", http.FileServer(http.Dir(frontendDir)))
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fallbackHTML))
		})
	}

	// CORS middleware
	handler := corsMiddleware(loggerMiddleware(mux))

	// Startup
	fmt.Print(banner)
	fmt.Println()
	fmt.Printf("  Config: %s\n", absPath)
	fmt.Printf("  Project: %s\n", projectRoot)
	fmt.Printf("  Open: http://localhost:%s\n", *port)
	if *public {
		if ip := getLocalIP(); ip != "" {
			fmt.Printf("  Public: http://%s:%s\n", ip, *port)
		}
	}
	fmt.Println()

	if !*noBrowser {
		go func() {
			time.Sleep(500 * time.Millisecond)
			openBrowser("http://localhost:" + *port)
		}()
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       time.Duration(envInt("CLAWDBOT_WEB_READ_TIMEOUT_SECONDS", 15)) * time.Second,
		WriteTimeout:      time.Duration(envInt("CLAWDBOT_WEB_WRITE_TIMEOUT_SECONDS", 300)) * time.Second,
		IdleTimeout:       time.Duration(envInt("CLAWDBOT_WEB_IDLE_TIMEOUT_SECONDS", 120)) * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

var startTime = time.Now()

func defaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".clawdbot", "config.json")
}

func envStatus(key string) string {
	if os.Getenv(key) != "" {
		return "connected"
	}
	return "not_configured"
}

func binaryStatus(name string) string {
	if _, err := exec.LookPath(name); err == nil {
		return "connected"
	}
	return "not_configured"
}

func urlStatus(value, expected string) string {
	if strings.TrimSpace(value) == "" {
		return "default_public"
	}
	if strings.TrimSpace(value) == expected {
		return "default_public"
	}
	return "custom"
}

// strategyParamsFromConfig maps the runtime strategy config into engine params.
func strategyParamsFromConfig(cfg *config.Config) strategy.StrategyParams {
	return strategy.StrategyParams{
		RSIOverbought:   cfg.Strategy.RSIOverbought,
		RSIOversold:     cfg.Strategy.RSIOversold,
		EMAFastPeriod:   cfg.Strategy.EMAFastPeriod,
		EMASlowPeriod:   cfg.Strategy.EMASlowPeriod,
		StopLossPct:     cfg.Strategy.StopLossPct,
		TakeProfitPct:   cfg.Strategy.TakeProfitPct,
		PositionSizePct: cfg.Strategy.PositionSizePct,
		UsePerps:        cfg.Strategy.UsePerps,
	}
}

// portfolioLimitsFromConfig derives account-level guardrails from config, using
// conservative defaults for limits the config does not yet express.
func portfolioLimitsFromConfig(cfg *config.Config) trading.PortfolioLimits {
	maxConcurrent := cfg.OODA.MaxPositions
	if maxConcurrent <= 0 {
		maxConcurrent = 3
	}
	totalExposure := cfg.Solana.MaxPositionSOL * float64(maxConcurrent)
	return trading.PortfolioLimits{
		MaxConcurrent:     maxConcurrent,
		MaxTotalExposure:  totalExposure,
		MaxPerAsset:       cfg.Solana.MaxPositionSOL,
		MaxDrawdownPct:    0.25,
		DailyLossLimitPct: 0.10,
	}
}

// demoSeries returns a deterministic sine-wave price series so the signal
// endpoint has something to evaluate without live market data.
func demoSeries() (closes, highs, lows []float64) {
	closes = make([]float64, 120)
	for i := range closes {
		closes[i] = 100 + 15*math.Sin(float64(i)/6.0) + 0.05*float64(i)
	}
	highs, lows = deriveHighsLows(closes)
	return closes, highs, lows
}

// deriveHighsLows synthesizes a ±2% intrabar range around each close.
func deriveHighsLows(closes []float64) (highs, lows []float64) {
	highs = make([]float64, len(closes))
	lows = make([]float64, len(closes))
	for i, c := range closes {
		highs[i] = c * 1.02
		lows[i] = c * 0.98
	}
	return highs, lows
}

// demoBars returns a deterministic OHLCV series for the backtest endpoint.
func demoBars(n int) []strategy.Bar {
	bars := make([]strategy.Bar, n)
	for i := range bars {
		base := 100 + 15*math.Sin(float64(i)/6.0) + 0.05*float64(i)
		bars[i] = strategy.Bar{Close: base, High: base * 1.02, Low: base * 0.98}
	}
	return bars
}

func loadRuntimeConfig(path string) (*config.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config.DefaultConfig(), nil
		}
		return nil, err
	}
	cfg := config.DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func redactedConfig(cfg *config.Config) config.Config {
	if cfg == nil {
		return *config.DefaultConfig()
	}
	out := *cfg
	out.ModelList = append([]config.ModelEntry(nil), cfg.ModelList...)
	for i := range out.ModelList {
		out.ModelList[i].APIKey = redactSecret(out.ModelList[i].APIKey)
	}
	out.Channels.Telegram.Token = redactSecret(out.Channels.Telegram.Token)
	out.Channels.Discord.Token = redactSecret(out.Channels.Discord.Token)
	out.Providers.OpenRouter.APIKey = redactSecret(out.Providers.OpenRouter.APIKey)
	out.Providers.Anthropic.APIKey = redactSecret(out.Providers.Anthropic.APIKey)
	out.Providers.OpenAI.APIKey = redactSecret(out.Providers.OpenAI.APIKey)
	out.Providers.Groq.APIKey = redactSecret(out.Providers.Groq.APIKey)
	out.Providers.Ollama.APIKey = redactSecret(out.Providers.Ollama.APIKey)
	out.Providers.NVIDIA.APIKey = redactSecret(out.Providers.NVIDIA.APIKey)
	out.Solana.HeliusAPIKey = redactSecret(out.Solana.HeliusAPIKey)
	out.Solana.BirdeyeAPIKey = redactSecret(out.Solana.BirdeyeAPIKey)
	out.Solana.JupiterAPIKey = redactSecret(out.Solana.JupiterAPIKey)
	out.Solana.AsterAPIKey = redactSecret(out.Solana.AsterAPIKey)
	out.Solana.AsterAPISecret = redactSecret(out.Solana.AsterAPISecret)
	out.Solana.WalletKeyPath = redactSecret(out.Solana.WalletKeyPath)
	out.Supabase.ServiceKey = redactSecret(out.Supabase.ServiceKey)
	return out
}

func redactSecret(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return "<redacted>"
}

func ecosystemLinks() map[string]string {
	return map[string]string{
		"runtime_repo": config.RuntimeRepoURL,
		"hub_repo":     config.HubRepoURL,
		"gateway":      config.GatewayURL,
		"terminal":     config.TerminalURL,
	}
}

var installLedgerMu sync.Mutex

type installFundingRequest struct {
	SOLLamports    uint64      `json:"solLamports"`
	CLAWDTokens    json.Number `json:"clawdTokens"`
	CLAWDMint      string      `json:"clawdMint"`
	CreateCLAWDATA bool        `json:"createClawdAta"`
}

type installRequest struct {
	InstallID         string                `json:"installId"`
	OS                string                `json:"os"`
	Arch              string                `json:"arch"`
	Version           string                `json:"version"`
	InstallComplete   string                `json:"installComplete"`
	CoreAI            string                `json:"coreAi"`
	Vulcan            string                `json:"vulcan"`
	AgentWalletPubkey string                `json:"agentWalletPubkey"`
	AgentDNAID        string                `json:"agentDnaId"`
	Funding           installFundingRequest `json:"funding"`
}

type installRecord struct {
	InstallID         string            `json:"installId"`
	RemoteIP          string            `json:"remoteIp"`
	UserAgent         string            `json:"userAgent"`
	OS                string            `json:"os,omitempty"`
	Arch              string            `json:"arch,omitempty"`
	Version           string            `json:"version,omitempty"`
	InstallComplete   string            `json:"installComplete,omitempty"`
	CoreAI            string            `json:"coreAi,omitempty"`
	Vulcan            string            `json:"vulcan,omitempty"`
	AgentWalletPubkey string            `json:"agentWalletPubkey,omitempty"`
	AgentDNAID        string            `json:"agentDnaId,omitempty"`
	FundingStatus     string            `json:"fundingStatus"`
	FundingError      string            `json:"fundingError,omitempty"`
	Funding           *birthfund.Result `json:"funding,omitempty"`
	CreatedAt         string            `json:"createdAt"`
}

func installAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req installRequest
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
		dec.UseNumber()
		if err := dec.Decode(&req); err != nil {
			http.Error(w, "invalid install payload", http.StatusBadRequest)
			return
		}

		installID := strings.TrimSpace(req.InstallID)
		if installID == "" {
			installID = randomInstallID()
		}
		ledgerPath := installLedgerPath()
		recipient := strings.TrimSpace(req.AgentWalletPubkey)
		record := installRecord{
			InstallID:         installID,
			RemoteIP:          clientIP(r),
			UserAgent:         r.UserAgent(),
			OS:                strings.TrimSpace(req.OS),
			Arch:              strings.TrimSpace(req.Arch),
			Version:           strings.TrimSpace(req.Version),
			InstallComplete:   strings.TrimSpace(req.InstallComplete),
			CoreAI:            strings.TrimSpace(req.CoreAI),
			Vulcan:            strings.TrimSpace(req.Vulcan),
			AgentWalletPubkey: recipient,
			AgentDNAID:        strings.TrimSpace(req.AgentDNAID),
			FundingStatus:     "skipped",
			CreatedAt:         time.Now().UTC().Format(time.RFC3339),
		}

		resp := map[string]any{
			"ok":            true,
			"installId":     installID,
			"zkrouterKey":   firstNonEmptyEnv("ZKROUTER_API_KEY", "clawdbot-free"),
			"zkrouterBase":  firstNonEmptyEnv("ZKROUTER_BASE_URL", config.ZkRouterBaseURL),
			"rpcUrl":        firstNonEmptyEnv("SOLANA_RPC_URL", firstNonEmptyEnv("HELIUS_RPC_URL", config.PublicRPCURL)),
			"fundingStatus": record.FundingStatus,
			"installLedger": ledgerPath,
		}

		if recipient == "" {
			record.FundingStatus = "skipped_no_wallet"
			resp["fundingStatus"] = record.FundingStatus
			_ = appendInstallRecord(ledgerPath, record)
			writeJSONResponse(w, resp)
			return
		}
		if !wallet.IsValidPubkey(recipient) {
			record.FundingStatus = "skipped_invalid_wallet"
			resp["fundingStatus"] = record.FundingStatus
			_ = appendInstallRecord(ledgerPath, record)
			writeJSONResponse(w, resp)
			return
		}

		if prior, ok := findPriorFunding(ledgerPath, installID, recipient); ok {
			record.FundingStatus = "already_recorded"
			record.Funding = prior.Funding
			resp["fundingStatus"] = prior.FundingStatus
			if prior.Funding != nil {
				resp["solSignature"] = prior.Funding.SOLSignature
				resp["clawdSignature"] = prior.Funding.CLAWDSignature
			}
			_ = appendInstallRecord(ledgerPath, record)
			writeJSONResponse(w, resp)
			return
		}

		if !webEnvBool("CLAWDBOT_INSTALL_FUNDING_ENABLED") {
			record.FundingStatus = "queued"
			resp["fundingStatus"] = record.FundingStatus
			_ = appendInstallRecord(ledgerPath, record)
			writeJSONResponse(w, resp)
			return
		}

		if ok, reason := installFundingWithinCaps(ledgerPath, record.RemoteIP); !ok {
			record.FundingStatus = "rate_limited"
			record.FundingError = reason
			resp["fundingStatus"] = record.FundingStatus
			resp["fundingError"] = reason
			_ = appendInstallRecord(ledgerPath, record)
			writeJSONResponse(w, resp)
			return
		}

		fundCfg := birthfund.FromEnv(recipient, config.DefaultWorkspacePath())
		fundCfg.Enabled = true
		fundCfg.Send = webEnvBool("CLAWDBOT_INSTALL_FUNDING_SEND") || webEnvBool("CLAWDBOT_BIRTH_FUNDING_SEND")
		fundCfg.InstallID = installID
		fundCfg.LedgerPath = firstNonEmptyEnv("CLAWDBOT_BIRTH_FUNDING_LEDGER", filepath.Join(filepath.Dir(ledgerPath), "funding.jsonl"))
		if req.Funding.SOLLamports > 0 {
			fundCfg.SOLAmount = strconv.FormatFloat(float64(req.Funding.SOLLamports)/1_000_000_000, 'f', 9, 64)
		}
		if strings.TrimSpace(req.Funding.CLAWDTokens.String()) != "" {
			fundCfg.CLAWDAmount = strings.TrimSpace(req.Funding.CLAWDTokens.String())
		}
		if strings.TrimSpace(req.Funding.CLAWDMint) != "" {
			fundCfg.CLAWDMint = strings.TrimSpace(req.Funding.CLAWDMint)
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(envInt("CLAWDBOT_INSTALL_FUNDING_TIMEOUT_SECONDS", 180))*time.Second)
		defer cancel()
		result, err := birthfund.Fund(ctx, fundCfg, birthfund.ExecRunner{})
		record.Funding = &result
		record.FundingStatus = result.Status
		resp["fundingStatus"] = result.Status
		resp["solSignature"] = result.SOLSignature
		resp["clawdSignature"] = result.CLAWDSignature
		if err != nil {
			record.FundingError = sanitizeFundingError(err.Error())
			resp["fundingError"] = record.FundingError
		}

		_ = appendInstallRecord(ledgerPath, record)
		writeJSONResponse(w, resp)
	}
}

func installsAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		adminToken := strings.TrimSpace(os.Getenv("CLAWDBOT_INSTALL_ADMIN_TOKEN"))
		if adminToken == "" || !constantTimeEqual(bearerToken(r), adminToken) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		limit := envInt("CLAWDBOT_INSTALLS_API_LIMIT", 100)
		if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
			if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 1000 {
				limit = parsed
			}
		}
		records := readInstallRecords(installLedgerPath())
		if len(records) > limit {
			records = records[len(records)-limit:]
		}
		writeJSONResponse(w, map[string]any{
			"ok":       true,
			"count":    len(records),
			"installs": records,
		})
	}
}

func writeJSONResponse(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(value)
}

func installLedgerPath() string {
	if path := strings.TrimSpace(os.Getenv("CLAWDBOT_INSTALL_LEDGER")); path != "" {
		return path
	}
	if info, err := os.Stat("/data"); err == nil && info.IsDir() {
		return "/data/installs.jsonl"
	}
	return filepath.Join(config.DefaultWorkspacePath(), "installs.jsonl")
}

func appendInstallRecord(path string, record installRecord) error {
	installLedgerMu.Lock()
	defer installLedgerMu.Unlock()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(data, '\n'))
	return err
}

func readInstallRecords(path string) []installRecord {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	records := make([]installRecord, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var record installRecord
		if err := json.Unmarshal([]byte(line), &record); err == nil {
			records = append(records, record)
		}
	}
	return records
}

func findPriorFunding(path, installID, recipient string) (installRecord, bool) {
	records := readInstallRecords(path)
	for i := len(records) - 1; i >= 0; i-- {
		record := records[i]
		if record.AgentWalletPubkey != recipient && record.InstallID != installID {
			continue
		}
		if record.Funding == nil {
			continue
		}
		if record.Funding.Status == "sent" || record.Funding.SOLSignature != "" || record.Funding.CLAWDSignature != "" {
			return record, true
		}
	}
	return installRecord{}, false
}

func installFundingWithinCaps(path, remoteIP string) (bool, string) {
	records := readInstallRecords(path)
	since := time.Now().UTC().Add(-24 * time.Hour)
	perIP := 0
	total := 0
	for _, record := range records {
		if record.Funding == nil {
			continue
		}
		if record.Funding.Status != "sent" && record.Funding.SOLSignature == "" && record.Funding.CLAWDSignature == "" {
			continue
		}
		createdAt, err := time.Parse(time.RFC3339, record.CreatedAt)
		if err != nil || createdAt.Before(since) {
			continue
		}
		total++
		if record.RemoteIP == remoteIP {
			perIP++
		}
	}
	maxPerIP := envInt("CLAWDBOT_INSTALL_FUNDING_MAX_PER_IP_DAY", 3)
	maxPerDay := envInt("CLAWDBOT_INSTALL_FUNDING_MAX_PER_DAY", 100)
	if maxPerIP > 0 && perIP >= maxPerIP {
		return false, fmt.Sprintf("daily per-IP funding cap reached (%d)", maxPerIP)
	}
	if maxPerDay > 0 && total >= maxPerDay {
		return false, fmt.Sprintf("daily global funding cap reached (%d)", maxPerDay)
	}
	return true, ""
}

func randomInstallID() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("cb_%d", time.Now().UnixNano())
	}
	return "cb_" + hex.EncodeToString(buf[:])
}

func clientIP(r *http.Request) string {
	if webEnvBool("CLAWDBOT_TRUST_PROXY_HEADERS") {
		for _, key := range []string{"Fly-Client-IP", "CF-Connecting-IP", "X-Real-IP"} {
			value := strings.TrimSpace(r.Header.Get(key))
			if value != "" {
				return value
			}
		}
		if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
			parts := strings.Split(forwarded, ",")
			if len(parts) > 0 {
				return strings.TrimSpace(parts[0])
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func constantTimeEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func bearerToken(r *http.Request) string {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	prefix := "Bearer "
	if strings.HasPrefix(auth, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(auth, prefix))
	}
	return ""
}

func firstNonEmptyEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func webEnvBool(key string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func sanitizeFundingError(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return strings.ReplaceAll(value, os.TempDir(), "<tmp>")
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	return cmd.Start()
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" {
			if !corsAllowedOrigin(r, origin) {
				http.Error(w, "cors origin not allowed", http.StatusForbidden)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsAllowedOrigin(r *http.Request, origin string) bool {
	configured := strings.TrimSpace(os.Getenv("CLAWDBOT_CORS_ORIGINS"))
	if configured != "" {
		for _, allowed := range strings.Split(configured, ",") {
			allowed = strings.TrimSpace(allowed)
			if allowed == "*" || strings.EqualFold(allowed, origin) {
				return true
			}
		}
	}

	parsed, err := url.Parse(origin)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	return sameHost(parsed.Host, r.Host)
}

func sameHost(a, b string) bool {
	ahost, aport := splitHostPort(a)
	bhost, bport := splitHostPort(b)
	return strings.EqualFold(ahost, bhost) && aport == bport
}

func splitHostPort(value string) (string, string) {
	host, port, err := net.SplitHostPort(value)
	if err == nil {
		return strings.Trim(strings.ToLower(host), "[]"), port
	}
	return strings.Trim(strings.ToLower(value), "[]"), ""
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s %s", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// findProjectRoot walks up from the config dir to find the go.mod file.
func findProjectRoot(configPath string) string {
	dir := filepath.Dir(configPath)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return dir
		}
		dir = parent
	}
}

// listGoPackages scans the pkg/ directory for Go packages (directories with .go files).
type PackageInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	FileCount   int    `json:"file_count"`
	Description string `json:"description,omitempty"`
}

func listGoPackages(projectRoot string) []PackageInfo {
	pkgDir := filepath.Join(projectRoot, "pkg")
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return nil
	}

	var pkgs []PackageInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pkgPath := filepath.Join(pkgDir, entry.Name())
		files, err := os.ReadDir(pkgPath)
		if err != nil {
			continue
		}
		goFileCount := 0
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".go") {
				goFileCount++
			}
		}
		if goFileCount == 0 {
			continue
		}
		info := PackageInfo{
			Name:      entry.Name(),
			Path:      "pkg/" + entry.Name(),
			FileCount: goFileCount,
		}
		pkgs = append(pkgs, info)
	}
	return pkgs
}

const fallbackHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ClawdBot OS — Console</title>
<link href="https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap" rel="stylesheet">
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{background:#020208;color:#c8d8e8;font-family:'Share Tech Mono',monospace;min-height:100vh;display:flex;align-items:center;justify-content:center}
.container{text-align:center;padding:2rem}
h1{color:#14F195;font-size:2rem;margin-bottom:1rem}
.status{color:#9945FF;margin:1rem 0}
.info{color:#556680;font-size:0.9em}
a{color:#00d4ff;text-decoration:none}
a:hover{text-decoration:underline}
</style>
</head>
<body>
<div class="container">
  <h1>🦞 ClawdBot OS</h1>
  <p class="status">Web Console Running</p>
  <p>API: <a href="/api/status">/api/status</a> | <a href="/api/dna">/api/dna</a> | <a href="/api/connectors">/api/connectors</a> | <a href="/api/trading/cockpit">/api/trading/cockpit</a> | <a href="/api/laws">/api/laws</a> | <a href="/api/doctor">/api/doctor</a></p>
  <p class="info">Build the frontend with: cd web/frontend && npm run build</p>
</div>
</body>
</html>`
