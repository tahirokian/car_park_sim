.PHONY: help build build-and-deploy deploy up start down stop restart logs ps top login-redis
.PHONY: login-prometheus login-web-backend login-rabbitmq login-queue-configurator
.PHONY: login-vehicle-entry login-vehicle-exit login-go-backend login-grafana

help:	## Show this help message
	@echo ""
	@echo "** Make sure you can run docker command without sudo **"
	@echo ""
	@echo "Usage: make option"
	@echo "Options:"
	@echo "	build                      Build services"
	@echo "	build-and-deploy           Build services and start containers"
	@echo "	deploy                     Create and start containers in detached mode"
	@echo "	down                       Stop and remove containers, networks"
	@echo "	login-go-backend           Login to go backend container"
	@echo "	login-grafana              Login to grafana container"
	@echo "	login-prometheus           Login to prometheus container"
	@echo "	login-queue-configurator   Login to queue configurator container"
	@echo "	login-rabbitmq             Login to rabbitmq container"
	@echo "	login-redis                Login to rabbitmq container"
	@echo "	login-vehicle-entry        Login to vehicle entry container"
	@echo "	login-vehicle-exit         Login to vehicle exit container"
	@echo "	login-web-backend          Login to web backend container"
	@echo "	logs                       View output from containers (last 200 log lines)"
	@echo "	ps                         List containers"
	@echo "	restart                    Stop and star the services"
	@echo "	start                      Starts existing containers for services"
	@echo "	stop                       Stops running containers without removing them"
	@echo "	top                        Display the running processes"
	@echo "	up                         Create and start containers"

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
