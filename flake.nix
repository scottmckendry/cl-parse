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
      packages.${system}.default = pkgs.buildGoModule {
        pname = "cl-parse";
        version = "0.6.0"; # x-release-please-version
        src = pkgs.fetchFromGitHub {
          owner = "scottmckendry";
          repo = "cl-parse";
          rev = "v0.6.0"; # x-release-please-version
          sha256 = "sha256-3AyMSKmNofI2+ezP8mgD/lmt2m8Zm0pNma6p7qAfFJ4";
        };
        vendorHash = "sha256-/SL6FE1rX1xkJ6vpVJUms7HUnLNY6qq66jaYUbGWKsM";
        goPackagePath = "github.com/scottmckendry/cl-parse";
        go = pkgs.go_1_24;
        doCheck = false; # skip tests (these rely on external network calls)

        # Install generated shell completions (bash, zsh, fish)
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
