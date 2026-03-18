export default function ContainerCard({ container, active, onClick }) {
  const isRunning = container.status.toLowerCase().includes('up')
  const name = container.name.replace(/^\//, '')

  return (
    <div className={`card ${active ? 'active' : ''}`} onClick={onClick}>
      <div className="card-header">
        <span className={`status-dot ${isRunning ? 'running' : 'stopped'}`} />
        <span className="card-name">{name}</span>
      </div>
      <div className="card-image">{container.image}</div>
      <div className={`card-status ${isRunning ? 'running' : 'stopped'}`}>
        {container.status}
      </div>
    </div>
  )
}
