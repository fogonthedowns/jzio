# Go related commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test ./...
GOGET=$(GOCMD) get -v -t -d ./...

# Detect the os so that we can build proper statically linked binary
OS := $(shell uname -s | awk '{print tolower($$0)}')

# Get a short hash of the git had for building images.
TAG = $$(git rev-parse --short HEAD)

# Name of actual binary to create
BINARY = app

# GOARCH tells go build which arch. to use while building a statically linked executable
GOARCH = amd64

# Setup the -ldflags option for go build here.
# While statically linking we want to inject version related information into the binary
LDFLAGS = -ldflags="$$()"

## make - show help
 HELP_FORMAT = \
	 %help; \
         while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^(\w+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/ }; \
         print "usage: make [target]\n\n"; \
     for (keys %help) { \
         print "$$_:\n"; $$sep = " " x (20 - length $$_->[0]); \
         print "  $$_->[0]$$sep$$_->[1]\n" for @{$$help{$$_}}; \
         print "\n"; }

help:           ##@miscellaneous Show this help.
	@echo
	@echo "GOPATH=$(GOPATH)"
	@echo
	@perl -e '$(HELP_FORMAT)' $(MAKEFILE_LIST)

.PHONY: run
run: bin ## This will cause "bin" target to be build first
	./$(BINARY)-$(OS)-$(GOARCH) # Execute the binary

# bin creates a platform specific statically linked binary. Platform sepcific because if you are on
# OS-X; linux binary will not work.
.PHONY: bin
bin:
	env CGO_ENABLED=0 GOOS=$(OS) GOARCH=${GOARCH} go build -a -installsuffix cgo ${LDFLAGS} -o ${BINARY}-$(OS)-${GOARCH} ./src ;

.PHONY: test
test:           ## Runs unit tests
	$(GOTEST)

.PHONY: cover
cover:          ## Generates a coverage report
	${GOCMD} test -coverprofile=coverage.out ./... && ${GOCMD} tool cover -html=coverage.out

.SILENT: clean
.PHONY: clean
clean:          ## Remove coverage report and the binary
	$(GOCLEAN) --testcache
	@rm -f ${BINARY}-$(OS)-${GOARCH}
	@rm -f coverage.out

.PHONY: goget
goget:   ## Runs go get
	$(GOGET)

