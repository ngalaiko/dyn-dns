run: build
	./bin/dyn-dns

build:
	go build -o ./bin/dyn-dns ./cmd/dyn-dns

test:
	go test -v ./app
