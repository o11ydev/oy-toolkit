{
  pkgs,
  tool,
  name,
  description,
}: let
  usage = pkgs.runCommand "${name} --help" {} "${tool}/bin/${name} --help &> $out";
  documentation = pkgs.writeTextFile {
    name = "${name}-doc.md";
    text = ''
      ---
      title: ${name}
      ---

      ${builtins.readFile ./docs/tools-top.md}

      ## Usage

      ```
      ${builtins.readFile usage}
      ```

      ## Description
      ${builtins.readFile description}

      ## Downloading

      {{< tabs "usage" >}}

      {{< tab "linux (wget)" >}}
      To execute **${name}** within Linux, run:
      ```
      wget https://github.com/o11ydev/oy-toolkit/releases/download/main/${name} -O ${name} && chmod +x ${name} && ./${name} --help
      ```
      {{< /tab >}}

      {{< tab "docker" >}}
      To execute **${name}** with docker, run:
      ```
      docker run quay.io/o11y/oy-toolkit:${name} --help
      ```
      {{< /tab >}}

      {{< tab "nix" >}}
      To execute **${name}** with nix, run:
      ```
      nix run github:o11ydev/oy-toolkit#${name} -- --help
      ```
      {{< /tab >}}

      {{< /tabs >}}
    '';
  };
in "cat ${documentation} >> $out/${name}.md"
