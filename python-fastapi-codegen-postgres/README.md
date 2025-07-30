# Python + FastAPI + CodeGen + Postgres Template

## What does this template contain?
- FastAPI server with an "items" API
- OpenAPI schema
- Code generation of models from the OpenAPI schema
- Postgres Database
- DB connection pooling and "items" CRUD via asyncpg
- DB migrations via Alembic
- Local dev environment with docker-compose
- Prometheus metrics exporter for the server
- Basic server logging setup
- API tests via pytest and FastAPI test client
- Automatic linting and pytest running with pre-commit

## Getting started

Requirements:
- docker
- uv
- httpie

### Local Development

Install dependencies (including dev dependencies)
```bash
uv venv
uv pip sync requirements.txt
```

Relock (if needed)
```bash
uv pip compile pyproject.toml --extra dev -o requirements.txt
```

### Docker Setup

Build images
```sh
docker compose build --no-cache
```

Run containers locally
```sh
docker compose up -d
```

Verify server is running by hitting the ping endpoint
```sh
http GET http://localhost:8000/ping
```

Run DB migrations
```sh
docker compose exec app alembic upgrade head
```

## Try out the "items" example API

POST item
```sh
http POST http://127.0.0.1:8000/items data:='{"name": "foo", "price": 3.14}'
```

GET item
```sh
http GET http://127.0.0.1:8000/items/1
```

PATCH item
```sh
http PATCH http://127.0.0.1:8000/items/1 data:='{"name": "bar", "price": 1.23}'
```

DELETE item
```sh
http DELETE http://127.0.0.1:8000/items/1
```

Spin-down containers when finished
```sh
docker compose down
```

## Model generation based on OpenAPI schema

```bash
make gen-models
```

## Pre-commit hooks (linting, formatting, testing)

```bash
make pre-commit-runner
```

### Running the docker containers will spin-up Swagger docs and Adminer

- Swagger docs:

    ```sh
    http://127.0.0.1:8000/docs
    ```

- Adminer DB management tool:

    ```sh
    http://127.0.0.1:8080/?pgsql=db&username=user&db=example_db&ns=public
    ```

## Database migrations

Alembic is used to manage raw SQL migrations. Migrations are not automatically
run when doing local development, but _are_ run automatically when a production
container is started.

The process for creating new migration files, dry-run testing, and application
of a new migration is as follows.

Make sure you have the containers up and running in a terminal tab:
```sh
docker compose up -d
```

Create a new migration file
```sh
uv run alembic revision -m "my new migration"
# and write your `upgrade` and `downgrade` queries
```

Dry run your migration
```sh
docker compose exec app alembic upgrade head --sql
```

Apply your new migration
```sh
docker compose exec app alembic upgrade head
```

(optional) Roll-back your migration
```sh
docker compose exec app alembic downgrade -1
```

## Other dev commands

See the [Makefile](./Makefile)
