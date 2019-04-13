.PHONY: build linux docker lint-all run test test-all vendor docker-push
default: build
build:
	go build -mod=vendor -o build/autocomplete -ldflags "-s -w" cmd/main.go
linux:
	GOOS=linux go build -mod=vendor -o build/autocomplete-linux-amd64 -ldflags "-s -w" cmd/main.go
docker:
	docker build -t titusjaka/autocomplete --rm .
docker-push:
	echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin && docker push titusjaka/autocomplete
lint-all:
	golangci-lint run
test:
	go test -mod=vendor ./...
test-all:
	go test -mod=vendor -tags=integration ./...
generate:
	go generate -mod=vendor -x ./...
run:
	go run -mod=vendor cmd/main.go
vendor:
	go mod tidy && go mod vendor
