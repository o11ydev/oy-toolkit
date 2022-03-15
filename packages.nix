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
      preBuild = "export buildDate=$(date -u '+%Y%m%d-%H:%M:%S%p'); echo OK > VERSION";

      ldflags = [
        "-X github.com/prometheus/common/version.Version=${builtins.readFile ./VERSION}"
        "-X github.com/prometheus/common/version.Branch=n/a"
        "-X github.com/prometheus/common/version.Revision=n/a"
        "-X github.com/prometheus/common/version.BuildUser=n/a}"
        "-X github.com/prometheus/common/version.BuildDate=n/a}"
      ];
    };
  packageList = (builtins.mapAttrs
    (name: value:
      basepkg name
    )
    (builtins.readDir ./cmd));
  dockerPackageList = (lib.mapAttrs'
    (name: value:
      lib.nameValuePair
        ("docker-${name}")
        (pkgs.dockerTools.buildImage {
          name = name;
          tag = "latest";
          contents = [ pkgs.bashInteractive (builtins.getAttr name packageList) ];
          config = {
            Entrypoint = [ "/bin/${name}" ];
          };
        }))
    (builtins.readDir ./cmd));
in
lib.recursiveUpdate
  (lib.recursiveUpdate packageList dockerPackageList)
{
  oy-toolkit = (basepkg "oy-toolkit");
  publish-script = (stdenv.mkDerivation {
    name = "release-script";
    phases = "buildPhase";
    unpackPhase = "true";
    buildPhase = pkgs.writeShellScript "publish" ''
    echo '#!/usr/bin/env bash -e' > $out
    echo 'echo ">> Login"' >> $out
    echo 'skopeo login $DOCKER_REGISTRY -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"' >> $out
  '' + (
      pkgs.lib.concatMapStrings (x: "\n" + x)
        (
          builtins.attrValues (
            builtins.mapAttrs
              (name: value: ''
                echo 'echo ">> ${name}"' >> $out
                echo 'skopeo --insecure-policy copy docker-archive://${builtins.getAttr name dockerPackageList} docker://$DOCKER_REGISTRY/$DOCKER_ORG/${pkgs.lib.removePrefix "docker-" name}:$DOCKER_TAG' >> $out
              '')
              dockerPackageList
          )
     )
     )
        + ''
    echo 'echo ">> Logout"' >> $out
    echo 'skopeo logout $DOCKER_REGISTRY' >> $out
        ''
    ;
    installPhase = "true";
  });
}
