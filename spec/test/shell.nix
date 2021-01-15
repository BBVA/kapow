{ pkgs ? import (builtins.fetchTarball {
      name = "nixos-20.09-2021-01-15";
      url = "https://github.com/nixos/nixpkgs/archive/cd63096d6d887d689543a0b97743d28995bc9bc3.tar.gz";
      sha256 = "1wg61h4gndm3vcprdcg7rc4s1v3jkm5xd7lw8r2f67w502y94gcy";
    }) {} }:
let
   environconfig = pkgs.python38Packages.buildPythonPackage rec {
     pname = "environconfig";
     version = "1.7.0";
   
     src = pkgs.python38Packages.fetchPypi {
       inherit pname version;
       sha256 = "087amqnqsx7d816adszd1424kma1kx9lfnzffr140wvy7a50vi86";
     };
       meta = {
         homepage = "https://github.com/buguroo/environconfig";
         description = "Environment variables made easy";
       };
   };

   pythonDependencies = [
     pkgs.python38Packages.behave
     pkgs.python38Packages.requests
     environconfig
   ];

   nodeDependencies = (pkgs.callPackage ./node-dependencies.nix {});
in
pkgs.mkShell {
  buildInputs = [
    pkgs.python38
    pythonDependencies
    pkgs.gnumake
    pkgs.which
    nodeDependencies.gherkin-lint
  ];
}
