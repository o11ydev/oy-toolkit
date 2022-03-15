# This variable is overriden by `nix develop`
O11Y_NIX_SHELL_ENABLED ?= 0

# Command used to run inside a `nix develop` shell.
# HOME is needed for `go build`.
NIX_DEVELOP = nix --extra-experimental-features nix-command develop --extra-experimental-features flakes -i --keep HOME

# This is true if we are in `nix develop` shell.
ifeq ($(O11Y_NIX_SHELL_ENABLED),1)
all: lint build

.PHONY: build
build: oy-toolkit

.PHONY: fmt
fmt:
	gofumpt -l -w --extra .

.PHONY: lint
lint:
	golangci-lint run

oy-%: rebuild
	@echo ">> Building oy-$*"
	@nix build ".#oy-$*"

.PHONY: tidy
tidy:
	go mod tidy

# Shortcut to force running go build each time.
.PHONY: rebuild
rebuild:

.PHONY: publish
publish:
	nix build ".#publish-script" -o ./publish.sh
	bash -eu ./publish.sh

# If we are not in a `nix develop` shell, automatically run into it.
else
default:
	@$(NIX_DEVELOP) -c $(MAKE)

%:
	@$(NIX_DEVELOP) -c $(MAKE) $*

shell:
	@$(NIX_DEVELOP)
endif

