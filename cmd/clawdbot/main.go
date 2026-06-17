// ClawdBot Go — Ultra-lightweight Solana Trading Intelligence
// Adapted from PicoClaw architecture for NVIDIA Orin Nano deployment
// Built by 8BIT Labs / Factory Division
//
// Copyright (c) 2026 8BIT Labs. All rights reserved.
// License: MIT

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/8bitlabs/clawdbot/pkg/agent"
	"github.com/8bitlabs/clawdbot/pkg/config"
	"github.com/8bitlabs/clawdbot/pkg/hardware"
	"github.com/8bitlabs/clawdbot/pkg/solana"
)

const (
	colorGreen  = "\033[1;38;2;20;241;149m"
	colorPurple = "\033[1;38;2;153;69;255m"
	colorTeal   = "\033[1;38;2;0;212;255m"
	colorAmber  = "\033[1;38;2;255;170;0m"
	colorRed    = "\033[1;38;2;255;64;96m"
	colorDim    = "\033[38;2;85;102;128m"
	colorReset  = "\033[0m"

	banner = "\r\n" +
		colorGreen + "    ███╗   ███╗ █████╗ ██╗    ██╗██████╗ " + colorPurple + "██████╗  ██████╗ ████████╗\n" +
		colorGreen + "    ████╗ ████║██╔══██╗██║    ██║██╔══██╗" + colorPurple + "██╔══██╗██╔═══██╗╚══██╔══╝\n" +
		colorGreen + "    ██╔████╔██║███████║██║ █╗ ██║██║  ██║" + colorPurple + "██████╔╝██║   ██║   ██║   \n" +
		colorGreen + "    ██║╚██╔╝██║██╔══██║██║███╗██║██║  ██║" + colorPurple + "██╔══██╗██║   ██║   ██║   \n" +
		colorGreen + "    ██║ ╚═╝ ██║██║  ██║╚███╔███╔╝██████╔╝" + colorPurple + "██████╔╝╚██████╔╝   ██║   \n" +
		colorGreen + "    ╚═╝     ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚═════╝ " + colorPurple + "╚═════╝  ╚═════╝    ╚═╝   \n" +
		colorReset + "\n" +
		colorDim + "    ┌─────────────────────────────────────────────────────────┐\n" +
		colorDim + "    │" + colorTeal + "  🦞 Sentient Solana Trading Intelligence" + colorDim + "                 │\n" +
		colorDim + "    │" + colorAmber + "  NVIDIA Orin Nano · <10MB RAM · Go Runtime" + colorDim + "             │\n" +
		colorDim + "    │" + colorGreen + "  $CLAWD :: Droids Lead The Way" + colorDim + "                          │\n" +
		colorDim + "    └─────────────────────────────────────────────────────────┘\n" +
		colorReset + "\n"

	lobster = colorRed + `              ,
             /|      __
            / |   ,-~ /
           Y :|  //  /
           | jj /( .^
           >-"~"-v"
          /       Y
         jo  o    |
        ( ~T~     j
         >._-' _./
        /   "~"  |
       Y     _,  |
      /| ;-"~ _  l
     / l/ ,-"~    \
     \//\/      .- \
      Y        /    Y
      l       I     !
      ]\      _\    /"\
     (" ~----( ~   Y.  )` + colorReset + "\n"
)

func NewClawdBotCommand() *cobra.Command {
	short := fmt.Sprintf("%s ClawdBot — Sentient Solana Trading Intelligence v%s", "🦞", config.GetVersion())

	cmd := &cobra.Command{
		Use:   "clawdbot",
		Short: short,
		Long: `ClawdBot Go — Ultra-lightweight autonomous trading agent for Solana.
Powered by the PicoClaw Go runtime, adapted for NVIDIA Orin Nano hardware.

Features:
  • OODA Loop (Observe → Orient → Decide → Act)
  • ClawVault persistent memory (known/learned/inferred)
  • ClawdBot Strategy: RSI + EMA cross + ATR signal engine
  • Solana: Jupiter swaps, Birdeye analytics, Helius RPC, Aster perps
  • Arduino Modulino® I2C: LEDs, buzzer, buttons, knob, sensors
  • Dexter deep research agent
  • Multi-channel: Telegram, Discord, CLI
  • <10MB RAM, boots in <1s on ARM64`,
		Example: "clawdbot agent -m \"What is SOL price?\"\nclawdbot ooda --interval 60\nclawdbot ooda --hw-bus 1\nclawdbot hardware scan\nclawdbot hardware demo",
	}

	cmd.AddCommand(
		NewAgentCommand(),
		NewGatewayCommand(),
		NewOnboardCommand(),
		NewStatusCommand(),
		NewOODACommand(),
		NewSolanaCommand(),
		NewHardwareCommand(),
		NewVersionCommand(),
	)

	return cmd
}

// ── Agent Command ────────────────────────────────────────────────────

