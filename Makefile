O11Y_NIX_SHELL_ENABLED ?= 0

build: pre
	go version

test: pre
	go test


.PHONY: pre
pre:
ifeq ($(O11Y_NIX_SHELL_ENABLED),0)
	$(error Please run this command in a `nix develop` shell)
endif

