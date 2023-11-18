build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

app-shell:
	docker compose run app sh

db-migrate-up:
	docker compose run app sh -c \
	'migrate -path=./migrations -database="$${DATABASE_URL}?sslmode=disable" up'

db-migrate-down-1:
	docker compose run app sh -c \
	'migrate -path=./migrations -database="$${DATABASE_URL}?sslmode=disable" down 1'

db-migrate-down-all:
	docker compose run app sh -c \
	'migrate -path=./migrations -database="$${DATABASE_URL}?sslmode=disable" down -all'

format:
	gofmt -l -s -w .

test:
	go test ./tests

gen-docs:
	swag init

cleanup-images-volumes:
	@read -p "Are you sure you want to clean up images and volumes? (yes/no): " answer; \
	if [ "$$answer" = "yes" ]; then \
		docker compose down --rmi all --volumes --remove-orphans; \
	else \
		echo "Cleanup canceled."; \
	fi
