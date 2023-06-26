build:
	docker compose build

up:
	docker compose up

app-shell:
	docker compose exec app bash

db-migrate:
	docker compose exec app alembic upgrade head

db-migrate-dry-run:
	docker compose exec app alembic upgrade head --sql

black:
	poetry run black .

isort:
	poetry run isort .

flake8:
	poetry run flake8 .

mypy:
	poetry run mypy .

lint: black flake8 isort mypy

test:
	poetry run pytest
