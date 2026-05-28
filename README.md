# gk

A terminal-based task and habit manager written in **Go**. The app runs as a TUI
built with Bubble Tea and stores data locally in a **libSQL** database.

## Features

- kanban board with `To Do`, `In Progress`, and `Done` columns
- create, edit, delete, and move tasks between columns
- weekly habit tracker
- mark habits as completed for a selected day
- track the current streak for each habit
- local data storage in `~/.gk`

## Requirements

- Go `1.24.2` or newer
- a terminal that supports TUI applications
- optional: [Task](https://taskfile.dev/) for running commands from `Taskfile.yml`
- optional: `golangci-lint` for linting

## Usage

```sh
go run .
```

Build the binary:

```sh
go build -o gk
./gk
```

With Task:

```sh
task build
./gk
```

## Keybindings

| Key            | Action                                                 |
| -------------- | ------------------------------------------------------ |
| `n`            | create a task or habit in the active section           |
| `e`            | edit the selected task                                 |
| `d`            | delete the selected task                               |
| `enter`        | move a task to the next column or toggle a habit check |
| `tab`          | switch between the task board and habit tracker        |
| `up` / `k`     | move up                                                |
| `down` / `j`   | move down                                              |
| `left` / `h`   | previous column or previous day in the habit tracker   |
| `right` / `l`  | next column or next day in the habit tracker           |
| `esc`          | return from a form                                     |
| `?`            | show help                                              |
| `q` / `ctrl+c` | quit                                                   |

## Data

The app creates the `~/.gk` directory and stores:

- `gk.db` - local database with tasks, habits, and habit checks
- `debug.log` - Bubble Tea debug log used while the app is running

## Tests and Linting

```sh
go test ./...
```

Or with Task:

```sh
task test
task lint
```

## Stack

- Go
- Bubble Tea
- Bubbles
- Lip Gloss
- go-libsql
