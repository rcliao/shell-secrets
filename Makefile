.PHONY: build test vet clean install

BINARY := shell-secrets
PKG := ./cmd/shell-secrets

build:
	go build -o $(BINARY) $(PKG)

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)

install:
	go install $(PKG)
