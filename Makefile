include dbcache.env
export $(shell env | grep MY_ENV_VAR)

build-post:
	docker build -t $(POSTGRES_IMAGE_NAME) dockerfiles/postgres
build-redis:
	docker build -t $(REDIS_IMAGE_NAME) dockerfiles/redis

run-post:
	docker run -d --name $(POSTGRES_CONTAINER_NAME) \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-p $(POSTGRES_HOST_PORT):$(POSTGRES_CONTAINER_PORT) $(POSTGRES_IMAGE_NAME)
run-redis: 
	docker run -d --name $(REDIS_CONTAINER_NAME) \
	-p $(REDIS_HOST_PORT):$(REDIS_CONTAINER_PORT) $(REDIS_IMAGE_NAME) 
run-services: run-post run-redis

stop-post: 
	docker stop $(POSTGRES_CONTAINER_NAME) 
	docker rm $(POSTGRES_CONTAINER_NAME) 
stop-redis: 
	docker stop $(REDIS_CONTAINER_NAME) 
	docker rm $(REDIS_CONTAINER_NAME)
stop: stop-post stop-redis 
 
logs-post: 
	docker logs -f $(POSTGRES_CONTAINER_NAME) 
logs-redis: 
	docker logs -f $(REDIS_CONTAINER_NAME) 

run-dbcache:
	sleep 3 && env $(shell xargs < dbcache.env) go run main.go
run: run-services run-dbcache