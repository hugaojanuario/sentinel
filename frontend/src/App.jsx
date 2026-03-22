import { useState, useEffect, useCallback } from 'react'
import { listContainers } from './api'
import ContainerCard from './components/ContainerCard'
import ContainerPanel from './components/ContainerPanel'

export default function App() {
  const [containers, setContainers] = useState([])
  const [selected, setSelected] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [systemExpanded, setSystemExpanded] = useState(false)

  const fetchContainers = useCallback(async () => {
    try {
      setError(null)
      const data = await listContainers()
      setContainers(data)
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchContainers()
    const interval = setInterval(fetchContainers, 15000)
    return () => clearInterval(interval)
  }, [fetchContainers])

  function handleSelect(container) {
    setSelected((prev) => (prev?.id === container.id ? null : container))
  }

  const systemContainers = containers.filter((c) =>
    c.name.toLowerCase().includes('sentinel')
  )
  const userContainers = containers.filter(
    (c) => !c.name.toLowerCase().includes('sentinel')
  )

  return (
    <div className="app">
      <header className="header">
        <div className="logo">
          <img
            src="https://github.com/user-attachments/assets/35f3d090-1e6a-40c1-9621-894ff89004ea"
            alt="Sentinel"
            className="logo-img"
            width="38"
            height="38"
          />
          <span className="logo-text">Sentinel</span>
        </div>
        <div className="header-right">
          <span className="container-count">{userContainers.length} containers</span>
          <button className="refresh-btn" onClick={fetchContainers}>↺ Refresh</button>
        </div>
      </header>

      <div className="main">
        <div className={`container-list ${selected ? 'shrunk' : ''}`}>
          {systemContainers.length > 0 && (
            <div className="system-group" style={{ gridColumn: '1 / -1' }}>
              <button
                className="system-toggle"
                onClick={() => setSystemExpanded((v) => !v)}
              >
                <span className={`system-arrow ${systemExpanded ? 'open' : ''}`}>›</span>
                <span>system ({systemContainers.length})</span>
              </button>
              {systemExpanded && (
                <div className="system-list">
                  {systemContainers.map((c) => (
                    <ContainerCard
                      key={c.id}
                      container={c}
                      active={selected?.id === c.id}
                      onClick={() => handleSelect(c)}
                    />
                  ))}
                </div>
              )}
            </div>
          )}
          {loading && <p className="status-msg">Loading containers...</p>}
          {error && <p className="status-msg error">{error}</p>}
          {!loading && !error && userContainers.length === 0 && (
            <p className="status-msg">No containers running.</p>
          )}
          {userContainers.map((c) => (
            <ContainerCard
              key={c.id}
              container={c}
              active={selected?.id === c.id}
              onClick={() => handleSelect(c)}
            />
          ))}
        </div>

        {selected && (
          <ContainerPanel
            container={selected}
            onClose={() => setSelected(null)}
            onRefresh={fetchContainers}
          />
        )}
      </div>

      <footer className="footer">
        <a
          href="https://github.com/hugaojanuario"
          target="_blank"
          rel="noreferrer"
          className="footer-link"
        >
          <svg className="github-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z" />
          </svg>
          @hugaojanuario
        </a>
        <span className="footer-credits">Desenvolvido por Hugo Januario &mdash; &copy;2026</span>
      </footer>
    </div>
  )
}
