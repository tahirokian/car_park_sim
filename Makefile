THIS_FILE := $(lastword $(MAKEFILE_LIST))

.PHONY: help build build-and-deploy deploy up start down stop restart logs ps top login-redis login-prometheus login-web-backend login-rabbitmq login-queue-configurator login-vehicle-entry login-vehicle-exit login-go-backend login-grafana

help:	## Show this help message
	@echo "Make sure you can run docker command without sudo"
	@echo ""
	@echo "Usage: make {build|build-and-deploy|deploy|up|start|stop|down|restart|logs|ps|top|login-redis|login-rabbitmq|login-web-backend|login-queue-configurator|login-vehicle-entry|login-vehicle-exit|login-go-backend|login-prometheus|login-grafana}"

build:
	docker-compose -f docker-compose.yml build

build-and-deploy:
	docker-compose -f docker-compose.yml up -d --build

deploy:
	docker-compose -f docker-compose.yml up -d

up:
	docker-compose -f docker-compose.yml up

down:
	docker-compose -f docker-compose.yml down

start:
	docker-compose -f docker-compose.yml start

stop:
	docker-compose -f docker-compose.yml stop

restart:
	docker-compose -f docker-compose.yml stop
	docker-compose -f docker-compose.yml up -d

logs:
	docker-compose -f docker-compose.yml logs --tail 200

ps:
	docker-compose -f docker-compose.yml ps

top:
	docker-compose -f docker-compose.yml top

login-redis:
	docker-compose -f docker-compose.yml exec redis sh

login-web-backend:
	docker-compose -f docker-compose.yml exec web_backend sh

login-rabbitmq:
	docker-compose -f docker-compose.yml exec rabbitmq sh

login-queue-configurator:
	docker-compose -f docker-compose.yml exec queue_configurator sh

login-vehicle-entry:
	docker-compose -f docker-compose.yml exec vehicle_entry sh

login-vehicle-exit:
	docker-compose -f docker-compose.yml exec vehicle_exit sh

login-go-backend:
	docker-compose -f docker-compose.yml exec go_backend sh

login-prometheus:
	docker-compose -f docker-compose.yml exec prometheus sh

login-grafana:
	docker-compose -f docker-compose.yml exec grafana sh
