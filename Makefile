BINARY := baltig
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build install snapshot clean completions

build:
	go build $(LDFLAGS) -o $(BINARY) .

install: build
	mv $(BINARY) /usr/local/bin/$(BINARY)

completions:
	mkdir -p completions
	go run . completion bash > completions/baltig.bash
	go run . completion zsh > completions/_baltig
	go run . completion fish > completions/baltig.fish

snapshot:
	goreleaser release --snapshot --clean

clean:
	rm -f $(BINARY)
	rm -rf dist/
