import { useState, useEffect, useRef } from 'react'

interface Connector {
  name: string
  status: string
  type: string
}

interface StatusInfo {
  status: string
  version: string
  agent: string
  config: string
  dna_path?: string
  uptime: string
  mode: string
  go_version: string
  go_os: string
  go_arch: string
  num_cpu: number
  goroutines: number
  public_links?: {
    runtime_repo: string
    hub_repo: string
    gateway: string
    terminal: string
  }
}

interface HealthInfo {
  status: string
  agent: string
}

interface PackageInfo {
  name: string
  path: string
  file_count: number
  description?: string
}

interface EnvInfo {
  AGENT_MODE: string
  HOSTNAME: string
  PWD: string
  SHELL: string
}

interface EcosystemInfo {
  runtime_repo: string
  hub_repo: string
  gateway: string
  terminal: string
}

interface CockpitInfo {
  mode: string
  watchlist: string[]
  readiness: {
    score: number
    grade: string
    status: string
    reasons: string[]
  }
  risk: {
    maxPositionSol: number
    positionSizePct: number
    stopLossPct: number
    takeProfitPct: number
    minSignalStrength: number
    minConfidence: number
  }
}

interface AgentDNA {
  agent: {
    name: string
    role: string
  }
  sequence: {
    length: number
    value: string
  }
  metrics: {
    gcContent: number
    pamSites: unknown[]
    tataBoxes: unknown[]
    utilityScore: number
    stabilityBand: string
  }
  proof: {
    dnaId: string
    sequenceSha256: string
  }
  attestation: {
    status: string
    network: string
    pdaSeed: string
  }
}

interface DNAInfo {
  path: string
  created: boolean
  dna: AgentDNA
}

interface SignalInfo {
  samples: number
  signal: {
    direction: string
    strength: number
    rsi: number
    emaFast: number
    emaSlow: number
    emaCross: string
    atr: number
    stopLoss: number
    takeProfit: number
    reasoning: string
  }
}

interface BacktestInfo {
  bars: number
  result: {
    trades: number
    wins: number
    losses: number
    winRate: number
    totalReturn: number
    avgReturn: number
    maxDrawdown: number
    profitFactor: number
    sharpe: number
    equityCurve: number[]
  }
}

