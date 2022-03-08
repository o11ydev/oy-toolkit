# This variable is overriden by `nix develop`
O11Y_NIX_SHELL_ENABLED ?= 0

# This is true if we are in `nix develop` shell..
ifeq ($(O11Y_NIX_SHELL_ENABLED),1)
build:
	go version

test:
	go test

# If we are not in a `nix develop` shell, automatically run into it.
else
default:
	@nix develop -c $(MAKE)

%:
	@nix develop -c $(MAKE) $*
endif

