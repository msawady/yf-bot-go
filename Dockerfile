FROM golang:1.14.3-alpine as builder
LABEL maintainer="msawady <msawady79@gmail.com>"

WORKDIR /opt/yf-bot-go

COPY . .

RUN apk add --no-cache ca-certificates git
RUN go mod download
RUN go build

FROM alpine

ENV YF_BOT_CONFIG "/opt/yf-bot-go/config.toml"

RUN apk add --no-cache ca-certificates
COPY --from=builder /opt/yf-bot-go/yf-bot-go /opt/yf-bot-go/yf-bot-go
COPY --from=builder /opt/yf-bot-go/config.toml /opt/yf-bot-go/config.toml
ENTRYPOINT ["/opt/yf-bot-go/yf-bot-go"]

