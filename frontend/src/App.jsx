import { useState, useEffect, useCallback } from 'react'
import { listContainers } from './api'
import ContainerCard from './components/ContainerCard'
import ContainerPanel from './components/ContainerPanel'

export default function App() {
  const [containers, setContainers] = useState([])
  const [selected, setSelected] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

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

  return (
    <div className="app">
      <header className="header">
        <span className="logo">Sentinel</span>
        <div className="header-right">
          <span className="container-count">{containers.length} containers</span>
          <button className="refresh-btn" onClick={fetchContainers}>↺ Refresh</button>
        </div>
      </header>

      <div className="main">
        <div className={`container-list ${selected ? 'shrunk' : ''}`}>
          {loading && <p className="status-msg">Loading containers...</p>}
          {error && <p className="status-msg error">{error}</p>}
          {!loading && !error && containers.length === 0 && (
            <p className="status-msg">No containers running.</p>
          )}
          {containers.map((c) => (
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
    </div>
  )
}
