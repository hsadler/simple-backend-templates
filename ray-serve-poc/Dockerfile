FROM python:3.12-slim-bookworm AS base

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
build-essential \
&& rm -rf /var/lib/apt/lists/*

# Install poetry
RUN pip install --no-cache-dir poetry==1.7.1

# Copy requirements and install Python dependencies
COPY pyproject.toml poetry.lock ./
RUN poetry install

# Copy the service code
COPY . .

# Expose the port that Ray Serve will use
EXPOSE 8000
