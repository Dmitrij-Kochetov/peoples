all: remove build run

rerun: stop remove build run

build:
	@docker compose build

run:
	@docker compose up -d --remove-orphans

remove:
	@docker compose down
	@docker compose rm

stop:
	@docker compose down
