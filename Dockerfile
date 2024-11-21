FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=386 go build -ldflags="-s -w" -trimpath -o nginfier

CMD ["cp", "/app/nginfier", "/dist/"]