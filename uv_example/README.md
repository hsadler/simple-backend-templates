# uv example

## Setup

create a virtual environment
```bash
uv venv
```

compile dependencies for local development
```bash
uv pip compile pyproject.toml --extra dev -o requirements.txt
```

install dependencies
```bash
uv pip sync requirements.txt
```

## Run the app with uvicorn

```bash
uv run uvicorn app:app --reload
```

## Run the app in docker

build
```bash
docker build -t uv_example .
```

run
```bash
docker run -p 8000:8000 uv_example
```
