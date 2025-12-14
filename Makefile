get-dependencies:
	go get ./cmd/app

.PHONY: get-dependencies\

lint:
	golangci-lint run

.PHONY: lint

build: get-dependencies
	go build -o bin/musicstreaming ./cmd/app/

.PHONY: build

run: build
	./bin/musicstreaming -loglevel=info

.PHONY: run

# TODO: refactor this
integration:
	cd ./compose && docker-compose up -d redis postgres && \
	cd .. && \
	go test ./... ; \
	EXIT_CODE=$$? ; \
	cd ./compose && docker-compose down ; \
	exit $$EXIT_CODE

.PHONY: integration