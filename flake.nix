{
  description = "ghx - GitHub eXtras";

  inputs = {
    # Use latest nixpkgs
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.05";

    # Flake helper utilities
    flake-utils.url = "github:numtide/flake-utils";

    # Go <=> Nix integration
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    # Unified polyglot source formatter
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    # Nix development shell helper
    devshell.url = "github:numtide/devshell";
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix, treefmt-nix, devshell }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        projName = "ghx";

        pkgs = import nixpkgs {
          inherit system;
          overlays = [ devshell.overlays.default gomod2nix.overlays.default ];
        };

        dprintWasmPluginUrl = name: version:
          "https://plugins.dprint.dev/${name}-${version}.wasm";
      in
      {
        # Configure source formatters for project
        formatter = treefmt-nix.lib.mkWrapper pkgs {
          projectRootFile = "flake.nix";
          programs = {
            # no-lambda-pattern-names is needed to preserve self input arg
            deadnix = { enable = true; no-lambda-pattern-names = true; };
            nixpkgs-fmt.enable = true;
            gofumpt.enable = true;
            dprint = {
              enable = true;
              config = {
                includes = [ "**/*.{json,md,toml}" ];
                excludes = [ "flake.lock" ];
                plugins = [
                  (dprintWasmPluginUrl "json" "0.17.4")
                  (dprintWasmPluginUrl "markdown" "0.15.3")
                  (dprintWasmPluginUrl "toml" "0.5.4")
                ];
              };
            };
          };
        };

        # Configure Nix deveopment shell
        devShells.default = pkgs.devshell.mkShell {
          motd = "";
          name = "${projName}-dev";
          packages = [ "go" gomod2nix ];
        };

        # Build executable binaries in this project (all `main` packages)
        packages.default = pkgs.buildGoApplication {
          pname = projName;
          version = pkgs.lib.removeSuffix "\n" (builtins.readFile ./VERSION);
          src = ./.;
          modules = ./gomod2nix.toml;
        };
      });
}
