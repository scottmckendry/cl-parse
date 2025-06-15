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
        version = "0.5.1"; # x-release-please-version
        src = pkgs.fetchFromGitHub {
          owner = "scottmckendry";
          repo = "cl-parse";
          rev = "v0.5.1"; # x-release-please-version
          sha256 = "sha256-3e4sB62wYLKmexfSqhC6i+xkdGJrOPAgeEYWgfpT8bY";
        };
        vendorHash = "sha256-NsiJoag3wZyIZYdY3H9LAcTKDAJhDv5CLYKHNCtKdi0";
        goPackagePath = "github.com/scottmckendry/cl-parse";
        go = pkgs.go_1_24;
        doCheck = false; # skip tests (these rely on external network calls)

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
