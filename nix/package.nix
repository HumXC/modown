{ lib
, buildGoModule
, ...
}:
buildGoModule {
  pname = "modown";
  version = "0.0.1";
  src = ./..;
  vendorHash = "sha256-hh3K3q0Al27AN0LihLRuft/jJtT+hqEnzRmsgVqhZXg=";
  meta = with lib; {
    description = "A command line tool for download packwiz mod pack";
    homepage = "https://github.com/humxc/modown";
    license = licenses.mit;
    mainProgram = "modown";
  };
}
