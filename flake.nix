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
    defaultPackage = pkgs.buildGoModule rec {
      pname = "o11ytools";
      version = "0.0.1";
      subPackages = [
        "cmd/runtrace"
        "cmd/o11y-collect-prom-data"
      ];

      src = ./.;

      vendorSha256 = null;
    };
    devShell = pkgs.mkShell {
      buildInputs = [
        pkgs.go_1_17
      ];
    };
  });
}
