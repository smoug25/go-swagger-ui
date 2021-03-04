VERSION=$(shell git describe --abbrev=0 --tags)
BUILD=$(shell git rev-parse --short HEAD)

# Inject the build version (commit hash) into the executable.
LDFLAGS := -ldflags "-X main.Build=$(BUILD) -X main.Version=$(VERSION)"

build: generate
	go build $(LDFLAGS) -v -i -o bin/swaggerui ./swagger.go

# Generate static vfs data
generate:
	go generate ./static || die "Failed to generate static staff."
