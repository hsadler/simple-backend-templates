# FastAPI template

## Getting started

Requirements:
- poetry
- httpie

install dependencies:

```sh
poetry install
```

run local server:

```sh
poetry run uvicorn src.main:app --reload
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