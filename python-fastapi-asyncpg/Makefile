build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

app-shell:
	docker compose exec app bash

db-migrate:
	docker compose exec app alembic upgrade head

db-migrate-dry-run:
	docker compose exec app alembic upgrade head --sql

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

cleanup-images-volumes:
	@read -p "Are you sure you want to clean up images and volumes? (yes/no): " answer; \
	if [ "$$answer" = "yes" ]; then \
		docker compose down --rmi all --volumes --remove-orphans; \
	else \
		echo "Cleanup canceled."; \
	fi
