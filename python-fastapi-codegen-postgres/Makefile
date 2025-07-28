.PHONY: build up down black flake8 isort mypy lint pytest pre-commit

build:
	docker compose build --no-cache

up:
	docker compose up -d

down:
	docker compose down

# Code generation

gen-models:
	docker compose run --rm app datamodel-codegen \
		--input openapi-schema.yaml \
		--input-file-type openapi \
		--output app/models.py \
		--target-python-version 3.13 \
		--field-constraints

# Database migrations

db-migrate:
	docker compose exec app alembic upgrade head

db-migrate-dry-run:
	docker compose exec app alembic upgrade head --sql

# Linting, formatting, etc.

black:
	uv run black .

flake8:
	uv run flake8 .

isort:
	uv run isort .

mypy:
	uv run mypy .

lint: black isort flake8 mypy

pytest:
	uv run pytest

pre-commit: lint pytest

# Cleanup

cleanup-images-volumes:
	@read -p "Are you sure you want to clean up images and volumes? (yes/no): " answer; \
	if [ "$$answer" = "yes" ]; then \
		docker compose down --rmi all --volumes --remove-orphans; \
	else \
		echo "Cleanup canceled."; \
	fi
