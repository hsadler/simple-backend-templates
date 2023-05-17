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

POST some items:

```sh
http POST http://127.0.0.1:8000/items \
items:='[{"name": "one", "price": 1.234}, {"name": "two", "price": 2.345}]'
```

GET a single item:

```sh
http GET http://127.0.0.1:8000/item/1
```

GET multiple items:

```sh
http GET 'http://127.0.0.1:8000/items' item_ids==1 item_ids==2
```
