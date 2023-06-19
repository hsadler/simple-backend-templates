FROM golang:1.19.10-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache git

RUN go install -mod=mod github.com/githubnemo/CompileDaemon
