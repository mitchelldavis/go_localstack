GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
BINARY_NAME=mybinary
BINARY_UNIX=$(BINARY_NAME)_unix

all: deps mockgen test
deps:
	$(GOINSTALL) github.com/golang/mock/mockgen
mockgen:
	mkdir -p pkg/mock_localstack
	mockgen -source=pkg/localstack/dockerWrapper.go > pkg/mock_localstack/dockerWrapper.go
test: 
	$(GOTEST) -v ./...
