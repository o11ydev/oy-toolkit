{ pkgs, ... }: with pkgs;
let
  basepkg = name: buildGo117Module.override
    {
      go = go.overrideAttrs (oldAttrs: rec {
        version = "1.17.8";
        src = fetchurl {
          url = "https://dl.google.com/go/go${version}.src.tar.gz";
          sha256 = "sha256-Lv/NiYFA2nmgYfN4TKT42LE9gR+yq+na0kBEQtq733o=";
        };
      });
    }
    {
      name = name;
      src = ./.;
      vendorSha256 = "sha256-0jqcIE+BZtj7Z7/C0sEPx0acSm0Dv6YxqVGklqmLweA=";
      subPackages = if name == "oy-toolkit" then [ ] else [ "./cmd/${name}" ];

      ldflags = [
        "-X github.com/prometheus/common/version.Version=${builtins.readFile ./VERSION}"
        "-X github.com/prometheus/common/version.Branch=n/a"
        "-X github.com/prometheus/common/version.Revision=n/a"
        "-X github.com/prometheus/common/version.BuildUser=${builtins.readFile (pkgs.runCommand "whoami" { } ("whoami > $out"))}"
        "-X github.com/prometheus/common/version.BuildDate=${builtins.readFile (pkgs.runCommand "date" { } ("date -u '+%Y%m%d-%H:%M:%S%p' > $out"))}"
      ];
    };
  packageList = (builtins.mapAttrs
    (name: value:
      basepkg name
    )
    (builtins.readDir ./cmd));
in
lib.recursiveUpdate { oy-toolkit = (basepkg "oy-toolkit"); }
  (lib.recursiveUpdate packageList
    (lib.mapAttrs'
      (name: value:
        lib.nameValuePair
          ("docker-${name}")
          (pkgs.dockerTools.buildImage {
            name = name;
            tag = "latest";
            contents = builtins.getAttr name packageList;
          }))
      (builtins.readDir ./cmd)))
