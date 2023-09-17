all: remove build run

rerun: stop remove build run

run-d: remove build run-detached

build:
	@docker compose build

run:
	@docker compose up --remove-orphans

run-detached:
	@docker compose up -d --remove-orphans

remove:
	@docker compose down
	@docker compose rm

stop:
	@docker compose down
