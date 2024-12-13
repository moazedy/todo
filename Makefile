.PHONY: run,test,run-build

run:
	docker compose up 

run-build:
	docker compose up --build

test:
	go test -v ./...
