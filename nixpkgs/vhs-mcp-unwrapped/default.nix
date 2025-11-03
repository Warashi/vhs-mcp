{ lib, buildGoLatestModule }:
buildGoLatestModule {
  pname = "vhs-mcp-unwrapped";
  version = "0.0.1";
  src = ../../.;
  vendorHash = "sha256-PBhfUBbp0n0KWDEuQQND4uIukqyOsTIL0MWQ8PWVXi0=";

  meta = with lib; {
    mainProgram = "vhs-mcp";
    description = "MCP Server for terminal recordings";
    homepage = "https://github.com/Warashi/vhs-mcp";
    license = licenses.mit;
    platforms = platforms.unix;
  };
}
