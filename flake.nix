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

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
        config = import ./go.nix;
      };
    in rec {
      packages = import ./packages.nix {inherit pkgs;};
      defaultPackage = packages.oy-toolkit;
      devShell = pkgs.mkShell rec {
        buildInputs = [
          pkgs.go
          pkgs.gofumpt
          pkgs.gotestsum
          pkgs.gotools
          pkgs.golangci-lint
          pkgs.git
          pkgs.nix
          pkgs.skopeo
          pkgs.hugo
          pkgs.alejandra
        ];
        # This variable is needed in our Makefile.
        O11Y_NIX_SHELL_ENABLED = "1";
      };
    });
  nixConfig.bash-prompt = "\[nix-develop\]$ ";
}
