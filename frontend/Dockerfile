FROM golang:1.23.1-alpine AS builder


RUN apk add --no-cache upx

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" lehrium-frontend

RUN upx --best --lzma lehrium-frontend

FROM alpine:latest

COPY --from=builder /app/lehrium-frontend /app/lehrium-frontend
COPY --from=builder /app/css /app/css
COPY --from=builder /app/imgs /app/imgs
COPY --from=builder /app/js /app/js
COPY --from=builder /app/sites /app/sites
WORKDIR /app
ENTRYPOINT ["./lehrium-frontend"]
