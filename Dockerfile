FROM golang:alpine as builder

RUN go env -w CGO_ENABLED=0\
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY . .

RUN go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags '-extldflags "-static" -s -w' -o main ./cmd/main.go

FROM alpine:latest

RUN apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone && \
    rm -rf /var/cache/apk/*

WORKDIR /data

COPY --from=builder /build/main /usr/bin/main

RUN chmod +x /usr/bin/main

ENTRYPOINT [ "/usr/bin/main" ]