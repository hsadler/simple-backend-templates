build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

test-add-numbers:
	curl -X POST "http://localhost:8000/add-numbers?x=5&y=3"

test-get-result:
	@read -p "Enter job ID: " job_id; \
	curl "http://localhost:8000/add-numbers/$$job_id"
