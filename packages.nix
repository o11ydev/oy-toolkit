{pkgs, ...}:
with pkgs; let
  basepkg = name:
    buildGoModule {
      name = name;
      src = stdenv.mkDerivation {
        name = "gosrc";
        srcs = [./go.mod ./go.sum ./cmd ./util ./wasm];
        phases = "installPhase";
        installPhase = ''
          mkdir $out
          for src in $srcs; do
            for srcFile in $src; do
              cp -r $srcFile $out/$(stripHash $srcFile)
            done
          done
        '';
      };
      CGO_ENABLED = 0;
      vendorSha256 = "sha256-1tEJR8F8AHjyYyv64zGSHx9O+JyHw0agbV8K/p2FVJ4=";
      #vendorSha256 = pkgs.lib.fakeSha256;
      subPackages =
        if name == "oy-toolkit"
        then []
        else ["./cmd/${name}"];

      ldflags = [
        "-X github.com/prometheus/common/version.Version=${pkgs.lib.removeSuffix "\n" (builtins.readFile ./VERSION)}"
        "-X github.com/prometheus/common/version.Branch=n/a"
        "-X github.com/prometheus/common/version.Revision=n/a"
        "-X github.com/prometheus/common/version.BuildUser=n/a"
        "-X github.com/prometheus/common/version.BuildDate=n/a"
      ];
    };
  packageList =
    builtins.mapAttrs
    (
      name: value:
        basepkg name
    )
    (builtins.readDir ./cmd);
  dockerPackageList =
    lib.mapAttrs'
    (name: value:
      lib.nameValuePair
      "docker-${name}"
      (pkgs.dockerTools.buildImage {
        name = name;
        tag = "latest";
        contents = [pkgs.bashInteractive (builtins.getAttr name packageList)];
        config = {
          Entrypoint = ["/bin/${name}"];
        };
      }))
    (builtins.readDir ./cmd);
in
  lib.recursiveUpdate
  (lib.recursiveUpdate packageList dockerPackageList)
  rec {
    oy-toolkit = basepkg "oy-toolkit";
    publish-script = (
      stdenv.mkDerivation {
        name = "release-script";
        phases = "buildPhase";
        buildPhase =
          pkgs.writeShellScript "publish" ''
          ''
          + (
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
      let
        theme = pkgs.fetchzip {
          url = "https://github.com/thegeeklab/hugo-geekdoc/releases/download/v0.29.4/hugo-geekdoc.tar.gz";
          sha256 = "sha256-Sypg9jBNbW4IeHTqDAq9ZpxgweW1BmFRFjDF51NSg/M=";
          stripRoot = false;
        };
        menu = {
          main = [
            {
              name = "tools";
              sub = builtins.map (x: {
                name = x;
                ref = "/" + x;
              }) (builtins.attrNames packageList);
            }
            {
              name = "HTTP client";
              ref = "/httpclient";
            }
            {
              name = "/metrics lint";
              ref = "/metricslint";
            }
          ];
        };
        menuFile = pkgs.writeTextFile {
          name = "menu";
          text = builtins.toJSON menu;
        };
        commandDocs = stdenv.mkDerivation {
          name = "commandDocs";
          phases = "buildPhase";
          buildPhase =
            pkgs.writeShellScript "commandDocs" ''
              mkdir $out
            ''
            + (
              pkgs.lib.concatMapStrings (x: "\n" + x)
              (
                builtins.attrValues (
                  builtins.mapAttrs
                  (name: value: let
                    description = pkgs.writeTextFile {
                      name = "description";
                      text = pkgs.lib.removePrefix "# ${name}\n" (builtins.readFile (./cmd + "/${name}/README.md"));
                    };
                  in (import ./tool-documentation.nix {
                    tool = builtins.getAttr name packageList;
                    name = name;
                    description = description;
                    pkgs = pkgs;
                  }))
                  packageList
                )
              )
            );
        };
        metricsLint = buildGoModule rec {
          name = "metricslint";
          src = stdenv.mkDerivation {
            name = "gosrc";
            srcs = [./go.mod ./go.sum ./cmd ./util ./wasm];
            phases = "installPhase";
            installPhase = ''
              mkdir $out
              for src in $srcs; do
                for srcFile in $src; do
                  cp -r $srcFile $out/$(stripHash $srcFile)
                done
              done
            '';
          };
          CGO_ENABLED = 0;
          vendorSha256 = "sha256-1tEJR8F8AHjyYyv64zGSHx9O+JyHw0agbV8K/p2FVJ4=";
          #vendorSha256 = pkgs.lib.fakeSha256;
          subPackages = ["wasm/${name}"];
          preBuild = ''
            export GOOS=js
            export GOARCH=wasm
          '';
        };
      in
        stdenv.mkDerivation {
          name = "documentation";
          src = ./docs;
          buildInputs = [pkgs.hugo];
          buildPhase = pkgs.writeShellScript "hugo" ''
            set -e
            cp ${metricsLint}/bin/js_wasm/metricslint static/metricslint.wasm
            chmod -x static/metricslint.wasm
            cp ${pkgs.go_1_18}/share/go/misc/wasm/wasm_exec.js static

            mkdir -p data/menu
            cp ${menuFile} data/menu/main.yml
            cp -r ${commandDocs}/* content
            cat data/menu/main.yml
            hugo --theme=${theme} -d $out
            echo o11y.tools > $out/CNAME
          '';
          installPhase = "true";
        }
    );
    nfpmPackages = let
      npmConfigurations = (
        builtins.mapAttrs (name: value:
          pkgs.writeTextFile {
            name = "npm-config-${name}";
            text = builtins.toJSON {
              name = name;
              arch = "amd64";
              version = pkgs.lib.removeSuffix "\n" (builtins.readFile ./VERSION);
              maintainer = "Julien Pivotto <roidelapluie@inuits.eu>";
              description = "The o11y toolkit is a collection of tools that are useful to manage and run an observability stack.";
              vendor = "o11y";
              contents = [
                {
                  src = (builtins.getAttr name packageList) + "/bin/${name}";
                  dst = "/bin/${name}";
                }
              ];
            };
          })
        (builtins.readDir ./cmd)
      );
    in
      stdenv.mkDerivation {
        name = "gosrc";
        buildInputs = [pkgs.nfpm oy-toolkit];
        phases = "installPhase";
        installPhase =
          ''
            mkdir $out
          ''
          + pkgs.lib.concatMapStrings (x: "\n" + x)
          (
            builtins.attrValues (
              builtins.mapAttrs
              (name: value: ''
                nfpm package --config ${(builtins.getAttr name npmConfigurations)} -p rpm -t $out/${name}.rpm
                nfpm package --config ${(builtins.getAttr name npmConfigurations)} -p deb -t $out/${name}.deb
              '')
              (builtins.readDir ./cmd)
            )
          );
      };
  }
