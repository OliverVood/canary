COMPOSE_BASE = -f compose.yaml
COMPOSE_POSTGRES = -f compose.postgres.yaml
COMPOSE_RABBITMQ = -f compose.rabbitmq.yaml

up-postgres:
	docker compose $(COMPOSE_BASE) $(COMPOSE_POSTGRES) up -d

up-rabbitmq:
	docker compose $(COMPOSE_BASE) $(COMPOSE_RABBITMQ) up -d

up:
	docker compose $(COMPOSE_BASE) $(COMPOSE_POSTGRES) $(COMPOSE_RABBITMQ) up -d

down:
	docker compose $(COMPOSE_BASE) $(COMPOSE_POSTGRES) $(COMPOSE_RABBITMQ) down
