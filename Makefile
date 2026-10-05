include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	@docker compose up -d todo-app-postgres

env-down:
	@docker compose down todo-app-postgres

env-cleanup:
	@read -p "Clean up environment volume files? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todo-app-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Environment files have been cleaned up"; \
	else \
		echo "Environment cleanup cancelled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "The required seq parameter is missing. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm todo-app-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "The required action parameter is missing. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todo-app-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todo-app-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

logs-cleanup:
	@read -p "Clean up all log files? Risk of logs loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Log files have been cleaned up"; \
	else \
		echo "Log cleanup cancelled"; \
	fi

todo-app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todo-app/main.go

todo-app-deploy:
	@docker compose up -d --build todo-app

ps:
	@docker compose ps

