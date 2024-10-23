FROM golang:1.23.1-alpine AS builder

RUN apk add --no-cache upx nodejs npm

RUN npm install -g uncss

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY . .

RUN mkdir temp
RUN find sites -name "*.html" -exec cp {} temp \;
RUN cp -R css temp/
RUN uncss temp/*.html > css/bulma.css
RUN rm -rf temp

COPY go.mod go.sum ./
RUN go mod download

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
