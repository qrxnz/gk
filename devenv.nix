{
  pkgs,
  lib,
  config,
  ...
}: {
  # https://devenv.sh/packages/
  packages = [
    pkgs.delve
    pkgs.gopls
    pkgs.nixd
    pkgs.nodePackages.prettier
    pkgs.alejandra
    pkgs.taplo
    pkgs.dockfmt
  ];

  # https://devenv.sh/languages/
  languages = {
    go.enable = true;
    nix.enable = true;
  };

  # https://devenv.sh/git-hooks/
  git-hooks = {
    enable = true;
    hooks = {
      treefmt = {
        enable = true;
      };
    };
  };

  # https://devenv.sh/treefmt/
  treefmt = {
    enable = true;
    config = {
      programs = {
        prettier.enable = true;
        alejandra.enable = true;
        taplo.enable = true;
        dockerfmt.enable = true;
      };
    };
  };

  # See full reference at https://devenv.sh/reference/options/
}
