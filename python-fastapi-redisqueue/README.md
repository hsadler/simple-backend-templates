# Python FastAPI Redis Queue

This project demonstrates a FastAPI application with Redis Queue for background task processing.

## Setup

### Local Development

```bash
poetry install
```

### Docker Setup

```bash
docker-compose build
```

## Run

Spin up the application:
```bash
docker-compose up -d
```

View logs:
```bash
docker-compose logs -f
```

Stop the application:
```bash
docker-compose down
```

## Test

Create a user:
```bash
curl -X POST http://localhost:8000/users/ -H "Content-Type: application/json" -d '{"name": "John Doe", "email": "john.doe@example.com"}'
```

Get a user:
```bash
curl http://localhost:8000/users/1
```
