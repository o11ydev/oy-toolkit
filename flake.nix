{
  description = "O11ytools";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
  flake-utils.lib.eachDefaultSystem (system:
  let pkgs = import nixpkgs { inherit system; };
  in
  {
    devShell = pkgs.mkShell {
      buildInputs = [
        pkgs.go_1_17
        pkgs.gofumpt
        pkgs.golangci-lint
        pkgs.git
      ];
      # This variable is needed in our Makefile.
      O11Y_NIX_SHELL_ENABLED = "1";
    };
  });
}
