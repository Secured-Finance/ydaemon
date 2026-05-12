FROM golang:1.24-bullseye

RUN mkdir -p /app
COPY . /app
WORKDIR /app
RUN go mod tidy
RUN go mod vendor
RUN go build -o yDaemon ./cmd
ENTRYPOINT /bin/sh -c "/app/yDaemon --chains \"$SUPPORTED_CHAIN_IDS\""
