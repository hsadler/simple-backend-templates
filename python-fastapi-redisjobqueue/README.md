# Python FastAPI Redis Job Queue (POC)

This project demonstrates a FastAPI application with Redis for background job processing.

## Setup

### Local Development

Install dependencies (including dev dependencies)
```bash
uv venv
uv pip sync requirements.txt
```

### Docker Setup

Build the application:
```bash
make build
```

Spin up the application:
```bash
make up
```

Stop the application:
```bash
make down
```

## Usage

### API Endpoints

- `POST /add-numbers?x=<float>&y=<float>`: Create a new addition job
- `GET /add-numbers/<job_id>`: Check the status and result of a job

### Example Usage

Create a job and parse the response to get the job ID
```bash
JOB_ID=$(curl -X POST "http://localhost:8000/add-numbers?x=5&y=3" | jq -r '.job_id')
```

Poll the job status until it is complete
```bash
curl "http://localhost:8000/add-numbers/$JOB_ID"
```

Or...

You can also use the provided Python client:
```bash
uv run python client.py add --x=5 --y=3
```
