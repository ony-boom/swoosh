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
          version = "0.1.4";
          vendorHash = "sha256-bwHGOu5EGUU7Uw8Fe5Yswv8tN9uxFgjtVpx4wncmHAI=";
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
          ];
        };

        nixosModules.default =
          { config, pkgs, ... }:
          {
            options.programs.swoosh.enable = {
              type = pkgs.lib.types.bool;
              default = false;
              description = "Enable Swoosh Audio Output Switcher user service";
            };

            config = pkgs.lib.mkIf config.programs.swoosh.enable {
              systemd.user.services.swoosh = {
                description = "Swoosh Audio Output Switcher";
                after = [
                  "graphical-session.target"
                  "pulseaudio.service"
                  "pipewire.service"
                ];
                wants = [
                  "pulseaudio.service"
                  "pipewire.service"
                ];
                partOf = [ "graphical-session.target" ];
                wantedBy = [ "default.target" ];
                serviceConfig = {
                  Type = "simple";
                  ExecStart = "${self.packages.${system}.default}/bin/swoosh";
                  Restart = "on-failure";
                  RestartSec = 3;
                  Environment = "DISPLAY=:0";
                  KillMode = "process";
                };
              };
            };
          };
      }
    );
}
