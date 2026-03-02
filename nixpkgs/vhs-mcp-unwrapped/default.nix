{ lib, buildGoLatestModule }:
buildGoLatestModule {
  pname = "vhs-mcp-unwrapped";
  version = "0.0.1";
  src = ../../.;
  vendorHash = "sha256-97y0kB3+k1rX8OEhVjE9+NzJil8tm5Ce/Y7WKJ/hAVc=";

  meta = with lib; {
    mainProgram = "vhs-mcp";
    description = "MCP Server for terminal recordings";
    homepage = "https://github.com/Warashi/vhs-mcp";
    license = licenses.mit;
    platforms = platforms.unix;
  };
}
