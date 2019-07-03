
.PHONY: test build

build:
	docker-compose build

test:
	docker-compose run app go test -cover ./...
