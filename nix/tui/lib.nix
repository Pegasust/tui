let
  inherit (inputs) nixpkgs;
  inherit (nixpkgs.lib) licenses;
  paisano-homepage = "https://github.com/paisano-nix/tui";
  paisano-description = (import (inputs.self + /flake.nix)) description;
  paisano-license = licenses.unlicense;
  paisano-version = "0.15.0+dev";
in {
  rebrand-paisano = {
    license ? paisano-license,
    description ? paisano-description,
    homepage ? paisano-homepage,
    pname ? "paisano",
    nativeBuildInputs ? [nixpkgs.installShellFiles],
    postInstall ? ''
      installShellCompletion --cmd paisano \
        --bash <($out/bin/paisano _carapace bash) \
        --fish <($out/bin/paisano _carapace fish) \
        --zsh <($out/bin/paisano _carapace zsh)
    '',
    version ? paisano-version,
  }:
    nixpkgs.buildGoModule rec {
      inherit version postInstall nativeBuildInputs;
      pname = "paisano";
      meta = {
        inherit homepage description license;
      };

      src = inputs.self + /src;

      vendorHash = "sha256-ja0nFWdWqieq8m6cSKAhE1ibeN0fODDCC837jw0eCnE=";

      ldflags = [
        "-s"
        "-w"
        "-X main.buildVersion=${version}"
      ];
    };
}
