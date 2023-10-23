FROM golang:1.19-alpine AS builder

RUN apk add build-base

ENV GOPROXY='https://goproxy.cn'

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GIT_COMMIT="unknown"
RUN make build

FROM alpine

WORKDIR /app

COPY --from=builder /src/chat2data ./chat2data

EXPOSE 8088

ENTRYPOINT ["./chat2data"]


