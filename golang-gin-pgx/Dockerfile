# Base image with common dependencies
FROM golang:1.21.3-bullseye as base

WORKDIR /app
COPY . .

RUN apt-get update -y
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Development stage
FROM base AS development
RUN apt-get install -y git
RUN go install -mod=mod github.com/githubnemo/CompileDaemon
EXPOSE 8000

# Production stage
FROM base AS production
RUN go build -o app .
ENV IS_PROD=true
EXPOSE 8000
ENTRYPOINT ["./app"]
