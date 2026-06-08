# Sentinel

Container management service built with **Go**.

<p align="center">
  <img src="https://github.com/user-attachments/assets/35f3d090-1e6a-40c1-9621-894ff89004ea" width="170">
</p>

## Overview

Sentinel is a lightweight service that runs on the same host as your containers and exposes a REST API to monitor and control them. It communicates directly with the Docker Engine via the official Go SDK and ships with a React frontend served through Nginx.

---

## Features

* List all running containers
* Restart containers
* Retrieve the last 50 lines of logs
* Inspect CPU and memory usage
* Stream logs in real time
* Web interface included
* Telegram alerts when a container stops, crashes, OOM-kills, or recovers

---

## Requirements

**To run with Docker (recommended):**
- [Docker](https://docs.docker.com/get-docker/) with Compose — Windows, Linux and macOS

**To run locally:**
- [Go 1.21+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/) — only needed for the frontend

---

## Running with Docker

Clone the repository:

```bash
git clone https://github.com/hugaojanuario/sentinel.git
cd sentinel
```

Copy the env file:

```bash
cp .env.example .env
```

Open `.env` and fill in the required values — especially the Telegram credentials if you want alerts:

```
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_ID=your_chat_id_here
```

Leave both empty to disable alerts.

> **Windows users:** open `.env` and set `DOCKER_SOCK=//var/run/docker.sock`

Start everything:

```bash
docker compose up --build
```

The frontend will be at `http://localhost` and the API at `http://localhost/api`.

---

## Running locally

**Backend:**

```bash
go mod tidy
go run cmd/api/main.go
```

Starts on `http://localhost:9090`. To use a different port:

```bash
PORT=8080 go run cmd/api/main.go
```

**Frontend:**

```bash
cd frontend
npm install
npm run dev
```

Starts on `http://localhost:5173` and proxies API calls to the backend automatically.

---

## API

### List containers

Returns all running containers.

```
GET /containers
```

```json
[
  {
    "id": "a1b2c3d4e5f6",
    "name": "/my-container",
    "image": "nginx:latest",
    "status": "Up 3 hours"
  }
]
```

---

### Restart container

```
POST /containers/:id/restart
```

```json
{ "message": "container restarted" }
```

---

### Container logs

Returns the last 50 lines of logs.

```
GET /containers/:id/logs
```

```
plain text response
```

---

### Container stats

Returns current CPU and memory usage from the Docker Engine.

```
GET /containers/:id/stats
```

```json
{
  "cpu_stats": { ... },
  "memory_stats": { ... },
  ...
}
```

---

### Live log stream

Streams logs in real time using chunked transfer encoding. The connection stays open until the container stops or the client disconnects.

```
GET /containers/:id/logs/stream
```

```
Content-Type: text/plain
Transfer-Encoding: chunked
```

---

## Project structure

```
cmd/
  api/
    main.go           # entry point

internal/
  router/             # route definitions
  controllers/        # request handling
  services/           # business logic
  docker/             # Docker Engine communication

frontend/
  src/                # React application
  nginx.conf          # production server config
  Dockerfile
```

---

## Telegram alerts

Sentinel monitors all containers every 30 seconds and sends a Telegram message when:

- A container **stops or crashes** — includes exit code and reason
- A container is **OOM-killed** — out of memory
- A container enters a **restart loop**
- A container **recovers** and is running again

To enable, create a bot via [@BotFather](https://t.me/BotFather) and set the credentials in `.env`:

```
TELEGRAM_BOT_TOKEN=123456789:ABC...
TELEGRAM_CHAT_ID=your_chat_or_group_id
```

To get your chat ID: open the bot, send `/start`, then visit:
```
https://api.telegram.org/bot<TOKEN>/getUpdates
```
and find `"chat": { "id": ... }` in the response.

---

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9090` | Backend listening port |
| `FRONTEND_PORT` | `80` | Frontend exposed port |
| `DOCKER_SOCK` | `/var/run/docker.sock` | Docker socket path |
| `TELEGRAM_BOT_TOKEN` | _(empty)_ | Telegram bot token — leave empty to disable alerts |
| `TELEGRAM_CHAT_ID` | _(empty)_ | Telegram chat or group ID to receive alerts |
| `CHECK_INTERVAL` | `30s` | How often containers are checked (e.g. `10s`, `1m`) |
| `ALERT_CPU_THRESHOLD` | `80` | CPU usage % threshold (reserved for future alerts) |
| `ALERT_MEM_THRESHOLD_MB` | `500` | Memory threshold in MB (reserved for future alerts) |
# 1000