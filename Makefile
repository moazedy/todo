.PHONY: run,test,run-build,benchmark

run:
	docker compose up 

run-build:
	docker compose up --build

test:
	go test -v ./...


benchmark:
	@if [ ! "$(shell docker ps -q -f name=minio-benchmark)" ]; then \
		echo "Starting minio container..."; \
		docker run -d --rm --name minio-benchmark -p 9000:9000 \
  		-e MINIO_ROOT_USER=mys3accesskey \
  		-e MINIO_ROOT_PASSWORD=mys3secretkey \
  		minio/minio server /data;\
	fi
	@if [ ! "$(shell docker ps -q -f name=postgres-benchmark)" ]; then \
		echo "Starting Postgres container..."; \
		docker run -d --rm --name postgres-benchmark \
			-e POSTGRES_PASSWORD=password -p 5432:5432 postgres; \
	fi	
	@echo "Creating benchmark database..."
	@docker exec postgres-benchmark psql -U postgres -tAc "SELECT 1 FROM pg_database WHERE datname='todo_benchmark';" | grep -q 1 || \
		docker exec postgres-benchmark psql -U postgres -c "CREATE DATABASE todo_benchmark;"
	@echo "Running benchmark tests..."
	go test -v -bench=. ./...
