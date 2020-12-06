VERSION := 0.1.0-dev0

all: shrt

go.mod: cmd/shrt/main.go *.go
	go mod tidy
	touch go.mod

shrt: go.mod
	go build -ldflags "-X main.version=${VERSION}" ./cmd/shrt

clean:
	rm -f shrt

.PHONY: all clean