.SILENT:

all: test-cover

test:
	go test -v ./...

test-cover:
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out
	rm cover.out
