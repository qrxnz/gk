{
  pkgs,
  lib,
  config,
  ...
}: {
  # https://devenv.sh/languages/
  languages.go.enable = true;

  # https://devenv.sh/packages/
  packages = [
    pkgs.treefmt
    pkgs.prettier
    pkgs.taplo
  ];

  # https://devenv.sh/git-hooks/
  git-hooks.hooks = {
    treefmt.enable = true;
  };

  # https://devenv.sh/reference/options/
  treefmt = {
    enable = true;
    config.programs = {
      gofmt.enable = true;
      prettier.enable = true;
      taplo.enable = true;
    };
  };
}
