.PHONY: run build test vet fmt swagger docker-build docker-dev

-include .env

run:
	go run ./cmd/api

build:
	go build -o bin/ ./cmd/...

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w cmd config db router handler service repository dto model middleware

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs --parseDependency --parseInternal

docker-build:
	docker build -t security-api .

# Lance l'API en local avec hot-reload (bind mount du code) et l'expose sur le port PORT (défini dans .env).
docker-dev:
	docker run --rm -it \
		-p $(PORT):8080 \
		--env-file .env \
		-e HTTP_ADDR=0.0.0.0:8080 \
		-v "$(PWD)":/app \
		-v /app/tmp \
		security-api
