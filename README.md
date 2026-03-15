# claude-code-statusline

A customizable status line for [Claude Code](https://docs.anthropic.com/en/docs/claude-code). Shows model, context usage, tokens, cost, git status, and more — right in your terminal.

## Install

```sh
go install github.com/saars/claude-code-statusline@latest
```

Or download a binary from [Releases](https://github.com/saars/claude-code-statusline/releases).

## Setup

Run the interactive wizard:

```sh
claude-code-statusline init
```

This lets you pick which components to show, choose a theme, and automatically wires it into Claude Code's `settings.json`.

## Configuration

Config lives at `~/.claude-code-statusline/config.toml`. The wizard generates this for you, but you can edit it manually.

## License

[MIT](LICENSE)
