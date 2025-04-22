# Ray Serve POC

This is a simple POC for Ray Serve. It is a simple service that loads a sentiment analysis model
from the Hugging Face Hub and serves it using Ray Serve.

## Setup

install dependencies
```bash
poetry install --no-root
```

## Run

build images
```bash
docker compose build
```

start service
```bash
docker compose up
```

## Test

test the service
```bash
curl -X POST http://localhost:8000/analyze -H "Content-Type: application/json" -d '{"text": "I love this product!"}'
```

## Stop

stop the service
```bash
docker compose down
```
