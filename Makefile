REDIS_CONTAINER_NAME = dbcache_redis
REDIS_IMAGE = redis:8
REDIS_HOST_PORT = 6379

POSTGRES_CONTAINER_NAME=dbcache_postgres
POSTGRES_IMAGE=postgres:17
POSTGRES_HOST_PORT=5432
POSTGRES_USER=local_dev
POSTGRES_PASSWORD=local_pw
POSTGRES_DB=local_db

run: run-post run-redis
run-post:
	docker run -d --name $(POSTGRES_CONTAINER_NAME) \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-p $(POSTGRES_HOST_PORT):5432 $(POSTGRES_IMAGE)
run-redis: 
	docker run -d --name $(REDIS_CONTAINER_NAME) \
	-p $(REDIS_HOST_PORT):6379 $(REDIS_IMAGE) 

stop: stop-post stop-redis
stop-post: 
	docker stop $(POSTGRES_CONTAINER_NAME) 
	docker rm $(POSTGRES_CONTAINER_NAME) 
stop-redis: 
	docker stop $(REDIS_CONTAINER_NAME) 
	docker rm $(REDIS_CONTAINER_NAME) 
 
logs-post: 
	docker logs -f $(POSTGRES_CONTAINER_NAME) 
logs-redis: 
	docker logs -f $(REDIS_CONTAINER_NAME) 