func NewAgentCommand() *cobra.Command {
	var message string

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Chat with ClawdBot agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			if message != "" {
				fmt.Printf("%s[CLAWDBOT]%s Processing: %s\n", colorGreen, colorReset, message)
				// TODO: Wire to LLM provider from config
				_ = cfg
				fmt.Printf("%s[CLAWDBOT]%s Agent mode ready. Model: %s\n", colorGreen, colorReset, cfg.Agents.Defaults.ModelName)
				return nil
			}

			// Interactive REPL mode
			fmt.Println(lobster)
			fmt.Printf("%s🦞 ClawdBot Interactive Mode%s\n", colorGreen, colorReset)
			fmt.Printf("%sModel: %s | Workspace: %s%s\n", colorDim, cfg.Agents.Defaults.ModelName, cfg.Agents.Defaults.Workspace, colorReset)
			fmt.Printf("%sType your message or use memory commands (!remember, !recall, !trades, !lessons)%s\n\n", colorDim, colorReset)

			return runInteractiveAgent(cfg)
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Single message to send")
	return cmd
}

// ── Gateway Command ──────────────────────────────────────────────────

func NewGatewayCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gateway",
		Short: "Start ClawdBot gateway (Telegram, Discord, WebSocket)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			fmt.Printf("%s🦞 ClawdBot Gateway starting...%s\n", colorGreen, colorReset)
			fmt.Printf("%sHost: %s:%d%s\n", colorDim, cfg.Gateway.Host, cfg.Gateway.Port, colorReset)

			// Print enabled channels
			if cfg.Channels.Telegram.Enabled {
				fmt.Printf("  %s✓%s Telegram\n", colorGreen, colorReset)
			}
			if cfg.Channels.Discord.Enabled {
				fmt.Printf("  %s✓%s Discord\n", colorGreen, colorReset)
			}

			// Print Solana connectors
			fmt.Printf("\n%sSolana Connectors:%s\n", colorAmber, colorReset)
			fmt.Printf("  Helius:  %s\n", boolIcon(cfg.Solana.HeliusAPIKey != ""))
			fmt.Printf("  Birdeye: %s\n", boolIcon(cfg.Solana.BirdeyeAPIKey != ""))
			fmt.Printf("  Jupiter: %s\n", boolIcon(cfg.Solana.JupiterEndpoint != ""))

			// TODO: Wire real gateway with OODA loop
			select {} // Block forever
		},
	}
}

// ── Onboard Command ──────────────────────────────────────────────────

func NewOnboardCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "onboard",
		Short: "Initialize ClawdBot config & workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(lobster)
			fmt.Printf("%s🦞 Welcome to ClawdBot!%s\n\n", colorGreen, colorReset)

			configPath := config.DefaultConfigPath()
			workspacePath := config.DefaultWorkspacePath()

			fmt.Printf("Creating config at:    %s%s%s\n", colorTeal, configPath, colorReset)
			fmt.Printf("Creating workspace at: %s%s%s\n", colorTeal, workspacePath, colorReset)

			if err := config.EnsureDefaults(); err != nil {
				return fmt.Errorf("onboard failed: %w", err)
			}

			fmt.Printf("\n%s✓ ClawdBot initialized!%s\n", colorGreen, colorReset)
			fmt.Printf("%sEdit %s to configure API keys.%s\n", colorDim, configPath, colorReset)
			fmt.Printf("\nQuick start:\n")
			fmt.Printf("  %sclawdbot agent -m \"Hello\"%s\n", colorGreen, colorReset)
			fmt.Printf("  %sclawdbot ooda --interval 60%s\n", colorGreen, colorReset)
			fmt.Printf("  %sclawdbot solana wallet%s\n", colorGreen, colorReset)
			return nil
		},
	}
}

// ── Status Command ───────────────────────────────────────────────────

