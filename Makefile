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

.PHONY: tests
tests:
	go test -tags=integration -v ./tests/...

protobuf:
	protoc --go_out=. --go-grpc_out=. --proto_path=./pkg/protobuf dbcache.proto