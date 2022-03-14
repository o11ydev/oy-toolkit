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
      let
        pkgs = import nixpkgs { inherit system; };
      in
      rec {
        packages = import ./packages.nix { inherit pkgs; };
        defaultPackage = packages.oy-runtrace;
        devShell = pkgs.mkShell rec {
          buildInputs = [
            pkgs.go_1_17
            pkgs.gofumpt
            pkgs.golangci-lint
            pkgs.git
            pkgs.strace
            pkgs.nix
          ];
          # This variable is needed in our Makefile.
          O11Y_NIX_SHELL_ENABLED = "1";
        };
      });
  nixConfig.bash-prompt = "\[nix-develop\]$ ";
}
