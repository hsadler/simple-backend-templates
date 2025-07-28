FROM python:3.13-slim AS base

WORKDIR /app

RUN apt-get update && apt-get install -y curl gcc libc6-dev ca-certificates && \
    curl -LsSf https://astral.sh/uv/install.sh | sh && \
    ln -s /root/.local/bin/uv /usr/local/bin/uv


FROM base AS development

COPY requirements.txt .

RUN uv venv
RUN uv pip install -r requirements.txt

COPY . .

EXPOSE 8000

CMD ["uv", "run", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]


FROM base AS production

COPY requirements.txt .

RUN uv venv
RUN uv pip install -r requirements.txt

COPY . .

ENV IS_PROD=true

EXPOSE 8000

CMD ["uv", "run", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
