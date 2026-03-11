# Sentinel

Container management service built with **Go**.

<p align="center">
  <img src="https://github.com/user-attachments/assets/35f3d090-1e6a-40c1-9621-894ff89004ea" width="170">
</p>

## Overview

Sentinel is a lightweight backend service designed to simplify the management of containers running on a host machine.

The application runs on the same server where containers are hosted and exposes a **REST API** that allows administrators or other services to monitor and control containers programmatically.

The project interacts directly with the Docker Engine using the official Go SDK, allowing operations such as listing containers, restarting them, retrieving logs and inspecting runtime statistics.

Sentinel aims to provide a simple and modular foundation for container monitoring and operational automation.

---

## Features

* List containers running on the host
* Restart containers
* Retrieve container logs
* Inspect container resource usage
* Manage containers through a REST API

---

## API Endpoints

### List containers

Returns all containers currently running on the host.

```
GET /containers
```

---

### Restart container

Restarts a specific container.

```
POST /containers/:id/restart
```

---

### Container logs

Returns recent logs from a container.

```
GET /containers/:id/logs
```

---

### Container statistics

Returns runtime statistics such as CPU and memory usage.

```
GET /containers/:id/stats
```

---

## Running the project

Clone the repository:

```
git clone https://github.com/your-username/sentinel.git
cd sentinel
```

Install dependencies:

```
go mod tidy
```

Run the service:

```
go run cmd/sentinel/main.go
```

The API will start on:

```
http://localhost:8080
```

---

## Project Structure

```
cmd/
  sentinel/
    main.go

internal/
  router/
  controllers/
  services/
  docker/
```

* **router**: defines API routes
* **controllers**: handles HTTP requests and responses
* **services**: contains business logic
* **docker**: communicates with the Docker Engine
