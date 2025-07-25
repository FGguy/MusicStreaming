build:
	go build -o bin/musicstreaming ./cmd/app/

.PHONY: build

run: build
	./bin/musicstreaming

.PHONY: run

integration:
	cd ./compose && docker-compose up -d redis postgres && \
	cd .. && \
	go test ./... ; \
	EXIT_CODE=$$? ; \
	cd ./compose && docker-compose down ; \
	exit $$EXIT_CODE

.PHONY: integration