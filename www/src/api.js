const BASE = '/api'

function authHeaders() {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

function unauthorized() {
  return Object.assign(new Error('Unauthorized'), { status: 401 })
}

export async function login(email, password) {
  const res = await fetch(`${BASE}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })
  if (!res.ok) throw new Error('Email ou senha inválidos')
  const data = await res.json()
  return data.token
}

export async function listContainers() {
  const res = await fetch(`${BASE}/containers`, { headers: authHeaders() })
  if (res.status === 401) throw unauthorized()
  if (!res.ok) throw new Error('Failed to fetch containers')
  return res.json()
}

export async function restartContainer(id) {
  const res = await fetch(`${BASE}/containers/${id}/restart`, {
    method: 'POST',
    headers: authHeaders(),
  })
  if (res.status === 401) throw unauthorized()
  if (!res.ok) throw new Error('Failed to restart container')
  return res.json()
}

export async function getContainerLogs(id) {
  const res = await fetch(`${BASE}/containers/${id}/logs`, { headers: authHeaders() })
  if (res.status === 401) throw unauthorized()
  if (!res.ok) throw new Error('Failed to fetch logs')
  return res.text()
}

export async function getContainerStats(id) {
  const res = await fetch(`${BASE}/containers/${id}/stats`, { headers: authHeaders() })
  if (res.status === 401) throw unauthorized()
  if (!res.ok) throw new Error('Failed to fetch stats')
  return res.json()
}

export function streamContainerLogs(id, onData, onEnd) {
  const controller = new AbortController()

  fetch(`${BASE}/containers/${id}/logs/stream`, {
    signal: controller.signal,
    headers: authHeaders(),
  })
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