func NewStatusCommand() *cobra.Command {
	var hwBus int

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show ClawdBot status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			fmt.Printf("%s🦞 ClawdBot Status%s\n\n", colorGreen, colorReset)
			fmt.Printf("Version:    %s\n", config.FormatVersion())
			buildTime, goVer := config.FormatBuildInfo()
			fmt.Printf("Go:         %s\n", goVer)
			fmt.Printf("Built:      %s\n", buildTime)
			fmt.Printf("Model:      %s\n", cfg.Agents.Defaults.ModelName)
			fmt.Printf("Workspace:  %s\n", cfg.Agents.Defaults.Workspace)
			fmt.Printf("OODA Int:   %ds\n", cfg.OODA.IntervalSeconds)
			fmt.Printf("Heartbeat:  %v (every %dm)\n", cfg.Heartbeat.Enabled, cfg.Heartbeat.Interval)

			fmt.Printf("\n%sStrategy:%s\n", colorPurple, colorReset)
			fmt.Printf("  Mode:     %s\n", cfg.OODA.Mode)
			fmt.Printf("  RSI:      oversold=%d overbought=%d\n",
				cfg.Strategy.RSIOversold, cfg.Strategy.RSIOverbought)
			fmt.Printf("  EMA:      fast=%d slow=%d\n",
				cfg.Strategy.EMAFastPeriod, cfg.Strategy.EMASlowPeriod)
			fmt.Printf("  SL/TP:    %.0f%% / %.0f%%\n",
				cfg.Strategy.StopLossPct*100, cfg.Strategy.TakeProfitPct*100)
			fmt.Printf("  AutoOpt:  %v\n", cfg.OODA.AutoOptimize)

			fmt.Printf("\n%sSolana Stack:%s\n", colorAmber, colorReset)
			fmt.Printf("  Helius:      %s\n", boolIcon(cfg.Solana.HeliusAPIKey != ""))
			fmt.Printf("  Birdeye:     %s\n", boolIcon(cfg.Solana.BirdeyeAPIKey != ""))
			fmt.Printf("  Birdeye WSS: %s\n", boolIcon(cfg.Solana.BirdeyeWSSURL != ""))
			fmt.Printf("  Jupiter:     %s\n", boolIcon(cfg.Solana.JupiterEndpoint != ""))
			fmt.Printf("  Aster DEX:   %s\n", boolIcon(cfg.Solana.AsterAPIKey != ""))
			fmt.Printf("  Wallet:      %s\n", truncate(cfg.Solana.WalletPubkey, 20))

			fmt.Printf("\n%sChannels:%s\n", colorPurple, colorReset)
			fmt.Printf("  Telegram: %s\n", boolIcon(cfg.Channels.Telegram.Enabled))
			fmt.Printf("  Discord:  %s\n", boolIcon(cfg.Channels.Discord.Enabled))

			fmt.Printf("\n%sHardware (I2C bus %d):%s\n", colorTeal, hwBus, colorReset)
			hwCfg := hardware.DefaultAdapterConfig()
			hwCfg.I2CBusNum = hwBus
			hw := hardware.NewHardwareAdapter(hwCfg, hardware.AgentControls{})
			if hw.IsConnected() {
				hw.PrintStatus()
			} else {
				fmt.Printf("  %s✗ No Modulino® sensors detected%s\n", colorRed, colorReset)
				fmt.Printf("  %sRun: clawdbot hardware scan%s\n", colorDim, colorReset)
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&hwBus, "hw-bus", 1, "I2C bus number to check for Modulino® hardware")
	return cmd
}

// ── OODA Command — fully wired ─────────────────────────────────────────

