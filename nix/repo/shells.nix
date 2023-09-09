/*
This file holds reproducible shells with commands in them.

They conveniently also generate config files in their startup hook.
*/
{
  mdbook = inputs.std.lib.dev.mkShell {nixago = [(inputs.std.lib.cfg.mdbook cell.config.mdbook)];};
  # Tool Homepage: https://numtide.github.io/devshell/
  default = inputs.std.lib.dev.mkShell {
    name = "Paisano TUI";

    # Tool Homepage: https://nix-community.github.io/nixago/
    # This is Standard's devshell integration.
    # It runs the startup hook when entering the shell.
    nixago = [
      inputs.std.lib.cfg.adrgen
      inputs.std.lib.cfg.conform
      (inputs.std.lib.cfg.treefmt cell.config.treefmt)
      (inputs.std.lib.cfg.editorconfig cell.config.editorconfig)
      (inputs.std.lib.cfg.githubsettings cell.config.githubsettings)
      (inputs.std.lib.cfg.lefthook cell.config.lefthook)
      (inputs.std.lib.cfg.mdbook cell.config.mdbook)
    ];

    commands =
      [
        {
          category = "release";
          package = inputs.nixpkgs.cocogitto;
        }
        {
          category = "rendering";
          package = inputs.nixpkgs.mdbook;
        }
        {
          package = inputs.nixpkgs.delve;
          category = "dev";
          name = "dlv";
        }
        {
          package = inputs.nixpkgs.go;
          category = "dev";
        }
        {
          package = inputs.nixpkgs.gotools;
          category = "dev";
        }
        {
          package = inputs.nixpkgs.gopls;
          category = "dev";
        }
        {
          package = inputs.cells.tui.app.default;
          category = "bleeding edge";
        }
      ]
      ++ inputs.nixpkgs.lib.optionals inputs.nixpkgs.stdenv.isLinux [
        {
          package = inputs.nixpkgs.golangci-lint;
          category = "dev";
        }
      ];
    env =
      [
        {
          name = "GO11MODULE";
          value = "auto";
        }
      ]
      ++ inputs.nixpkgs.lib.optionals inputs.nixpkgs.stdenv.isDarwin [
        {
          name = "PATH";
          # the pinned version of devshell or std doesn't support `prefix`
          # that well (it got resolved into `:$PATH`)
          eval = let
            inherit (inputs.nixpkgs) xcbuild;
            xcbuild_path = inputs.nixpkgs.lib.makeBinPath [
              xcbuild
              "${xcbuild}/Toolchains/XcodeDefault.xctoolchain"
            ];
          in "${xcbuild_path}:$PATH";
        }
      ];
  };

  trad = let
    inherit (inputs.nixpkgs) mkShell go gotools gopls cocogitto mdbook delve;
    app = inputs.cells.tui.app.default;
  in
    mkShell {
      buildInputs = [go gotools gopls cocogitto mdbook delve app];
    };
}
