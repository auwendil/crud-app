MONGO_CONTAINER_NAME := mongodb_test
MONGO_IP := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(MONGO_CONTAINER_NAME))

POSTGRES_CONTAINER_NAME := postgresql_test
POSTGRES_IP := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(POSTGRES_CONTAINER_NAME))
POSTGRES_PW := postgrespw

run:
	go run ./cmd/api

run-mongo:
	go run ./cmd/api --db_type=mongodb

test:
	go test ./...

coverage:
	go test -cover ./...

build:
	go build ./cmd/api

docker-build:
	docker build -f ./Dockerfile -t crud-app .

docker-run-with-mongo: get-mongo-ip
	docker run -it -p 3000:3000 --rm --name crud-app crud-app /opt/crud-app --db_type=mongodb --conn_string="mongodb://$(MONGO_IP):27017"

docker-run-with-postgres: get-postgres-ip
	docker run -it -p 3000:3000 --rm --name crud-app crud-app /opt/crud-app --db_type=postgresql --conn_string="postgresql://postgres:$(POSTGRES_PW)@$(POSTGRES_IP):5432"

start-postgres:
	docker run \
		-e POSTGRES_PASSWORD='$(POSTGRES_PW)' \
		-p 32768:5432 \
		--name $(POSTGRES_CONTAINER_NAME) \
		-v ./sql/postgresql_init.sql:/opt/postgresql_init.sql \
		-v ./sql/postgresql_create_table.sql:/opt/postgresql_create_table.sql \
		-d \
		postgres:15.4

get-postgres-ip:
	$(eval POSTGRES_IP := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(POSTGRES_CONTAINER_NAME)))

init-postgres:
	docker exec $(POSTGRES_CONTAINER_NAME) psql -U postgres -a -f /opt/postgresql_init.sql
	docker exec $(POSTGRES_CONTAINER_NAME) psql -U postgres -d books -a -f /opt/postgresql_create_table.sql

stop-postgres:
	docker stop postgresql_test

start-mongodb:
	docker run \
		-p 27017:27017 \
		--name $(MONGO_CONTAINER_NAME) \
		-d \
		mongo:6.0

get-mongo-ip:
	$(eval MONGO_IP := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(MONGO_CONTAINER_NAME)))

stop-mongodb:
	docker stop $(MONGO_CONTAINER_NAME)

wait-for-postgres:
	# wait for some time so postgres will start and creating db and tables will be available
	# TODO: wait for container status / db ping instead of hardcoded value
	sleep 2

start-db: start-postgres start-mongodb wait-for-postgres init-postgres

start-with-postgres: docker-build start-db get-postgres-ip docker-run-with-postgres

start-with-mongo: docker-build start-db get-mongo-ip docker-run-with-mongo

