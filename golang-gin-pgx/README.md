# Golang + Gin + pgx Template

## Getting started

Requirements:
- docker
- goenv
- httpie

install dependencies:

```sh
go install
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
http://localhost:8000/status
```

generate API docs:
```sh
swag init
```

visit local Swagger API docs
```sh
http://localhost:8000/swagger/index.html
```

before you commit code, make sure to lint:

```sh
gofmt -l -s -w .
```

<!-- ## Try out the "items" example API:

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
``` -->
