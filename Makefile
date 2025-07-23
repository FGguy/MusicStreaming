integration:
	docker-compose up -d redis postgres&& \
	cd ./src && \
	go test ./... ; \
	EXIT_CODE=$$? ; \
	docker-compose down ; \
	exit $$EXIT_CODE

.PHONY: integration