export default function App() {
  const [status, setStatus] = useState<StatusInfo | null>(null)
  const [health, setHealth] = useState<HealthInfo | null>(null)
  const [connectors, setConnectors] = useState<Connector[]>([])
  const [packages, setPackages] = useState<PackageInfo[]>([])
  const [envInfo, setEnvInfo] = useState<EnvInfo | null>(null)
  const [ecosystem, setEcosystem] = useState<EcosystemInfo | null>(null)
  const [cockpit, setCockpit] = useState<CockpitInfo | null>(null)
  const [signal, setSignal] = useState<SignalInfo | null>(null)
  const [backtest, setBacktest] = useState<BacktestInfo | null>(null)
  const [dnaInfo, setDNAInfo] = useState<DNAInfo | null>(null)
  const [configText, setConfigText] = useState<string>('')
  const [showConfig, setShowConfig] = useState(false)
  const [logs, setLogs] = useState<string[]>(['🦞 ClawdBot Console ready.'])
  const logRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const fetchAll = () => {
      fetch('/api/status').then(r => r.json()).then(setStatus).catch(() => {})
      fetch('/api/health').then(r => r.json()).then(setHealth).catch(() => {})
      fetch('/api/connectors').then(r => r.json()).then(setConnectors).catch(() => {})
      fetch('/api/packages').then(r => r.json()).then(setPackages).catch(() => {})
      fetch('/api/env').then(r => r.json()).then(setEnvInfo).catch(() => {})
      fetch('/api/ecosystem').then(r => r.json()).then(setEcosystem).catch(() => {})
      fetch('/api/trading/cockpit').then(r => r.json()).then(setCockpit).catch(() => {})
      fetch('/api/trading/signal').then(r => r.json()).then(setSignal).catch(() => {})
      fetch('/api/trading/backtest').then(r => r.json()).then(setBacktest).catch(() => {})
      fetch('/api/dna').then(r => r.json()).then(setDNAInfo).catch(() => {})
    }
    fetchAll()
    const interval = setInterval(fetchAll, 10000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    if (logRef.current) logRef.current.scrollTop = logRef.current.scrollHeight
  }, [logs])

  const connected = connectors.filter(c => c.status === 'connected').length
  const total = connectors.length

  return (
    <div className="app">
      <header className="header">
        <h1>🦞 CLAWDBOT OS</h1>
        <div className="header-right">
          <span className={`health-dot ${health?.status === 'ok' ? 'ok' : 'err'}`}>●</span>
          <span className="status">{status?.status ?? 'connecting...'}</span>
        </div>
      </header>

      <div className="dashboard">
        {/* System Status Panel */}
        <div className="panel">
          <h2>System Status</h2>
          <div style={{display:'grid',gridTemplateColumns:'1fr 1fr',gap:'12px'}}>
            <div>
              <div className="value">{status?.version ?? '—'}</div>
              <div className="label">Version</div>
            </div>
            <div>
              <div className="value">{status?.mode || 'simulated'}</div>
              <div className="label">Mode</div>
            </div>
            <div>
              <div className="value">{status?.uptime ?? '—'}</div>
              <div className="label">Uptime</div>
            </div>
            <div>
              <div className="value">{status?.go_version ?? '—'}</div>
              <div className="label">Go Version</div>
            </div>
            <div>
              <div className="value">{status?.go_os ?? '—'}/{status?.go_arch ?? '—'}</div>
              <div className="label">Platform</div>
            </div>
            <div>
              <div className="value">{status?.num_cpu ?? '—'}</div>
              <div className="label">CPU Cores</div>
            </div>
            <div>
              <div className="value">{status?.goroutines ?? '—'}</div>
              <div className="label">Goroutines</div>
            </div>
            <div>
              <div className="value">{connected}/{total}</div>
              <div className="label">Connectors</div>
            </div>
          </div>
          <div style={{marginTop:'12px',paddingTop:'12px',borderTop:'1px solid var(--border)'}}>
            <div className="label">Config path</div>
            <div className="value" style={{fontSize:'0.7rem',wordBreak:'break-all'}}>{status?.config ?? '—'}</div>
          </div>
        </div>

        {/* Trading Cockpit Panel */}
        <div className="panel">
          <h2>Trading Cockpit</h2>
          {cockpit ? (
            <>
              <div style={{display:'grid',gridTemplateColumns:'1fr 1fr',gap:'12px'}}>
                <div>
                  <div className="value">{cockpit.readiness.grade} / {cockpit.readiness.score}</div>
                  <div className="label">Readiness</div>
                </div>
                <div>
                  <div className="value">{cockpit.mode || 'simulated'}</div>
                  <div className="label">Mode</div>
                </div>
                <div>
                  <div className="value">{cockpit.watchlist.length}</div>
                  <div className="label">Watchlist</div>
                </div>
                <div>
                  <div className="value">{cockpit.risk.maxPositionSol}</div>
                  <div className="label">Max SOL</div>
                </div>
              </div>
              <div style={{marginTop:'12px',paddingTop:'12px',borderTop:'1px solid var(--border)'}}>
                <div className="label">Risk limits</div>
                <div className="value" style={{fontSize:'0.75rem'}}>
                  size {(cockpit.risk.positionSizePct * 100).toFixed(1)}% · SL {(cockpit.risk.stopLossPct * 100).toFixed(1)}% · TP {(cockpit.risk.takeProfitPct * 100).toFixed(1)}%
                </div>
              </div>
            </>
          ) : (
            <div style={{color:'var(--text-dim)'}}>Loading cockpit...</div>
          )}
        </div>

        {/* Strategy Engine Panel — live signal + backtest from the Go engine */}
        <div className="panel">
          <h2>Strategy Engine</h2>
          {signal && backtest ? (
            <>
              <div style={{display:'grid',gridTemplateColumns:'1fr 1fr',gap:'12px'}}>
                <div>
                  <div className="value" style={{
                    color: signal.signal.direction === 'long' ? 'var(--solana-green, #14F195)'
                      : signal.signal.direction === 'short' ? '#ff5d6c' : 'var(--text-dim)'
                  }}>
                    {signal.signal.direction.toUpperCase()}
                  </div>
                  <div className="label">Live Signal</div>
                </div>
                <div>
                  <div className="value">{(signal.signal.strength * 100).toFixed(0)}%</div>
                  <div className="label">Strength</div>
                </div>
                <div>
                  <div className="value">{signal.signal.rsi.toFixed(1)}</div>
                  <div className="label">RSI</div>
                </div>
                <div>
                  <div className="value">{signal.signal.emaCross}</div>
                  <div className="label">EMA Cross</div>
                </div>
              </div>
              <div style={{marginTop:'12px',paddingTop:'12px',borderTop:'1px solid var(--border)'}}>
                <div className="label">Backtest ({backtest.bars} bars, {backtest.result.trades} trades)</div>
                <Sparkline data={backtest.result.equityCurve} />
                <div style={{display:'grid',gridTemplateColumns:'1fr 1fr 1fr',gap:'8px',marginTop:'8px'}}>
                  <div>
                    <div className="value" style={{fontSize:'0.9rem'}}>{(backtest.result.winRate * 100).toFixed(0)}%</div>
                    <div className="label">Win Rate</div>
                  </div>
                  <div>
                    <div className="value" style={{fontSize:'0.9rem',color: backtest.result.totalReturn >= 0 ? 'var(--solana-green, #14F195)' : '#ff5d6c'}}>
                      {(backtest.result.totalReturn * 100).toFixed(1)}%
                    </div>
                    <div className="label">Total Return</div>
                  </div>
                  <div>
                    <div className="value" style={{fontSize:'0.9rem',color:'#ff9f45'}}>-{(backtest.result.maxDrawdown * 100).toFixed(1)}%</div>
                    <div className="label">Max DD</div>
                  </div>
                  <div>
                    <div className="value" style={{fontSize:'0.9rem'}}>{backtest.result.sharpe.toFixed(2)}</div>
                    <div className="label">Sharpe</div>
                  </div>
                  <div>
                    <div className="value" style={{fontSize:'0.9rem'}}>{backtest.result.profitFactor.toFixed(2)}</div>
                    <div className="label">Profit Factor</div>
                  </div>
                  <div>
                    <div className="value" style={{fontSize:'0.9rem'}}>{backtest.result.wins}/{backtest.result.losses}</div>
                    <div className="label">W / L</div>
                  </div>
                </div>
              </div>
            </>
          ) : (
            <div style={{color:'var(--text-dim)'}}>Running strategy engine...</div>
          )}
        </div>

        {/* Agent DNA Panel */}
        <div className="panel">
          <h2>Agent DNA</h2>
          {dnaInfo ? (
            <>
              <div style={{display:'grid',gridTemplateColumns:'1fr 1fr',gap:'12px'}}>
                <div>
                  <div className="value">{dnaInfo.dna.metrics.utilityScore}/100</div>
                  <div className="label">Utility</div>
                </div>
                <div>
                  <div className="value">{dnaInfo.dna.sequence.length}</div>
                  <div className="label">Bases</div>
                </div>
                <div>
                  <div className="value">{dnaInfo.dna.metrics.gcContent.toFixed(2)}%</div>
                  <div className="label">GC Content</div>
                </div>
                <div>
                  <div className="value">{dnaInfo.dna.metrics.pamSites.length}/{dnaInfo.dna.metrics.tataBoxes.length}</div>
                  <div className="label">PAM / TATA</div>
                </div>
              </div>
              <div style={{marginTop:'12px',paddingTop:'12px',borderTop:'1px solid var(--border)'}}>
                <div className="label">DNA ID</div>
                <div className="value" style={{fontSize:'0.75rem',wordBreak:'break-all'}}>{dnaInfo.dna.proof.dnaId}</div>
                <div className="label" style={{marginTop:'8px'}}>Attestation</div>
                <div className="value" style={{fontSize:'0.75rem'}}>{dnaInfo.dna.attestation.status} · {dnaInfo.dna.attestation.network}</div>
                <div className="dna-sequence">{dnaInfo.dna.sequence.value.slice(0, 96)}...</div>
              </div>
            </>
          ) : (
            <div style={{color:'var(--text-dim)'}}>Loading agent DNA...</div>
          )}
        </div>

        {/* Connectors Panel */}
        <div className="panel">
          <h2>Connectors</h2>
          {connectors.map(c => (
            <div key={c.name} className="connector-row">
              <span className="connector-name">{c.name}</span>
              <span className="connector-type">{c.type}</span>
              <span className={`connector-status ${c.status}`}>{c.status}</span>
            </div>
          ))}
          {connectors.length === 0 && (
            <div style={{color:'var(--text-dim)'}}>Loading connectors...</div>
          )}
          <div style={{marginTop:'12px',paddingTop:'12px',borderTop:'1px solid var(--border)'}}>
            <button className="btn-config" onClick={() => {
              if (!showConfig) {
                fetch('/api/config').then(r => r.text()).then(setConfigText).catch(() => setConfigText('Failed to load config'))
              }
              setShowConfig(!showConfig)
            }}>
              {showConfig ? 'Hide' : 'View'} Config
            </button>
            {showConfig && (
              <pre className="config-json">{configText}</pre>
            )}
          </div>
        </div>

        {/* Ecosystem Panel */}
        <div className="panel">
          <h2>Ecosystem</h2>
          {ecosystem ? (
            <div className="env-list">
              {Object.entries(ecosystem).map(([key, val]) => (
                <div key={key} className="env-row">
                  <span className="env-key">{key}</span>
                  <a className="env-val" href={val} target="_blank" rel="noreferrer">
                    {val}
                  </a>
                </div>
              ))}
            </div>
          ) : (
            <div style={{color:'var(--text-dim)'}}>Loading ecosystem links...</div>
          )}
        </div>

        {/* Go Packages Panel */}
        <div className="panel">
          <h2>Go Packages ({packages.length})</h2>
          <div className="package-grid">
            {packages.map(pkg => (
              <div key={pkg.name} className="package-item">
                <span className="package-name">{pkg.name}</span>
                <span className="package-path">{pkg.path}</span>
                <span className="package-files">{pkg.file_count} file{pkg.file_count !== 1 ? 's' : ''}</span>
              </div>
            ))}
            {packages.length === 0 && (
              <div style={{color:'var(--text-dim)'}}>Loading packages...</div>
            )}
          </div>
        </div>

        {/* Environment Panel */}
        <div className="panel">
          <h2>Environment</h2>
          {envInfo ? (
            <div className="env-list">
              {Object.entries(envInfo).map(([key, val]) => (
                <div key={key} className="env-row">
                  <span className="env-key">{key}</span>
                  <span className="env-val">{val || '—'}</span>
                </div>
              ))}
            </div>
          ) : (
            <div style={{color:'var(--text-dim)'}}>Loading environment...</div>
          )}
        </div>

        {/* System Log Panel */}
        <div className="panel chat-panel">
          <h2>System Log</h2>
          <div className="chat-messages" ref={logRef}>
            {logs.map((m, i) => (
              <div key={`${i}-${m}`} className="log-line">{m}</div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

// Sparkline renders an equity curve as an inline SVG, green when the run ends up
// and red when it ends down. No external chart dependency.
function Sparkline({ data, width = 260, height = 48 }: { data: number[]; width?: number; height?: number }) {
  if (!data || data.length < 2) {
    return <div style={{color:'var(--text-dim)',fontSize:'0.75rem'}}>no equity data</div>
  }
  const min = Math.min(...data)
  const max = Math.max(...data)
  const span = max - min || 1
  const pts = data.map((v, i) => {
    const x = (i / (data.length - 1)) * width
    const y = height - ((v - min) / span) * height
    return `${x.toFixed(1)},${y.toFixed(1)}`
  }).join(' ')
  const up = data[data.length - 1] >= data[0]
  const stroke = up ? 'var(--solana-green, #14F195)' : '#ff5d6c'
  return (
    <svg width="100%" viewBox={`0 0 ${width} ${height}`} preserveAspectRatio="none"
      role="img" aria-label="Backtest equity curve"
      style={{display:'block',marginTop:'6px',overflow:'visible'}}>
      <title>Backtest equity curve</title>
      <polyline points={pts} fill="none" stroke={stroke} strokeWidth="1.5"
        strokeLinejoin="round" strokeLinecap="round" />
    </svg>
  )
}
