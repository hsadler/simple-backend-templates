# Golang + Gin + pgx Template

A simple "items" CRUD API implementation using:
- [Golang](https://golang.org/) - Programming language
- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go

## Requirements

- Go 1.24+
- Docker and Docker Compose
- Make

## Getting Started

### Running application

Build images
```bash
docker compose build
```

Run containers locally
```bash
docker compose up
```

Verify server is running
```bash
http GET http://localhost:8000/status
```

Running the docker containers will also spin-up:
- [Swagger docs](http://localhost:8000/docs/index.html)
- [Adminer](http://127.0.0.1:8080/?pgsql=db&username=user&db=example_db&ns=public)

## Try out the "items" example API

POST an item
```bash
http POST http://127.0.0.1:8000/api/items data:='{"name": "foo", "price": 3.14}'
```

GET a single item
```bash
http GET http://127.0.0.1:8000/api/items/1
```

GET multiple items
```bash
http GET 'http://127.0.0.1:8000/api/items' item_ids==1 item_ids==2
```

### Development

Install dependencies
```bash
go mod download
```

Make sure the latest version of the "swag" documentation generator is installed
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Generate API docs from the code annotations
```bash
make gen-docs
```

### Database migrations

First, have all docker-compose containers running with `make up`.

Create a new migration
```bash
docker compose run app migrate create -ext sql -dir ./migrations -seq <migration_name>
```

Write your "up" and "down" SQL into the new migration files.

Run all migrations
```bash
make db-migrate-up
```

## Other dev commands

See the Makefile.
