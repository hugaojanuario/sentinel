import { useState, useEffect, useRef } from 'react'
import {
  getContainerStats,
  getContainerLogs,
  restartContainer,
  streamContainerLogs,
} from '../api'

export default function ContainerPanel({ container, onClose, onRefresh }) {
  const [tab, setTab] = useState('stats')
  const [stats, setStats] = useState(null)
  const [statsLoading, setStatsLoading] = useState(false)
  const [logs, setLogs] = useState('')
  const [logsLoading, setLogsLoading] = useState(false)
  const [streamText, setStreamText] = useState('')
  const [streaming, setStreaming] = useState(false)
  const [restarting, setRestarting] = useState(false)
  const stopStreamRef = useRef(null)
  const streamEndRef = useRef(null)

  function stopStream() {
    if (stopStreamRef.current) {
      stopStreamRef.current()
      stopStreamRef.current = null
    }
    setStreaming(false)
  }

  function startStream() {
    setStreaming(true)
    stopStreamRef.current = streamContainerLogs(
      container.id,
      (chunk) => setStreamText((prev) => prev + chunk),
      () => setStreaming(false)
    )
  }

  useEffect(() => {
    setStats(null)
    setLogs('')
    setStreamText('')
    setTab('stats')
    stopStream()
  }, [container.id])

  useEffect(() => {
    stopStream()

    if (tab === 'stats') {
      setStatsLoading(true)
      getContainerStats(container.id)
        .then((data) => setStats(parseStats(data)))
        .catch(() => setStats(null))
        .finally(() => setStatsLoading(false))
    }

    if (tab === 'logs') {
      setLogsLoading(true)
      getContainerLogs(container.id)
        .then((data) => setLogs(data))
        .catch(() => setLogs('Failed to load logs.'))
        .finally(() => setLogsLoading(false))
    }

    if (tab === 'stream') {
      setStreamText('')
      startStream()
    }

    return () => stopStream()
  }, [tab, container.id])

  useEffect(() => {
    streamEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [streamText])

  async function handleRestart() {
    setRestarting(true)
    try {
      await restartContainer(container.id)
      onRefresh()
    } finally {
      setRestarting(false)
    }
  }

  const name = container.name.replace(/^\//, '')

  return (
    <div className="panel">
      <div className="panel-header">
        <div>
          <div className="panel-name">{name}</div>
          <div className="panel-image">{container.image}</div>
        </div>
        <div className="panel-actions">
          <button className="btn-restart" onClick={handleRestart} disabled={restarting}>
            {restarting ? 'Restarting...' : '↺ Restart'}
          </button>
          <button className="btn-close" onClick={onClose}>✕</button>
        </div>
      </div>

      <div className="tabs">
        {['stats', 'logs', 'stream'].map((t) => (
          <button
            key={t}
            className={`tab ${tab === t ? 'active' : ''}`}
            onClick={() => setTab(t)}
          >
            {t.charAt(0).toUpperCase() + t.slice(1)}
          </button>
        ))}
      </div>

      <div className="panel-content">
        {tab === 'stats' && (
          <div className="stats-view">
            {statsLoading && <p className="status-msg">Loading stats...</p>}
            {!statsLoading && stats && (
              <>
                <div className="stat-row">
                  <span className="stat-label">CPU Usage</span>
                  <div className="stat-bar-wrap">
                    <div className="stat-bar cpu" style={{ width: `${Math.min(stats.cpu, 100)}%` }} />
                    <span className="stat-value">{stats.cpu.toFixed(2)}%</span>
                  </div>
                </div>
                <div className="stat-row">
                  <span className="stat-label">Memory</span>
                  <div className="stat-bar-wrap">
                    <div className="stat-bar mem" style={{ width: `${Math.min(stats.memPercent, 100)}%` }} />
                    <span className="stat-value">
                      {stats.memUsage} / {stats.memLimit} ({stats.memPercent.toFixed(1)}%)
                    </span>
                  </div>
                </div>
              </>
            )}
            {!statsLoading && !stats && (
              <p className="status-msg error">Stats not available.</p>
            )}
            <button
              className="btn-secondary"
              onClick={() => setTab('stats')}
              style={{ marginTop: 8 }}
            >
              ↺ Refresh
            </button>
          </div>
        )}

        {tab === 'logs' && (
          <div className="logs-view">
            {logsLoading && <p className="status-msg">Loading logs...</p>}
            {!logsLoading && (
              <>
                <pre className="log-output">{logs || 'No logs available.'}</pre>
                <button className="btn-secondary" onClick={() => setTab('logs')}>
                  ↺ Refresh
                </button>
              </>
            )}
          </div>
        )}

        {tab === 'stream' && (
          <div className="stream-view">
            <div className="stream-controls">
              <span className={`stream-status ${streaming ? 'active' : ''}`}>
                {streaming ? '● Live' : '○ Stopped'}
              </span>
              {streaming ? (
                <button className="btn-secondary" onClick={stopStream}>Stop</button>
              ) : (
                <button
                  className="btn-secondary"
                  onClick={() => { setStreamText(''); startStream() }}
                >
                  Start
                </button>
              )}
              <button className="btn-secondary" onClick={() => setStreamText('')}>
                Clear
              </button>
            </div>
            <pre className="log-output stream">
              {streamText || 'Waiting for logs...'}
              <span ref={streamEndRef} />
            </pre>
          </div>
        )}
      </div>
    </div>
  )
}

function parseStats(raw) {
  try {
    const cpuDelta =
      raw.cpu_stats.cpu_usage.total_usage - raw.precpu_stats.cpu_usage.total_usage
    const systemDelta =
      raw.cpu_stats.system_cpu_usage - raw.precpu_stats.system_cpu_usage
    const numCPUs =
      raw.cpu_stats.online_cpus || raw.cpu_stats.cpu_usage.percpu_usage?.length || 1
    const cpu = (cpuDelta / systemDelta) * numCPUs * 100

    const memUsageBytes = raw.memory_stats.usage
    const memLimitBytes = raw.memory_stats.limit
    const memPercent = (memUsageBytes / memLimitBytes) * 100

    return {
      cpu,
      memUsage: formatBytes(memUsageBytes),
      memLimit: formatBytes(memLimitBytes),
      memPercent,
    }
  } catch {
    return null
  }
}

function formatBytes(bytes) {
  if (bytes >= 1024 ** 3) return (bytes / 1024 ** 3).toFixed(1) + ' GB'
  if (bytes >= 1024 ** 2) return (bytes / 1024 ** 2).toFixed(1) + ' MB'
  return (bytes / 1024).toFixed(1) + ' KB'
}
