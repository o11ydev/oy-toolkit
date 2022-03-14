{ pkgs, ... }:
builtins.mapAttrs
  (name: value:
    pkgs.buildGoModule {
      name = name;
      src = ./.;
      vendorSha256 = "sha256-0jqcIE+BZtj7Z7/C0sEPx0acSm0Dv6YxqVGklqmLweA=";
      subPackages = [ "./cmd/${name}" ];

      ldflags = [
        "-X github.com/prometheus/common/version.Version=${builtins.readFile ./VERSION}"
        "-X github.com/prometheus/common/version.Branch=n/a"
        "-X github.com/prometheus/common/version.Revision=n/a"
        "-X github.com/prometheus/common/version.BuildUser=${builtins.readFile (pkgs.runCommand "whoami" { } ("whoami > $out"))}"
        "-X github.com/prometheus/common/version.BuildDate=${builtins.readFile (pkgs.runCommand "date" { } ("date -u '+%Y%m%d-%H:%M:%S%p' > $out"))}"
      ];
    })
  (builtins.readDir ./cmd)
