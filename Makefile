.PHONY: run-containers
run-containers:
	docker run -d -p 8080:8080 hello && \
    docker run -d -p 8081:8080 hello && \
    docker run -d -p 8082:8080 hello

.PHONY: remove-containers
remove-containers:
	docker rm -f $$(docker ps -a -q)

.PHONY: load-test
load-test:
	hey http://localhost:3000

.PHONY: build-hello
build-hello:
	docker build hello-service/ -t hello