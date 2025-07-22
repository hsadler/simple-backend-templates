# uv example

## Setup

create a virtual environment
```bash
uv venv
```

compile dependencies for production
```bash
uv pip compile pyproject.toml -o requirements.txt
```

compile dependencies for local development
```bash
uv pip compile pyproject.toml --extra dev -o requirements.dev.txt
```

install dependencies
```bash
uv pip sync requirements.dev.txt
```

## Run the app with uvicorn

```bash
uv run uvicorn app:app --reload
```
