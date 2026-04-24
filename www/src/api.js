const BASE = '/api'

export async function listContainers() {
  const res = await fetch(`${BASE}/containers`)
  if (!res.ok) throw new Error('Failed to fetch containers')
  return res.json()
}

export async function restartContainer(id) {
  const res = await fetch(`${BASE}/containers/${id}/restart`, { method: 'POST' })
  if (!res.ok) throw new Error('Failed to restart container')
  return res.json()
}

export async function getContainerLogs(id) {
  const res = await fetch(`${BASE}/containers/${id}/logs`)
  if (!res.ok) throw new Error('Failed to fetch logs')
  return res.text()
}

export async function getContainerStats(id) {
  const res = await fetch(`${BASE}/containers/${id}/stats`)
  if (!res.ok) throw new Error('Failed to fetch stats')
  return res.json()
}

export function streamContainerLogs(id, onData, onEnd) {
  const controller = new AbortController()

  fetch(`${BASE}/containers/${id}/logs/stream`, { signal: controller.signal })
    .then((res) => {
      const reader = res.body.getReader()
      const decoder = new TextDecoder()

      function read() {
        reader
          .read()
          .then(({ done, value }) => {
            if (done) { onEnd(); return }
            onData(stripDockerHeaders(decoder.decode(value, { stream: true })))
            read()
          })
          .catch((err) => {
            if (err.name !== 'AbortError') onEnd()
          })
      }

      read()
    })
    .catch((err) => {
      if (err.name !== 'AbortError') onEnd()
    })

  return () => controller.abort()
}

// Docker multiplexed stream adds 8-byte binary headers per frame.
// Strip non-printable characters (except tab, newline, carriage return).
function stripDockerHeaders(text) {
  return text.replace(/[^\x09\x0A\x0D\x20-\x7E\x80-\xFF]/g, '')
}
