{
  pkgs ? import <nixpkgs> { },
}:
let
  callPackage = pkgs.lib.callPackageWith (pkgs // packages);
  packages = rec {
    default = vhs-mcp;
    vhs = callPackage ./vhs { };
    vhs-mcp = callPackage ./vhs-mcp { };
    vhs-mcp-unwrapped = callPackage ./vhs-mcp-unwrapped { };
    vhs-unwrapped = callPackage ./vhs-unwrapped { };
  };
in
packages
