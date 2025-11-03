{
  lib,
  buildGoLatestModule,
  fetchFromGitHub,
}:
buildGoLatestModule rec {
  pname = "vhs-unwrapped";
  version = "v0.10.0";
  src = fetchFromGitHub {
    owner = "charmbracelet";
    repo = "vhs";
    rev = "v${version}";
    hash = "sha256-ZnE5G8kfj7qScsT+bZg90ze4scpUxeC6xF8dAhdUUCo=";
  };
  vendorHash = "sha256-jmabOEFHduHzOBAymnxQrvYzXzxKnS1RqZZ0re3w63Y=";

  meta = with lib; {
    mainProgram = "vhs";
    description = "Your CLI home video recorder";
    homepage = "https://github.com/charmbracelet/vhs";
    license = licenses.mit;
    platforms = platforms.unix;
  };
}
