.PHONY: build up down black flake8 isort mypy lint pytest pre-commit

build:
	docker compose build --no-cache

up:
	docker compose up -d app db adminer

down:
	docker compose down

# Code generation

gen-models:
	docker compose run --rm gen-models

# Database migrations

db-migrate:
	docker compose exec app alembic upgrade head

db-migrate-dry-run:
	docker compose exec app alembic upgrade head --sql

# Linting, formatting, testing

lint:
	black . && mypy .

pytest:
	pytest

test-runner:
	docker compose run --rm test-runner

# Cleanup

cleanup-images-volumes:
	@read -p "Are you sure you want to clean up images and volumes? (yes/no): " answer; \
	if [ "$$answer" = "yes" ]; then \
		docker compose down --rmi all --volumes --remove-orphans; \
	else \
		echo "Cleanup canceled."; \
	fi
