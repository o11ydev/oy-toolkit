{
  pkgs,
  tool,
  name,
  description,
}: let
  usage = pkgs.runCommand "${name} --help" {} "${tool}/bin/${name} --help &> $out";
  documentation = pkgs.writeTextFile {
    name = "${name}-doc.md";
    text = let
      firstParagraphDescription = builtins.elemAt (builtins.split "\n\n" (builtins.readFile description)) 0;
      frontmatterDescription = pkgs.lib.removePrefix "\n*${name}*" firstParagraphDescription;
    in ''
      ---
      title: ${name}
      geekdocRepo: "https://github.com/o11ydev/oy-toolkit"
      geekdocEditPath: "edit/main/cmd/${name}"
      geekdocFilePath: "README.md"
      tool: ${name}
      description: ${assert frontmatterDescription != firstParagraphDescription; builtins.toJSON frontmatterDescription}
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

      {{< tab "linux (deb)" >}}
      **${name}** is available as a `.deb` package:
      https://github.com/o11ydev/oy-toolkit/releases/download/main/${name}.deb
      {{< /tab >}}

      {{< tab "linux (yum)" >}}
      **${name}** is available as a `.rpm` package:
      https://github.com/o11ydev/oy-toolkit/releases/download/main/${name}.rpm
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
