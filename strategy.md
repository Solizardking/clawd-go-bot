# ClawdBot Strategy

This is the active strategy state for the Go runtime repository
`github.com/Solizardking/clawdbot-go`. It is intentionally runtime-local:
the broader ecosystem hub is `github.com/solizardking/solana-clawd`,
and public user-facing access routes through `https://zk.x402.wtf`
and `https://cheshireterminal.ai`.

Last updated: 2026-03-11T00:00:00.000Z
Best metric: 0.0000 (baseline)

## Active Parameters

```json
{
  "rsiOverbought": 70,
  "rsiOversold": 30,
  "emaFastPeriod": 20,
  "emaSlowPeriod": 50,
  "minVolume24h": 100000,
  "minLiquidity": 50000,
  "maxSlippage": 0.02,
  "stopLossPct": 0.08,
  "takeProfitPct": 0.20,
  "positionSizePct": 0.10,
  "fundingRateThreshold": 0.0005,
  "usePerps": true
}
```

## Change Log
(empty — baseline)
