# This variable is overriden by `nix develop`
O11Y_NIX_SHELL_ENABLED ?= 0

# Command used to run inside a `nix develop` shell.
NIX_DEVELOP = nix develop -i --keep HOME -c

# This is true if we are in `nix develop` shell..
ifeq ($(O11Y_NIX_SHELL_ENABLED),1)
.PHONY: build
build: oy-runtrace

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	gofumpt -l -w .

oy-%: rebuild
	go build ./cmd/oy-$*


# Shortcut to force running go build each time.
.PHONY: rebuild
rebuild:

# If we are not in a `nix develop` shell, automatically run into it.
else
default:
	@$(NIX_DEVELOP) $(MAKE)

%:
	@$(NIX_DEVELOP) $(MAKE) $*
endif

