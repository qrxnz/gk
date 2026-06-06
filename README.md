# gk

<h3 align="center">
  <div>
    <a href="https://github.com/qrxnz/gk/issues">
        <img src="https://img.shields.io/github/issues/qrxnz/gk?color=fab387&labelColor=303446&style=for-the-badge">
    </a>
    <a href="https://github.com/qrxnz/gk/stargazers">
        <img src="https://img.shields.io/github/stars/qrxnz/gk?color=ca9ee6&labelColor=303446&style=for-the-badge">
    </a>
    <a href="https://github.com/qrxnz/gk/actions/workflows/go.yml">
        <img src="https://img.shields.io/github/actions/workflow/status/qrxnz/gk/go.yml?color=a6e3a1&labelColor=303446&style=for-the-badge&label=Go">
    </a>
    <a href="https://github.com/qrxnz/gk/blob/main/.github/LICENCE">
        <img src="https://img.shields.io/static/v1.svg?style=for-the-badge&label=License&message=MIT&logoColor=ca9ee6&colorA=313244&colorB=cba6f7"/>
    </a>
    <br>
    </div>
</h3>

> A terminal-based task and habit manager written in Go

**gk** is a lightweight, fast, and keyboard-driven TUI designed for managing daily tasks and habits, built in **Go** using the **Bubble Tea** framework and a local **libSQL** database. Built with developers and terminal enthusiasts in mind, it allows you to stay organized without ever leaving the command line or lifting your hands off the keyboard.

## 🧰 Features

- kanban board with `To Do`, `In Progress`, and `Done` columns
- create, edit, delete, and move tasks between columns
- weekly habit tracker
- mark habits as completed for a selected day
- track the current streak for each habit
- local data storage in `~/.gk`

## 📋 Requirements

- Go `1.26.2` or newer
- a terminal that supports TUI applications

## 🛠️ Installation

### 📦 Binary Releases

Pre-compiled binaries for Linux, Windows, and macOS are available on the [Releases](https://github.com/qrxnz/gk/releases) page.

### 🐹Using Go

You can install `gk` directly using `go install`:

```bash
go install github.com/qrxnz/gk@latest
```

### 🏗️ Build from Source

To build from source, you need to have [Go](https://go.dev/) installed.

```bash
git clone https://github.com/qrxnz/gk.git
cd gk
go build -o gk .
```

Alternatively, if you have [Task](https://taskfile.dev/) installed, you can use:

```bash
task build
```

### ❄️ Using Nix

- **Run without installing**

```bash
nix run github:qrxnz/gk
```

- **Add to a Nix Flake**

Add input in your flake like:

```nix
{
 inputs = {
   gk = {
     url = "github:qrxnz/gk";
     inputs.nixpkgs.follows = "nixpkgs";
   };
 };
}
```

With the input added you can reference it directly:

```nix
{ inputs, system, ... }:
{
  # NixOS
  environment.systemPackages = [ inputs.gk.packages.${pkgs.system}.default ];
  # home-manager
  home.packages = [ inputs.gk.packages.${pkgs.system}.default ];
}
```

- **Install imperatively**

```bash
nix profile install github:qrxnz/gk
```

## 📖 Usage

Simply run the application from your terminal:

```sh
gk
```

### ⌨️ Keybindings

| Key            | Action                                                 |
| -------------- | ------------------------------------------------------ |
| `n`            | create a task or habit in the active section           |
| `e`            | edit the selected task/habit                           |
| `d`            | delete the selected task/habit                         |
| `enter`        | move a task to the next column or toggle a habit check |
| `tab`          | switch between the task board and habit tracker        |
| `up` / `k`     | move up                                                |
| `down` / `j`   | move down                                              |
| `left` / `h`   | previous column or previous day in the habit tracker   |
| `right` / `l`  | next column or next day in the habit tracker           |
| `esc`          | return from a form                                     |
| `?`            | show help                                              |
| `q` / `ctrl+c` | quit                                                   |

## 👨🏻‍💻 Development

This project uses [Nix](https://nixos.org/) with flakes and [direnv](https://direnv.net/) to provide a reproducible development environment.

1. **Clone the repository**

   ```sh
   git clone https://github.com/qrxnz/gk.git
   cd gk
   ```

1. **Activate the environment**
   If you have Nix and direnv installed, the environment should be activated automatically when you enter the directory. If not, run:

   ```sh
   direnv allow
   ```

1. **Available Commands**
   This project uses `go-task` as a command runner. Here are the most common commands:
   - `task build`: Build a production binary.
   - `task test`: Run unit tests.
   - `task lint`: Run the linter and fix issues.

## 🗒️ Credits

### 🎨 Inspiration

I was inspired by:

- [charm-and-friends/kancli](https://github.com/charm-and-friends/kancli)

## 📜 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.
