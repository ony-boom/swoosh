{
  description = "Swoosh flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {

        packages.default = pkgs.buildGoModule {
          src = self;
          name = "swoosh";
          version = "0.1.0";
          vendorHash = "sha256-bwHGOu5EGUU7Uw8Fe5Yswv8tN9uxFgjtVpx4wncmHAI=";
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
          ];
        };
      }
    );
}
