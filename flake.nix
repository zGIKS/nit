{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-linux"
        "aarch64-darwin"
      ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          nit = pkgs.buildGoModule {
            pname = "nit";
            version = "dev";
            src = ./.;

            vendorHash = "sha256-D49hZAGzz7JPeGvG2l5ax2YotM8E1Ek1smlVH43gLjc=";

            meta = with pkgs.lib; {
              website = "https://github.com/zGIKS/nit";
              license = licenses.mit;
              mainProgram = "nit";
            };
          };
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              gotools
              go-tools
            ];
          };
        }
      );

      defaultPackage = forAllSystems (system: self.packages.${system}.nit);
    };
}
