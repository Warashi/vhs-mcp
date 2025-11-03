{
  makeWrapper,
  stdenv,
  vhs,
  vhs-mcp-unwrapped,
  lib,
}:
stdenv.mkDerivation rec {
  pname = "vhs-mcp";
  version = vhs-mcp-unwrapped.version;

  nativeBuildInputs = [ makeWrapper ];

  dontUnpack = true;

  installPhase = ''
    makeWrapper ${lib.getExe vhs-mcp-unwrapped} $out/bin/vhs-mcp --prefix PATH : ${lib.makeBinPath [ vhs ]}
  '';

  meta = with lib; {
    mainProgram = "vhs-mcp";
    description = "MCP Server for VHS recordings";
    homepage = "https://github.com/Warashi/vhs-mcp";
    license = licenses.mit;
    platforms = platforms.unix;
  };
}
