COMPOSE_BASE = -f compose.yaml
COMPOSE_POSTGRES = -f compose.postgres.yaml

up-postgres:
	docker compose $(COMPOSE_BASE) $(COMPOSE_POSTGRES) up -d

up:
	docker compose $(COMPOSE_BASE) $(COMPOSE_POSTGRES) up -d

down:
	docker compose $(COMPOSE_BASE) $(COMPOSE_POSTGRES) down
