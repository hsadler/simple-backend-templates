# Golang + Gin + pgx Template

## Getting started

Requirements:
- golang 1.24.0 or higher
- docker
- httpie


Make sure the latest version of the "swag" documentation generator is installed
```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

Ensure you have the following in your `.zshrc` file or similar
```sh
PATH=$(go env GOPATH)/bin:$PATH
```

Install dependencies
```sh
go install
```

Build images
```sh
docker compose build
```

Run containers locally
```sh
docker compose up
```

Run DB migrations
```sh
make db-migrate-up
```

Verify server is running by hitting the status endpoint
```sh
http GET http://localhost:8000/status
```

## Try out the "items" example API

POST an item
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

## DB migrations

First, you must have all docker-compose containers running with `make up`.

To create a new DB migration...

Create migration SQL files
```sh
docker compose run app migrate create -ext sql -dir ./migrations -seq <migration_name>
```

Write your "up" and "down" SQL into the new migration files.

Run the migrations on the DB
```sh
docker compose run app sh -c \
'migrate -path=./migrations -database="$DATABASE_URL?sslmode=disable" up'
```

## Other dev commands

Generate API docs
```sh
swag init
```

Before you commit code, make sure to lint
```sh
gofmt -l -s -w .
```

### Running the docker containers will spin-up Swagger docs and Adminer

- Visit Swagger docs here:

    ```sh
    http://localhost:8000/docs/index.html
    ```

- Visit Adminer DB management tool here:

    ```sh
    http://127.0.0.1:8080/?pgsql=db&username=user&db=example_db&ns=public
    ```
