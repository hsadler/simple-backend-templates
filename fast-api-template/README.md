# FastAPI template

## Getting started

Requirements:
- poetry
- httpie

install dependencies:

```sh
poetry install
```

build the server image:

```sh
docker compose build
```

run server container locally:

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

## Try out the "items" example API:

POST some items:

```sh
http POST http://127.0.0.1:8000/items \
items:='[{"name": "one", "price": 1}, {"name": "two", "price": 2}]'
```

GET all items

```sh
http GET http://127.0.0.1:8000/items
```
