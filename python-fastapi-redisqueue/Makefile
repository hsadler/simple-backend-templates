up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

build:
	docker-compose build

test-endpoint:
	curl -X POST http://localhost:8000/users/ -H "Content-Type: application/json" -d '{"name": "John Doe", "email": "john.doe@example.com"}'

# For local development without Docker
local-up:
	poetry run uvicorn app:app --reload

local-down:
	poetry run uvicorn app:app --reload
