FROM golang:1.19-buster AS builder

RUN mkdir -p /src
COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:buster-slim

RUN mkdir -p /app

COPY --from=builder /src/bin/task-server /app/
VOLUME /data/conf

WORKDIR /app

CMD ["./task-server", "-conf", "/data/conf"]