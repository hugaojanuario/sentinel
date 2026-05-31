variable "db_user" {
  description = "Postgres username"
  type = string
}

variable "db_password" {
  description = "Postgres password"
  type = string
}

variable "jwt_secret" {
  description = "Jwt secret key"
  type = string
  sensitive = true
}