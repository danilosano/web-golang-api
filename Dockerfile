FROM golang:1.22.2

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -C cmd/server -o server .

ENTRYPOINT [ "/app/cmd/server/server" ]