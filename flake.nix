{
  description = "Nova Cloud development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {self, nixpkgs, flake-utils}:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        go-migrate = pkgs.go-migrate.overrideAttrs(oldAttrs: {
            tags = ["sqlite"];
        });
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            go-migrate
            protobuf
            protoc-gen-go
            protoc-gen-go-grpc
            mockgen
            sqlc
          ];
        };
      });
}
