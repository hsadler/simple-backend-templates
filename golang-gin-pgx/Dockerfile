# Base image with common dependencies
FROM golang:1.19.10-alpine as base
WORKDIR /app
COPY . .

# Development stage
FROM base AS development
RUN apk add --no-cache git
RUN go install -mod=mod github.com/githubnemo/CompileDaemon
EXPOSE 8000

# Production stage
FROM base AS production
RUN go build -o app .
ENV IS_PROD=true
EXPOSE 8000
ENTRYPOINT ["./app"]
