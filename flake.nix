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
        packages = {
          oy-toolkit = pkgs.buildGoModule {
            name = "oy-toolkit";
            src = ./.;
            vendorSha256 = "sha256-0jqcIE+BZtj7Z7/C0sEPx0acSm0Dv6YxqVGklqmLweA=";

            ldflags = [
              "-X github.com/prometheus/common/version.Version=${builtins.readFile ./VERSION}"
              "-X github.com/prometheus/common/version.Branch=n/a"
              "-X github.com/prometheus/common/version.Revision=n/a"
              "-X github.com/prometheus/common/version.BuildUser=${builtins.readFile (pkgs.runCommand "whoami" { } ("whoami > $out"))}"
              "-X github.com/prometheus/common/version.BuildDate=${builtins.readFile (pkgs.runCommand "date" { } ("date -u '+%Y%m%d-%H:%M:%S%p' > $out"))}"
            ];
          };
        };
        defaultPackage = packages.oy-toolkit;
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
