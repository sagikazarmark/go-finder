{
  description = "Finder library for Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    devenv.url = "github:cachix/devenv";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devenv.flakeModule
      ];

      systems = [ "x86_64-linux" "aarch64-darwin" ];

      perSystem = { config, self', inputs', pkgs, system, ... }: {
        devenv.shells = {
          default = {
            languages = {
              go.enable = true;
            };

            packages = with pkgs; [
              go-task
            ];

            # https://github.com/cachix/devenv/issues/528#issuecomment-1556108767
            containers = pkgs.lib.mkForce { };
          };
        };
      };
    };
}
