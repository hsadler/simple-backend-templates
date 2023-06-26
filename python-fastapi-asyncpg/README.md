# Python + FastAPI + asyncpg Template

## Getting started

Requirements:
- docker
- poetry
- httpie

install dependencies:

```sh
poetry install
```

build images:

```sh
docker compose build
```

run containers locally:

```sh
docker compose up
```

verify server is running by hitting the status endpoint:

```sh
http GET http://localhost:8000/status
```

## Try out the "items" example API

POST an items:

```sh
http POST http://127.0.0.1:8000/api/items item:='{"name": "foo", "price": 3.14}'
```

GET a single item:

```sh
http GET http://127.0.0.1:8000/api/items/1
```

GET multiple items:

```sh
http GET 'http://127.0.0.1:8000/api/items' item_ids==1 item_ids==2
```

## Database migrations

Alembic is used to manage raw SQL migrations. Migrations are not automatically
run when doing local development, but are run automatically when a production
container is started.

The process for creating new migration files, dry-run testing, and application
of a new migration is as follows.

make sure you have the containers up and running in a terminal tab:
```sh
docker compose up
```

open a poetry shell:

```sh
poetry shell
```

create a new migration file:

```sh
alembic revision -m "my new migration"
# and also write your `upgrade` and `downgrade` queries
```

dry run your migration:
```sh
docker compose exec app alembic upgrade head --sql
```

apply your new migration:
```sh
docker compose exec app alembic upgrade head
```

(optional) roll-back your migration:
```sh
docker compose exec app alembic downgrade -1
```

## Other dev commands

enter poetry shell:

```sh
poetry shell
```

before you commit code, make sure to lint:

```sh
make lint
```

if you get isort errors, run the command alone to fix:

```sh
poetry run isort .
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
