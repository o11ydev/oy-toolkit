# This variable is overriden by `nix develop`
O11Y_NIX_SHELL_ENABLED ?= 0

# Command used to run inside a `nix develop` shell.
# HOME is needed for `go build`.
NIX_DEVELOP = nix --extra-experimental-features nix-command develop --extra-experimental-features flakes -i --keep HOME

PROMLDFLAGS = \
	-X github.com/prometheus/common/version.Version=$(shell cat VERSION) \
	-X github.com/prometheus/common/version.Revision=$(shell git rev-parse --short HEAD) \
	-X github.com/prometheus/common/version.Branch=$(shell git branch | cut -c 3-) \
	-X github.com/prometheus/common/version.BuildUser=$(shell whoami) \
	-X github.com/prometheus/common/version.BuildDate=$(shell date -u '+%Y%m%d-%H:%M:%S%p') \

# This is true if we are in `nix develop` shell.
ifeq ($(O11Y_NIX_SHELL_ENABLED),1)
all: lint test build

.PHONY: build
build: oy-runtrace

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	gofumpt -l -w --extra .

.PHONY: lint
lint:
	golangci-lint run

oy-%: rebuild
	@echo ">> Building oy-$*"
	@nix build ".#oy-$*"
	@echo ">> Running oy-$* --version"
	@nix run ".#oy-$*" -- --version

.PHONY: tidy
tidy:
	go mod tidy

# Shortcut to force running go build each time.
.PHONY: rebuild
rebuild:

# If we are not in a `nix develop` shell, automatically run into it.
else
default:
	@$(NIX_DEVELOP) -c $(MAKE)

%:
	@$(NIX_DEVELOP) -c $(MAKE) $*

shell:
	@$(NIX_DEVELOP)
endif

