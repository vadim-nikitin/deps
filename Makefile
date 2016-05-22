# ==============================================================================
# deps - Dependency management utility for Linux executables
# ==============================================================================

# Go compiler
GO=go

# GOPATH uses directory of the project
export GOPATH := $(GOPATH):$(shell pwd)

all: install

install:
	@echo Installing with GOPATH: $(GOPATH)
	$(GO) install -v deps
	
build:
	@echo Building with GOPATH: $(GOPATH)
	$(GO) build -v deps

clean:
	rm bin/deps
