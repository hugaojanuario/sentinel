FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sentinel ./cmd/api

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/sentinel .

EXPOSE 9090
CMD ["./sentinel"]
