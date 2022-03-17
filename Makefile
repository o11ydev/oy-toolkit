# This variable is overriden by `nix develop`
O11Y_NIX_SHELL_ENABLED ?= 0

# Command used to run inside a `nix develop` shell.
# HOME is needed for `go build`.
NIX_DEVELOP = nix --extra-experimental-features nix-command develop --extra-experimental-features flakes -i --keep HOME --keep DOCKER_USERNAME --keep DOCKER_PASSWORD --keep DOCKER_REGISTRY --keep DOCKER_REPOSITORY --keep DOCKER_TAG_SUFFIX --keep DOCKER_ORG

# Docker settings
export DOCKER_ORG ?= o11y
export DOCKER_PASSWORD ?= none
export DOCKER_REGISTRY ?= quay.io
export DOCKER_REPOSITORY ?= oy-toolkit
export DOCKER_TAG_SUFFIX ?=
export DOCKER_USERNAME ?= none

# This is true if we are in `nix develop` shell.
ifeq ($(O11Y_NIX_SHELL_ENABLED),1)
all: lint build

.PHONY: build
build: oy-toolkit

.PHONY: fmt
fmt:
	@gofumpt -l -w --extra .
	@goimports -w -local github.com/o11ydev/oy-toolkit .
	@alejandra -q *.nix

.PHONY: lint
lint:
	@golangci-lint run
	@alejandra -q --check *.nix

oy-%: rebuild
	@echo ">> Building oy-$*"
	@nix build ".#oy-$*"

.PHONY: tidy
tidy:
	@go mod tidy

# Shortcut to force running go build each time.
.PHONY: rebuild
rebuild:

.PHONY: publish
publish:
	@echo ">> Creating publishing script"
	@nix build ".#publish-script" -o ./publish.sh
	@echo ">> Running publishing script"
	@bash -eu ./publish.sh

.PHONY: documentation
documentation:
	@echo ">> Generating documentation"
	@nix build ".#documentation" -o ./documentation

.PHONY: vendorhash
vendorhash:
	@go mod vendor
	@nix hash path vendor

# If we are not in a `nix develop` shell, automatically run into it.
else
default:
	@$(NIX_DEVELOP) -c $(MAKE)

%:
	@$(NIX_DEVELOP) -c $(MAKE) $*

shell:
	@$(NIX_DEVELOP)
endif

