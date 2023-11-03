build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

format:
	gofmt -l -s -w .

gen-docs:
	swag init

cleanup-images-volumes:
	@read -p "Are you sure you want to clean up images and volumes? (yes/no): " answer; \
	if [ "$$answer" = "yes" ]; then \
		docker compose down --rmi all --volumes --remove-orphans; \
	else \
		echo "Cleanup canceled."; \
	fi
