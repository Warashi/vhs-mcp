{
  ffmpeg,
  makeWrapper,
  stdenv,
  ttyd,
  vhs-unwrapped,
  lib,
}:
stdenv.mkDerivation rec {
  pname = "vhs";
  version = vhs-unwrapped.version;

  nativeBuildInputs = [ makeWrapper ];

  dontUnpack = true;

  installPhase = ''
    makeWrapper ${lib.getExe vhs-unwrapped} $out/bin/vhs --prefix PATH : ${
      lib.makeBinPath [
        ffmpeg
        ttyd
      ]
    }
  '';

  meta = with lib; {
    mainProgram = "vhs";
    description = "Your CLI home video recorder";
    homepage = "https://github.com/charmbracelet/vhs";
    license = licenses.mit;
    platforms = platforms.unix;
  };
}
