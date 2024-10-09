FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app .

FROM debian:stretch-slim

# 静态文件
COPY ./config.yml /config.yml

COPY --from=builder /build/app /

RUN chmod 755 app
ENV BUILD_ID=dontKillMe

ENTRYPOINT ["nohup", "./app", ">", "./app.log", "2>&1", "&"]