func NewOODACommand() *cobra.Command {
	var (
		interval int
		hwBus    int
		noHW     bool
		simMode  bool
	)

	cmd := &cobra.Command{
		Use:   "ooda",
		Short: "Start autonomous OODA trading loop",
		Long: `Start the Observe-Orient-Decide-Act autonomous trading cycle.
The agent will continuously:
  1. OBSERVE: Pull Helius on-chain + Birdeye OHLCV + Aster perps
  2. ORIENT:  RSI/EMA/ATR strategy evaluation + ClawVault recall
  3. DECIDE:  Signal scoring (strength × confidence threshold)
  4. ACT:     Open/close positions, store vault entries, adjust params

Hardware integration (when --hw-bus is set):
  Pixels  → live status (idle/signal/trade/win/loss)
  Buzzer  → audio alerts on signals, trades, wins, losses
  Button A → trigger immediate cycle
  Button B → toggle simulated/live mode
  Button C → emergency stop (closes all positions)
  Knob    → real-time RSI threshold tuning (twist to adjust, press to reset)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			if interval > 0 {
				cfg.OODA.IntervalSeconds = interval
			}
			if simMode {
				cfg.OODA.Mode = "simulated"
			}

			fmt.Printf("%s🔄 ClawdBot OODA Loop%s\n", colorGreen, colorReset)
			fmt.Printf("%sMode: %s | Interval: %ds | Watchlist: %d tokens%s\n",
				colorDim, cfg.OODA.Mode, cfg.OODA.IntervalSeconds,
				len(cfg.OODA.Watchlist), colorReset)
			fmt.Printf("%sStrategy: RSI(%d/%d) EMA(%d/%d) SL=%.0f%% TP=%.0f%%%s\n",
				colorDim,
				cfg.Strategy.RSIOversold, cfg.Strategy.RSIOverbought,
				cfg.Strategy.EMAFastPeriod, cfg.Strategy.EMASlowPeriod,
				cfg.Strategy.StopLossPct*100, cfg.Strategy.TakeProfitPct*100,
				colorReset)

			// ── Build hooks ────────────────────────────────────────────────
			var hooks agent.AgentHooks = &consoleHooks{}
			var hwAdapter *hardware.HardwareAdapter

			if !noHW {
				hwCfg := hardware.DefaultAdapterConfig()
				hwCfg.I2CBusNum = hwBus
				hwAdapter = hardware.NewHardwareAdapter(hwCfg, hardware.AgentControls{})
				if hwAdapter.IsConnected() {
					fmt.Printf("%s🎛  Hardware: %v%s\n", colorTeal, hwAdapter.ConnectedSensors(), colorReset)
					hooks = agent.NewMultiHooks(&consoleHooks{}, hwAdapter)
				} else {
					fmt.Printf("%s🎛  Hardware: not connected (stub mode)%s\n", colorDim, colorReset)
				}
			}

			fmt.Println()

			// ── Create agent ───────────────────────────────────────────────
			ooda := agent.NewOODAAgent(cfg, hooks)

			// ── Wire hardware controls back to agent ───────────────────────
			if hwAdapter != nil && hwAdapter.IsConnected() {
				hwCfg := hardware.DefaultAdapterConfig()
				hwCfg.I2CBusNum = hwBus
				controls := hardware.AgentControls{
					TriggerCycle:  ooda.TriggerCycle,
					SetMode:       ooda.SetMode,
					EmergencyStop: ooda.Stop,
					AdjustRSI:     ooda.AdjustRSI,
				}
				hwAdapter = hardware.NewHardwareAdapter(hwCfg, controls)
				hwAdapter.Start()
				defer hwAdapter.Stop()
			}

			// ── Start agent ────────────────────────────────────────────────
			if err := ooda.Start(); err != nil {
				return fmt.Errorf("agent start: %w", err)
			}

			// ── Wait for SIGINT/SIGTERM ─────────────────────────────────────
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			sig := <-sigCh

			fmt.Printf("\n%s[OODA] Signal %s — shutting down gracefully...%s\n",
				colorAmber, sig, colorReset)
			ooda.Stop()

			stats := ooda.GetStats()
			fmt.Printf("\n%s📊 Final Stats:%s\n", colorGreen, colorReset)
			fmt.Printf("  Cycles:   %v\n", stats["cycles"])
			fmt.Printf("  Trades:   %v closed\n", stats["closed_trades"])
			fmt.Printf("  Win Rate: %.1f%%\n", stats["win_rate"])
			fmt.Printf("  Avg PnL:  %.2f%%\n", stats["avg_pnl_pct"])

			return nil
		},
	}

	cmd.Flags().IntVar(&interval, "interval", 0, "OODA cycle interval in seconds (overrides config)")
	cmd.Flags().IntVar(&hwBus, "hw-bus", 1, "I2C bus number for Modulino® hardware")
	cmd.Flags().BoolVar(&noHW, "no-hw", false, "Disable hardware integration")
	cmd.Flags().BoolVar(&simMode, "sim", false, "Force simulated mode (no live trades)")
	return cmd
}

// ── Solana Command ───────────────────────────────────────────────────

func NewSolanaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solana",
		Short: "Solana tools (wallet, Birdeye, Helius DAS/SPL)",
		Long: `Solana CLI suite for:
  • Wallet/balance checks
  • Birdeye research, trending, token search
  • Helius DAS methods (assets, owner assets, search, proofs)
  • Helius SPL/RPC methods (token balances, supply, holders, generic RPC)`,
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "wallet",
			Short: "Show wallet info and balance",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("config error: %w", err)
				}
				fmt.Printf("%s💰 Solana Wallet%s\n", colorGreen, colorReset)
				fmt.Printf("Pubkey:  %s\n", cfg.Solana.WalletPubkey)
				fmt.Printf("RPC:     %s\n", truncate(cfg.Solana.HeliusRPCURL, 40))

				if cfg.Solana.WalletPubkey != "" && cfg.Solana.HeliusAPIKey != "" {
					timeout := cfg.Solana.HeliusTimeoutSeconds
					if timeout <= 0 {
						timeout = 20
					}
					retries := cfg.Solana.HeliusRetries
					if retries <= 0 {
						retries = 3
					}

					hc := solana.NewHeliusClientWithOptions(
						cfg.Solana.HeliusAPIKey,
						cfg.Solana.HeliusRPCURL,
						cfg.Solana.HeliusWSSURL,
						cfg.Solana.HeliusNetwork,
						time.Duration(timeout*float64(time.Second)),
						retries,
						750*time.Millisecond,
					)

					balance, err := hc.GetBalance(cfg.Solana.WalletPubkey)
					if err != nil {
						fmt.Printf("%sBalance lookup failed:%s %v\n", colorRed, colorReset, err)
					} else {
						fmt.Printf("Balance: %s%.6f SOL%s (%d lamports)\n", colorTeal, balance.SOL, colorReset, balance.Lamports)
					}
				} else {
					fmt.Printf("%sSet HELIUS_API_KEY and wallet_pubkey to fetch live balance.%s\n", colorDim, colorReset)
				}

				return nil
			},
		},
		&cobra.Command{
			Use:   "research [mint]",
			Short: "Deep research a Solana token via Birdeye",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("config error: %w", err)
				}
				if cfg.Solana.BirdeyeAPIKey == "" {
					return fmt.Errorf("BIRDEYE_API_KEY not set")
				}
				client := solana.NewBirdeyeClient(cfg.Solana.BirdeyeAPIKey)
				mint := args[0]
				fmt.Printf("%s🔬 Researching token: %s%s\n\n", colorTeal, mint, colorReset)

				// Metadata
				if meta, err := client.GetTokenMetadata(mint); err == nil {
					fmt.Printf("%s── Metadata ──%s\n", colorAmber, colorReset)
					fmt.Printf("  Name:     %s (%s)\n", meta.Name, meta.Symbol)
					fmt.Printf("  Decimals: %d\n", meta.Decimals)
					if meta.Extensions.Website != "" {
						fmt.Printf("  Website:  %s\n", meta.Extensions.Website)
					}
					if meta.Extensions.Twitter != "" {
						fmt.Printf("  Twitter:  %s\n", meta.Extensions.Twitter)
					}
				}

				// Market Data
				if md, err := client.GetTokenMarketData(mint); err == nil {
					fmt.Printf("\n%s── Market Data ──%s\n", colorAmber, colorReset)
					fmt.Printf("  Price:       $%.8f\n", md.Price)
					fmt.Printf("  Market Cap:  $%.0f\n", md.MarketCap)
					fmt.Printf("  FDV:         $%.0f\n", md.FDV)
					fmt.Printf("  Liquidity:   $%.0f\n", md.Liquidity)
					fmt.Printf("  Holders:     %d\n", md.Holder)
				}

				// Trade Data
				if td, err := client.GetTokenTradeData(mint); err == nil {
					fmt.Printf("\n%s── Trade Data (24h) ──%s\n", colorAmber, colorReset)
					fmt.Printf("  Volume:      $%.0f\n", td.Volume24hUSD)
					fmt.Printf("  Trades:      %d (buy: %d / sell: %d)\n", td.Trade24h, td.Buy24h, td.Sell24h)
					fmt.Printf("  Price Chg:   %.2f%%\n", td.PriceChange24hPct)
					fmt.Printf("  Wallets:     %d unique\n", td.UniqueWallet24h)
				}

				// Security
				if sec, err := client.GetTokenSecurity(mint); err == nil {
					fmt.Printf("\n%s── Security ──%s\n", colorAmber, colorReset)
					fmt.Printf("  Mutable:       %v\n", sec.IsMutable)
					fmt.Printf("  Top10 Hold%%:   %.2f%%\n", sec.Top10Percentage)
					fmt.Printf("  Mint Auth:     %s\n", sec.HasMintAuth)
					fmt.Printf("  Freeze Auth:   %s\n", sec.HasFreezeAuth)
				}

				return nil
			},
		},
		&cobra.Command{
			Use:   "trending",
			Short: "Show trending Solana tokens (Birdeye)",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("config error: %w", err)
				}
				if cfg.Solana.BirdeyeAPIKey == "" {
					return fmt.Errorf("BIRDEYE_API_KEY not set")
				}
				client := solana.NewBirdeyeClient(cfg.Solana.BirdeyeAPIKey)

				fmt.Printf("%s🌐 Trending Solana Tokens%s\n\n", colorGreen, colorReset)
				tokens, err := client.GetTrendingV3(20)
				if err != nil {
					return fmt.Errorf("birdeye trending: %w", err)
				}
				for i, t := range tokens {
					chgColor := colorGreen
					if t.PriceChange24hPct < 0 {
						chgColor = colorRed
					}
					fmt.Printf("  %2d. %s%-8s%s $%.6f  %s%+.2f%%%s  MCap: $%.0f  Vol: $%.0f\n",
						i+1, colorTeal, t.Symbol, colorReset,
						t.Price, chgColor, t.PriceChange24hPct, colorReset,
						t.MarketCap, t.Volume24hUSD)
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "search [keyword]",
			Short: "Search for Solana tokens by name or symbol",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("config error: %w", err)
				}
				if cfg.Solana.BirdeyeAPIKey == "" {
					return fmt.Errorf("BIRDEYE_API_KEY not set")
				}
				client := solana.NewBirdeyeClient(cfg.Solana.BirdeyeAPIKey)
				keyword := args[0]

				fmt.Printf("%s🔍 Searching: %s%s\n\n", colorTeal, keyword, colorReset)
				results, err := client.SearchToken(keyword, 10)
				if err != nil {
					return fmt.Errorf("birdeye search: %w", err)
				}
				for _, r := range results {
					fmt.Printf("  %s%-8s%s %s $%.8f  Liq: $%.0f\n",
						colorTeal, r.Symbol, colorReset, r.Name, r.Price, r.Liquidity)
					fmt.Printf("    %s%s%s\n", colorDim, r.Address, colorReset)
				}
				if len(results) == 0 {
					fmt.Printf("  %sNo results found%s\n", colorDim, colorReset)
				}
				return nil
			},
		},
		NewSolanaDASCommand(),
		NewSolanaSPLCommand(),
	)

	return cmd
}

func NewSolanaDASCommand() *cobra.Command {
	defaults := mustLoadConfigDefaults()

	cmd := &cobra.Command{
		Use:   "das",
		Short: "Helius DAS methods (asset, owner-assets, search)",
	}

	getAsset := &cobra.Command{
		Use:   "get-asset [id]",
		Short: "DAS getAsset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			showFungible, _ := cmd.Flags().GetBool("show-fungible")
			result, err := client.GetAsset(args[0], displayOptionsFromFlags(showFungible, false, false))
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	getAsset.Flags().Bool("show-fungible", false, "Include fungible token fields")
	addHeliusCommonFlags(getAsset, defaults)

	getAssetBatch := &cobra.Command{
		Use:   "get-asset-batch [id...]",
		Short: "DAS getAssetBatch",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			result, err := client.GetAssetBatch(args)
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	addHeliusCommonFlags(getAssetBatch, defaults)

	assetProof := &cobra.Command{
		Use:   "asset-proof [id]",
		Short: "DAS getAssetProof",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			result, err := client.GetAssetProof(args[0])
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	addHeliusCommonFlags(assetProof, defaults)

	ownerAssets := &cobra.Command{
		Use:   "owner-assets [owner]",
		Short: "DAS getAssetsByOwner",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			owner := strings.TrimSpace(cfg.Solana.WalletPubkey)
			if len(args) > 0 {
				owner = strings.TrimSpace(args[0])
			}
			if owner == "" {
				return fmt.Errorf("owner address required (pass [owner] or set solana.wallet_pubkey in config)")
			}

			page, _ := cmd.Flags().GetInt("page")
			limit, _ := cmd.Flags().GetInt("limit")
			tokenType, _ := cmd.Flags().GetString("token-type")
			showFungible, _ := cmd.Flags().GetBool("show-fungible")
			showNativeBalance, _ := cmd.Flags().GetBool("show-native-balance")
			showInscription, _ := cmd.Flags().GetBool("show-inscription")

			result, err := client.GetAssetsByOwner(
				owner,
				page,
				limit,
				tokenType,
				displayOptionsFromFlags(showFungible, showNativeBalance, showInscription),
			)
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	ownerAssets.Flags().Int("page", 1, "Page number (DAS pages start at 1)")
	ownerAssets.Flags().Int("limit", 100, "Page size")
	ownerAssets.Flags().String("token-type", "", "Optional token type: fungible|nonFungible|regularNft|compressedNft|all")
	ownerAssets.Flags().Bool("show-fungible", false, "Include fungible token fields")
	ownerAssets.Flags().Bool("show-native-balance", false, "Include native SOL balance fields")
	ownerAssets.Flags().Bool("show-inscription", false, "Include inscription fields")
	addHeliusCommonFlags(ownerAssets, defaults)

	search := &cobra.Command{
		Use:   "search",
		Short: "DAS searchAssets using raw JSON params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			rawParams, _ := cmd.Flags().GetString("params")
			params, err := parseJSONMap(rawParams)
			if err != nil {
				return err
			}

			result, err := client.SearchAssets(params)
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	search.Flags().String("params", "{}", "JSON object for searchAssets params")
	_ = search.MarkFlagRequired("params")
	addHeliusCommonFlags(search, defaults)

	assetSignatures := &cobra.Command{
		Use:   "asset-signatures [id]",
		Short: "DAS getSignaturesForAsset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			page, _ := cmd.Flags().GetInt("page")
			limit, _ := cmd.Flags().GetInt("limit")
			result, err := client.GetSignaturesForAsset(args[0], page, limit)
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	assetSignatures.Flags().Int("page", 1, "Page number")
	assetSignatures.Flags().Int("limit", 100, "Page size")
	addHeliusCommonFlags(assetSignatures, defaults)

	cmd.AddCommand(
		getAsset,
		getAssetBatch,
		assetProof,
		ownerAssets,
		search,
		assetSignatures,
	)

	return cmd
}

func NewSolanaSPLCommand() *cobra.Command {
	defaults := mustLoadConfigDefaults()

	cmd := &cobra.Command{
		Use:   "spl",
		Short: "Helius SPL + generic RPC methods",
	}

	tokenBalance := &cobra.Command{
		Use:   "token-balance [token-account]",
		Short: "RPC getTokenAccountBalance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			result, err := client.GetTokenAccountBalance(args[0])
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	addHeliusCommonFlags(tokenBalance, defaults)

	tokenAccounts := &cobra.Command{
		Use:   "token-accounts [owner]",
		Short: "RPC getTokenAccountsByOwner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			mint, _ := cmd.Flags().GetString("mint")
			programID, _ := cmd.Flags().GetString("program-id")
			encoding, _ := cmd.Flags().GetString("encoding")

			result, err := client.GetTokenAccountsByOwner(args[0], programID, mint, encoding)
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	tokenAccounts.Flags().String("mint", "", "Optional mint filter (overrides program-id)")
	tokenAccounts.Flags().String("program-id", solana.TokenProgramID, "Token program ID when mint is not provided")
	tokenAccounts.Flags().String("encoding", "jsonParsed", "Response encoding (jsonParsed|base64)")
	addHeliusCommonFlags(tokenAccounts, defaults)

	tokenSupply := &cobra.Command{
		Use:   "token-supply [mint]",
		Short: "RPC getTokenSupply",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			result, err := client.GetTokenSupply(args[0])
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	addHeliusCommonFlags(tokenSupply, defaults)

	tokenLargest := &cobra.Command{
		Use:   "token-largest [mint]",
		Short: "RPC getTokenLargestAccounts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			result, err := client.GetTokenLargestAccounts(args[0])
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	addHeliusCommonFlags(tokenLargest, defaults)

	rpc := &cobra.Command{
		Use:   "rpc [method]",
		Short: "Generic RPC passthrough",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			client, err := newHeliusClientForCLI(cmd, cfg)
			if err != nil {
				return err
			}

			rawParams, _ := cmd.Flags().GetString("params")
			params, err := parseJSONAny(rawParams)
			if err != nil {
				return err
			}

			result, err := client.RPCAny(args[0], params)
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
	rpc.Flags().String("params", "{}", "JSON params (object or array)")
	addHeliusCommonFlags(rpc, defaults)

	cmd.AddCommand(
		tokenBalance,
		tokenAccounts,
		tokenSupply,
		tokenLargest,
		rpc,
	)

	return cmd
}

func mustLoadConfigDefaults() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		return config.DefaultConfig()
	}
	return cfg
}

func addHeliusCommonFlags(cmd *cobra.Command, cfg *config.Config) {
	network := cfg.Solana.HeliusNetwork
	if network == "" {
		network = "mainnet"
	}
	timeout := cfg.Solana.HeliusTimeoutSeconds
	if timeout <= 0 {
		timeout = 20
	}
	retries := cfg.Solana.HeliusRetries
	if retries <= 0 {
		retries = 3
	}

	cmd.Flags().String("api-key", cfg.Solana.HeliusAPIKey, "Helius API key (or set HELIUS_API_KEY)")
	cmd.Flags().String("network", network, "Helius network (mainnet|devnet)")
	cmd.Flags().String("endpoint", cfg.Solana.HeliusRPCURL, "Optional custom Helius RPC endpoint")
	cmd.Flags().Float64("timeout", timeout, "RPC timeout in seconds")
	cmd.Flags().Int("retries", retries, "RPC retry attempts")
}

func newHeliusClientForCLI(cmd *cobra.Command, cfg *config.Config) (*solana.HeliusClient, error) {
	apiKey, _ := cmd.Flags().GetString("api-key")
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		apiKey = strings.TrimSpace(cfg.Solana.HeliusAPIKey)
	}
	if apiKey == "" {
		return nil, fmt.Errorf("missing Helius API key (set HELIUS_API_KEY or pass --api-key)")
	}

	network, _ := cmd.Flags().GetString("network")
	network = strings.TrimSpace(network)
	if network == "" {
		network = strings.TrimSpace(cfg.Solana.HeliusNetwork)
	}
	if network == "" {
		network = "mainnet"
	}

	endpoint, _ := cmd.Flags().GetString("endpoint")
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		endpoint = strings.TrimSpace(cfg.Solana.HeliusRPCURL)
	}

	timeout, _ := cmd.Flags().GetFloat64("timeout")
	if timeout <= 0 {
		timeout = cfg.Solana.HeliusTimeoutSeconds
	}
	if timeout <= 0 {
		timeout = 20
	}

	retries, _ := cmd.Flags().GetInt("retries")
	if retries <= 0 {
		retries = cfg.Solana.HeliusRetries
	}
	if retries <= 0 {
		retries = 3
	}

	return solana.NewHeliusClientWithOptions(
		apiKey,
		endpoint,
		cfg.Solana.HeliusWSSURL,
		network,
		time.Duration(timeout*float64(time.Second)),
		retries,
		750*time.Millisecond,
	), nil
}

func displayOptionsFromFlags(showFungible, showNativeBalance, showInscription bool) map[string]any {
	opts := map[string]any{}
	if showFungible {
		// Keep both keys for compatibility with Helius doc variants.
		opts["showFungible"] = true
		opts["showFungibleTokens"] = true
	}
	if showNativeBalance {
		opts["showNativeBalance"] = true
	}
	if showInscription {
		opts["showInscription"] = true
	}
	if len(opts) == 0 {
		return nil
	}
	return opts
}

func parseJSONMap(raw string) (map[string]any, error) {
	raw = sanitizeJSONInput(raw)
	if raw == "" {
		return map[string]any{}, nil
	}
	var params map[string]any
	if err := json.Unmarshal([]byte(raw), &params); err != nil {
		return nil, fmt.Errorf("invalid JSON object: %w", err)
	}
	if params == nil {
		params = map[string]any{}
	}
	return params, nil
}

func parseJSONAny(raw string) (any, error) {
	raw = sanitizeJSONInput(raw)
	if raw == "" {
		return map[string]any{}, nil
	}
	var params any
	if err := json.Unmarshal([]byte(raw), &params); err != nil {
		return nil, fmt.Errorf("invalid JSON params: %w", err)
	}
	if params == nil {
		params = map[string]any{}
	}
	return params, nil
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func sanitizeJSONInput(raw string) string {
	s := strings.TrimSpace(raw)
	if len(s) >= 2 {
		if (s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '`' && s[len(s)-1] == '`') {
			s = s[1 : len(s)-1]
		}
	}
	return strings.TrimSpace(s)
}

// ── Version Command ──────────────────────────────────────────────────

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("clawdbot %s\n", config.FormatVersion())
			buildTime, goVer := config.FormatBuildInfo()
			if buildTime != "" {
				fmt.Printf("built:  %s\n", buildTime)
			}
			fmt.Printf("go:     %s\n", goVer)
		},
	}
}

// ── Console hooks (AgentHooks → terminal output) ───────────────────────

type consoleHooks struct{ agent.NoopHooks }

func (c *consoleHooks) OnAgentStart(mode string, wl []string) {
	fmt.Printf("%s[OODA]%s Agent started (mode=%s watchlist=%v)\n",
		colorGreen, colorReset, mode, wl)
}
func (c *consoleHooks) OnCycleStart(n int, sol float64) {
	if sol > 0 {
		fmt.Printf("%s[OODA]%s Cycle #%d | SOL=$%.2f\n", colorTeal, colorReset, n, sol)
	} else {
		fmt.Printf("%s[OODA]%s Cycle #%d\n", colorTeal, colorReset, n)
	}
}
func (c *consoleHooks) OnSignalDetected(sym, dir string, str, conf float64) {
	fmt.Printf("%s[OODA]%s 📡 SIGNAL %s %s (strength=%.2f conf=%.2f)\n",
		colorPurple, colorReset, dir, sym, str, conf)
}
func (c *consoleHooks) OnTradeOpen(sym, dir string, price, sol float64) {
	fmt.Printf("%s[OODA]%s 📈 OPEN %s %s at $%.6f (%.4f SOL)\n",
		colorGreen, colorReset, dir, sym, price, sol)
}
func (c *consoleHooks) OnTradeClose(sym, dir string, pnl float64, outcome, reason string) {
	col := colorGreen
	if outcome == "loss" {
		col = colorRed
	}
	fmt.Printf("%s[OODA]%s 📉 CLOSE %s %s PnL=%s%.2f%%%s (%s)\n",
		col, colorReset, dir, sym, col, pnl, colorReset, reason)
}
func (c *consoleHooks) OnLearningCycle(wr, pnl float64, count int) {
	fmt.Printf("%s[OODA]%s 🧠 Learning: wr=%.1f%% pnl=%.2f%% trades=%d\n",
		colorPurple, colorReset, wr*100, pnl, count)
}
func (c *consoleHooks) OnParamsUpdated(reason string) {
	fmt.Printf("%s[OODA]%s ⚡ Params: %s\n", colorAmber, colorReset, reason)
}
func (c *consoleHooks) OnError(ctx string, err error) {
	fmt.Printf("%s[OODA]%s ❌ %s: %v\n", colorRed, colorReset, ctx, err)
}
func (c *consoleHooks) OnHeartbeat(cycleCount, openPos int) {
	fmt.Printf("%s[OODA]%s 💓 cycle=%d open=%d\n", colorDim, colorReset, cycleCount, openPos)
}

// ── Interactive REPL ─────────────────────────────────────────────────

func runInteractiveAgent(cfg *config.Config) error {
	// Minimal REPL — real implementation will wire LLM + tools
	reader := os.Stdin
	buf := make([]byte, 4096)
	for {
		fmt.Printf("%s🦞 > %s", colorGreen, colorReset)
		n, err := reader.Read(buf)
		if err != nil {
			return nil
		}
		input := string(buf[:n-1]) // trim newline

		switch {
		case input == "exit" || input == "quit":
			fmt.Printf("%s💤 ClawdBot sleeping. Vault saved.%s\n", colorDim, colorReset)
			return nil
		case input == "!trades":
			fmt.Printf("%s📊 Trade history: (not yet implemented)%s\n", colorDim, colorReset)
		case input == "!lessons":
			fmt.Printf("%s🧠 Learned patterns: (not yet implemented)%s\n", colorDim, colorReset)
		case len(input) > 10 && input[:10] == "!remember ":
			fmt.Printf("%s💾 Stored to ClawVault: %s%s\n", colorGreen, input[10:], colorReset)
		case len(input) > 8 && input[:8] == "!recall ":
			fmt.Printf("%s🔍 Searching memory: %s%s\n", colorTeal, input[8:], colorReset)
		default:
			fmt.Printf("%s[CLAWDBOT]%s Processing with %s...\n", colorGreen, colorReset, cfg.Agents.Defaults.ModelName)
			fmt.Printf("%s(LLM integration pending — connect your API keys in config)%s\n", colorDim, colorReset)
		}
	}
}

// ── Helpers ──────────────────────────────────────────────────────────

func boolIcon(b bool) string {
	if b {
		return colorGreen + "✓" + colorReset
	}
	return colorRed + "✗" + colorReset
}

func truncate(s string, maxLen int) string {
	if s == "" {
		return colorDim + "(not set)" + colorReset
	}
	if len(s) > maxLen {
		return s[:maxLen] + "…"
	}
	return s
}

func main() {
	fmt.Print(banner)
	cmd := NewClawdBotCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
