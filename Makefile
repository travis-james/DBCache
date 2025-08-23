include dbcache.env
export $(shell env | grep MY_ENV_VAR)

build-dbs:
	docker compose --env-file dbcache.env build
run-dbs:
	docker compose --env-file dbcache.env up -d
stop-dbs:
	docker compose --env-file dbcache.env down

run-server:
	env $(shell xargs < dbcache.env) go run ./cmd/main.go -mode=server
run-client:
	go run ./cmd/main.go -mode=client
run-dev:
	env $(shell xargs < dbcache.env) go run ./cmd/main.go -mode=dev

.PHONY: tests
tests:
	@echo "Starting DB containers..."
	@$(MAKE) run-dbs
	@sleep 3 # Wait for dbs to finish.
	@trap '$(MAKE) stop-dbs' EXIT; \
	go test -tags=integration -v ./tests/...

protobuf:
	protoc \
	--go_out=./pkg/protobuf --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/protobuf --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=./pkg/protobuf --grpc-gateway_opt=paths=source_relative \
	--proto_path=./pkg/protobuf \
	--proto_path=./proto_deps \
	dbcache.proto