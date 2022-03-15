{ pkgs, ... }:
with pkgs; [
  (
    go.overrideAttrs
  (
    old: rec {
#      version = "1.17.8";
#      src = fetchurl {
#        url = "https://dl.google.com/go/go${version}.src.tar.gz";
#        sha256 = "sha256-Lv/NiYFA2nmgYfN4TKT42LE9gR+yq+na0kBEQtq733o=";
#      };
    }
    )
    )
]
