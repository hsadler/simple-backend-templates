# Golang + Gin + pgx Template

## Getting started

Requirements:
- docker
- goenv
- httpie

Ensure this is in your `.zshrc` file or similar
```sh
eval "$(goenv init -)"
export GOBIN=$(go env GOPATH)/bin
```

Install golang for this project
```sh
goenv install
```

Verify the correct go version is running
```sh
go version
# return should specify 1.21.3
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
