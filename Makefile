.PHONY: docker-build
docker-build:
	@docker build -t tractor . && \
	docker run -p 8080:8080 -d --name tractor-container tractor

.PHONY: docker-cleanup
docker-cleanup:
	@docker container kill tractor-container && \
	docker container prune && \
	docker image prune

.PHONY: app
app:
	@cd client && \
	yarn build && \
	yarn start

.PHONY: server
server:
	@cd src && \
	go mod vendor && \
	go build cmd/tractor/main.go && \
	./main 
