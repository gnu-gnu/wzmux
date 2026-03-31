# wzmux

[한국어](README.md)

WezTerm multiplexer for Claude Code. Launch, monitor, and manage multiple Claude Code agents from your terminal.

## Features

- **Multi-agent management** — Spawn multiple Claude Code sessions in WezTerm tabs
- **Live dashboard** — Bubble Tea TUI showing real-time agent status
- **Session resume** — Pick up where you left off by agent name
- **Auto status tracking** — Tab titles update with status emoji (⚙ running, 🟡 waiting, ✅ done, 🔴 error)
- **Single binary** — No runtime dependencies besides WezTerm

## Requirements

- [WezTerm](https://wezfurlong.org/wezterm/)
- [Claude Code](https://docs.anthropic.com/en/docs/claude-code)

## Install

```bash
go install github.com/gnugnu/wzmux@latest
```

Or build from source:

```bash
git clone https://github.com/gnugnu/wzmux.git
cd wzmux
go install .
```

## Setup

Register Claude Code hooks (one-time):

```bash
wzmux setup
```

This adds hook entries to `~/.claude/settings.json`. Your existing settings are preserved and backed up to `settings.json.bak`.

## Usage

### Launch a new agent

```bash
wzmux new backend
wzmux new frontend "refactor the login page"
```

Each agent gets a dedicated WezTerm tab with a dashboard split pane. Duplicate names are auto-suffixed (`backend`, `backend-1`, `backend-2`, ...).

### Resume a previous session

```bash
wzmux resume backend
```

Continues the Claude Code conversation from where it left off.

### List active agents

```bash
wzmux ls
```

### Open the dashboard

```bash
wzmux dashboard
```

Press `1`-`9` to jump to an agent tab. Press `q` to quit.

### Kill agents

```bash
wzmux kill backend
wzmux kill --all
```

## How it works

```
Claude Code agent
       ↓ (hook events)
wzmux hook <event>          ← Go binary handles hooks directly
       ↓
/tmp/claude-agent-status/  ← status JSON files
       ↑
wzmux dashboard             ← reads status, renders TUI
```

1. `wzmux setup` registers `wzmux hook <event>` as Claude Code hooks
2. When Claude Code runs inside WezTerm, hooks fire and wzmux updates the tab title + writes a status file
3. The dashboard polls status files every 2 seconds
4. Outside WezTerm, hooks silently do nothing

## Uninstall

```bash
wzmux uninstall
```

Removes wzmux hooks from `~/.claude/settings.json` and cleans up status files. Other settings are untouched.

## License

MIT
