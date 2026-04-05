BINARY ?= sshnotify
VERSION ?= dev
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build test fmt install

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/sshnotify

test:
	go test ./...

fmt:
	gofmt -w .

install: build
	install -d /usr/local/bin
	install -m 0755 $(BINARY) /usr/local/bin/$(BINARY)
	install -d /etc/sshnotify
	@test -f /etc/sshnotify/config.yaml || install -m 0600 config.example.yaml /etc/sshnotify/config.yaml
