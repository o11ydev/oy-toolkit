{ pkgs, ... }: with pkgs;
let
  basepkg = name: buildGo117Module.override
    {
#      go = go.overrideAttrs (oldAttrs: rec {
#        version = "1.17.8";
#        src = fetchurl {
#          url = "https://dl.google.com/go/go${version}.src.tar.gz";
#          sha256 = "sha256-Lv/NiYFA2nmgYfN4TKT42LE9gR+yq+na0kBEQtq733o=";
#        };
#      });
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
        "-X github.com/prometheus/common/version.BuildUser=n/a"
        "-X github.com/prometheus/common/version.BuildDate=n/a"
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
    buildPhase = pkgs.writeShellScript "publish" ''
    '' + (
      pkgs.lib.concatMapStrings (x: "\n" + x)
        (
          builtins.attrValues (
            builtins.mapAttrs
              (name: value: ''
                echo -e "\n\n## ${name} ##\n" >> $out
                echo 'echo ">> ${name}"' >> $out
                echo 'skopeo --insecure-policy copy --dest-username "$DOCKER_USERNAME" --dest-password "$DOCKER_PASSWORD" docker-archive://${builtins.getAttr name dockerPackageList} docker://$DOCKER_REGISTRY/$DOCKER_ORG/$DOCKER_REPOSITORY:${pkgs.lib.removePrefix "docker-" name}$DOCKER_TAG_SUFFIX' >> $out
              '')
              dockerPackageList
          )
        )
        );
      }
  );
  documentation = (
    let theme = pkgs.fetchzip{
        url = "https://github.com/thegeeklab/hugo-geekdoc/releases/download/v0.27.4/hugo-geekdoc.tar.gz";
        sha256 = "sha256-TtnpqLRaanninztiv85ASEsiO6/ciVmnjS4zotkdCaY=";
        stripRoot=false;
      };
      menu = {
        main = [
          { name = "tools"; sub = (builtins.map(x: {name = x; ref="/"+x;}) (builtins.attrNames packageList)); }
        ];
      };
      menuFile = pkgs.writeTextFile { name = "menu"; text = builtins.toJSON(menu); };
  commandDocs = (stdenv.mkDerivation {
    name = "commandDocs";
    phases = "buildPhase";
    buildPhase = pkgs.writeShellScript "commandDocs" ''
    mkdir $out
    '' + (
      pkgs.lib.concatMapStrings (x: "\n" + x)
        (
          builtins.attrValues (
            builtins.mapAttrs
            (name: value: let
              description = pkgs.writeTextFile { name = "description"; text = (pkgs.lib.removePrefix "# ${name}\n" (builtins.readFile (./cmd + "/${name}/README.md"))); };
            in
              ''
              echo --- >> $out/${name}.md
              echo title: ${name} >> $out/${name}.md
              echo --- >> $out/${name}.md
              echo >> $out/${name}.md
              echo '## Usage' >> $out/${name}.md
              echo >> $out/${name}.md
              echo '```' >> $out/${name}.md
              ${builtins.getAttr name packageList}/bin/${name} --help &>> $out/${name}.md
              echo '```' >> $out/${name}.md
              echo '## Description' >> $out/${name}.md
              cat ${description} >> $out/${name}.md
              echo >> $out/${name}.md
              echo '## Downloading' >> $out/${name}.md
              echo '{{< tabs "usage" >}}' >> $out/${name}.md
              echo '{{< tab "docker" >}}' >> $out/${name}.md
              echo 'To execute **${name}** with docker, run:' >> $out/${name}.md
              echo '```' >> $out/${name}.md
              echo 'docker run quay.io/o11y/oy-toolkit:${name} --help' >> $out/${name}.md
              echo '```' >> $out/${name}.md
              echo '{{< /tab >}}' >> $out/${name}.md
              echo '{{< tab "nix" >}}' >> $out/${name}.md
              echo 'To execute **${name}** with nix, run:' >> $out/${name}.md
              echo '```' >> $out/${name}.md
              echo 'nix run github:o11ydev/oy-toolkit#${name} -- --help' >> $out/${name}.md
              echo '```' >> $out/${name}.md
              echo '{{< /tab >}}' >> $out/${name}.md
              echo '{{< /tabs >}}' >> $out/${name}.md
              echo >> $out/${name}.md
              '')
              packageList
          )
        )
        );
      });
    in
    stdenv.mkDerivation {
    name = "documentation";
    src = ./docs;
    buildInputs = [ pkgs.hugo ];
    buildPhase = pkgs.writeShellScript "hugo" ''
      mkdir -p data/menu
      cp ${menuFile} data/menu/main.yml
      cp -r ${commandDocs}/* content
      cat data/menu/main.yml
      hugo --theme=${theme} -d $out
    '';
    installPhase = "true";
  });
}
