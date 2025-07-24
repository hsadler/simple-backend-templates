# Python + FastAPI + asyncpg Template

## What does this template contain?
- FastAPI server with an "items" API
- Postgres Database
- DB connection pooling and "items" CRUD via asyncpg
- DB migrations via Alembic
- Basic server logging setup
- Swagger docs for all exposed server endpoints
- API tests via pytest and FastAPI test client
- Local dev environment with docker-compose
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
docker compose build
```

Run containers locally
```sh
docker compose up -d
```

Verify server is running by hitting the status endpoint
```sh
http GET http://localhost:8000/status
```

Run DB migrations
```sh
docker compose exec app alembic upgrade head
```

## Try out the "items" example API

POST an items
```sh
http POST http://127.0.0.1:8000/api/items data:='{"name": "foo", "price": 3.14}'
```

GET a single item
```sh
http GET http://127.0.0.1:8000/api/items/1
```

GET multiple items
```sh
http GET 'http://127.0.0.1:8000/api/items' item_ids==1 item_ids==2
```

Spin-down containers when finished
```sh
docker compose down
```

### Running the docker containers will spin-up Swagger docs and Adminer

- Visit Swagger docs here:

    ```sh
    http://127.0.0.1:8000/docs
    ```

- Visit Adminer DB management tool here:

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

Open a poetry shell
```sh
poetry shell
```

Create a new migration file
```sh
alembic revision -m "my new migration"
# and also write your `upgrade` and `downgrade` queries
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

If you get isort errors, run this command alone to fix
```sh
poetry run isort .
```
