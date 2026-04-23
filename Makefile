IMAGE_NAME=go-microservice

.PHONY: build start test

build:
	DOCKER_BUILDKIT=1 docker build \
		--no-cache \
		-f Dockerfile \
		-t $(IMAGE_NAME):dev \
		.

start:
	docker compose up --build

install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	swag init \
      -g ./cmd/main.go \
      -o docs \
      --parseInternal \
      --parseDependency