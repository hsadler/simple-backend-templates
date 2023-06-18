# Base image with common dependencies
FROM python:3.9-slim-bullseye AS base
WORKDIR /app
RUN pip install poetry
COPY pyproject.toml poetry.lock ./

# Development stage
FROM base AS development
RUN poetry config virtualenvs.create false && \
    poetry install --no-interaction --no-ansi
COPY . .
EXPOSE 8000
CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]

# Production stage
FROM base AS production
RUN poetry config virtualenvs.create false && \
    poetry install --no-dev --no-interaction --no-ansi
COPY . .
EXPOSE 8000
CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000"]
