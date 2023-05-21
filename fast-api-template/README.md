# FastAPI template

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

## Other dev commands:

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

## Try out the "items" example API:

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
