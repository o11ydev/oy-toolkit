{
  packageOverrides = pkgs: {
    go = pkgs.go_1_18.overrideAttrs (oldAttrs: rec {
      version = "1.18";
      src = pkgs.fetchurl {
        url = "https://dl.google.com/go/go${version}.src.tar.gz";
        sha256 = "sha256-OPQj20zINIg/K1I0QoL6ejn7uTZQ3GKhH98L5kCb2tY=";
      };
    });
  };
}
