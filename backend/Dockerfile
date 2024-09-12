FROM golang:1.23 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" lehrium-backend

FROM debian:bookworm-slim

COPY --from=builder /app/lehrium-backend /app/lehrium-backend

WORKDIR /app
ENTRYPOINT ["./lehrium-backend"]
