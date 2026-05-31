terraform {
  required_providers {
    docker = {
        source = "kreuzwerker/docker"
        version = "~> 3.0"
    }
  }
}

provider "docker" {}

resource "docker_network" "sentinel" {
  name = "sentinel_network"
}

resource "docker_volume" "postgres_db" {
  name = "sentinel_postgres_db"
}

resource "docker_container" "postgres" {
  name = "sentinel_postgres"
  image = "postgres:16-alpine"

  env = [
    "POSTGRES_USER=${var.db_user}",
    "POSTGRES_PASSWORD=${var.db_password}",
    "POSTGRES_DB=sentinel",
  ]

  volumes {
    volume_name = docker_volume.postgres_db.name
    container_path = "/var/lib/postgresql/data"
  }

  networks_advanced {
    name = docker_network.sentinel.name
  }  

  restart = "unless-stopped"
}

resource "docker_image" "sentinel_backend" {
  name = "devhugojanuario/sentinel:latest"
}

resource "docker_container" "backend" {
  name  = "sentinel_backend"
  image = docker_image.sentinel_backend.image_id

  env = [
    "PORT=9090",
    "DB_HOST=sentinel_postgres",
    "DB_PORT=5432",
    "DB_NAME=sentinel",
    "DB_USER=${var.db_user}",
    "DB_PASSWORD=${var.db_password}",
    "DB_SSLMODE=disable",
    "JWT_SECRET=${var.jwt_secret}",
  ]

  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }

  networks_advanced {
    name = docker_network.sentinel.name
  }

  depends_on = [docker_container.postgres]

  restart = "unless-stopped"
}