all: docker-up

build:
	CGO_ENABLED="0" go build -o main

docker: build
	docker build -t "charlesbases/swagger-ui" .
	@rm -rf main

docker-up: docker
	docker-compose up -d

