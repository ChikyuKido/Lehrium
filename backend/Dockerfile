FROM golang:1.23.1-alpine AS builder

RUN apk add --no-cache upx

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o lehrium-backend

RUN upx --best --lzma lehrium-backend

FROM alpine:latest

COPY --from=builder /app/lehrium-backend /app/lehrium-backend

WORKDIR /app
ENTRYPOINT ["./lehrium-backend"]
