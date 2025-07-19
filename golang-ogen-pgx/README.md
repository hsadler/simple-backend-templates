# Golang + ogen + pgx

This is a simple CRUD API for items implemented using:
- [Golang](https://golang.org/) - Programming language
- [ogen](https://github.com/ogen-go/ogen) - OpenAPI code generator for Go
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go

## Requirements

- Go 1.21+
- Docker and Docker Compose
- Make

## Getting Started

### Running application

Spin up the application with:
```bash
make up
```

This will start:
- The API server on port 8000
- PostgreSQL database on port 5433
- Adminer (database management tool) on port 8080

Spin down the application with:
```bash
make down
```

### Development Setup

1. Install dependencies:

```bash
go mod download
```

2. Generate API code from OpenAPI spec:

```bash
make openapi-generate
```

### Running the docker containers will spin-up Adminer

- Visit Adminer DB management tool here:

    ```sh
    http://127.0.0.1:8080/?pgsql=db&username=user&db=example_db&ns=public
    ```