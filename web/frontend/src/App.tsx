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

export default function App() {
  const [status, setStatus] = useState<StatusInfo | null>(null)
  const [health, setHealth] = useState<HealthInfo | null>(null)
  const [connectors, setConnectors] = useState<Connector[]>([])
  const [packages, setPackages] = useState<PackageInfo[]>([])
  const [envInfo, setEnvInfo] = useState<EnvInfo | null>(null)
  const [ecosystem, setEcosystem] = useState<EcosystemInfo | null>(null)
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

        {/* Ecosystem Panel */}
        <div className="panel">
          <h2>Ecosystem</h2>
          {ecosystem ? (
            <div className="env-list">
              {Object.entries(ecosystem).map(([key, val]) => (
                <div key={key} className="env-row">
                  <span className="env-key">{key}</span>
                  <a className="env-val" href={val} target="_blank" rel="noreferrer">{val}</a>
                </div>
              ))}
            </div>
          ) : (
            <div style={{color:'var(--text-dim)'}}>Loading ecosystem links...</div>
          )}
        </div>

        {/* System Log Panel */}
        <div className="panel chat-panel">
          <h2>System Log</h2>
          <div className="chat-messages" ref={logRef}>
            {logs.map((m, i) => (
              <div key={i} className="log-line">{m}</div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
