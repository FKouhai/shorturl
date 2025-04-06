{
  description = "nix flake that builds url-shortner and provides the development environment used";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    gomod2nix,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [gomod2nix.overlays.default];
          config.allowUnfree = true;
        };
        shorturl = pkgs.buildGoModule {
          name = "shorturl";
          src = ./.;
          vendorHash = "sha256-T6/aGn/UG329oAjFuSqHp/jmY2ZCbhhxdA1c9gkM4Cs=";
          proxyVendor = true;
          doCheck = false;
          postInstall = ''
            mv $out/bin/url_shortener $out/bin/shorturl
          '';
        };
        dockerImage = pkgs.dockerTools.buildImage {
          name = "shorturl";
          tag = "latest";
          created = "now";
          copyToRoot = [shorturl];
          config = {
            Cmd = ["${shorturl}/bin/shorturl"];
          };
        };
      in
        with pkgs; {
          inherit shorturl dockerImage;
          defaultPackage = shorturl;
          devShells.default = mkShell {
            buildInputs = [
              go
              air
              sqlite
              sqlc
            ];
          };
        }
    );
}
