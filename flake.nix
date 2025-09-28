{
  description = "Parse most changelog formats into common structured data formats like JSON, TOML, and YAML";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      packages.${system}.default =
        let
          baseVersion = (
            let
              cmdFile = builtins.readFile ./cmd/cmd.go;
              m = builtins.match ".*const VERSION = \"([0-9]+\.[0-9]+\.[0-9]+)\".*" cmdFile;
            in
            if m == null then "0.0.0" else builtins.elemAt m 0
          );
          shortRev = builtins.substring 0 7 (self.rev or "dev");
        in
        pkgs.buildGoModule {
          pname = "cl-parse";
          version = "${baseVersion}-${shortRev}";
          src = self;
          vendorHash = "sha256-/SL6FE1rX1xkJ6vpVJUms7HUnLNY6qq66jaYUbGWKsM";
          goPackagePath = "github.com/scottmckendry/cl-parse";
          go = pkgs.go_1_24;
          doCheck = true;

          nativeBuildInputs = [ pkgs.installShellFiles ];
          postInstall = ''
            installShellCompletion --cmd cl-parse \
              --bash <($out/bin/cl-parse completion bash) \
              --zsh <($out/bin/cl-parse completion zsh) \
              --fish <($out/bin/cl-parse completion fish)
          '';

          meta = with pkgs.lib; {
            description = "Parse most changelog formats into common structured data formats like JSON, TOML, and YAML";
            homepage = "https://github.com/scottmckendry/cl-parse";
            license = licenses.mit;
            maintainers = [ "scottmckendry" ];
          };
        };
      defaultPackage.${system} = self.packages.${system}.default;
    };
}
