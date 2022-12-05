{pkgs, ...}: let
  go = pkgs.go_1_19;
in {
  packageOverrides = pkgs: {
    buildGoModule = pkgs.buildGoModule.override {go = go;};
    go = go;
  };
}
