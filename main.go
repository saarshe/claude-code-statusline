package main

import (
	"fmt"
	"io"
	"os"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/render"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/wizard"
)

var version = "dev"

func run(args []string, stdin io.Reader) (string, int) {
	for _, arg := range args {
		if arg == "--version" {
			return fmt.Sprintf("claude-code-statusline %s\n", version), 0
		}
		if arg == "--help" {
			return "Usage: claude-code-statusline [--version] [--help] [init]\n\n" +
				"  Reads Claude Code JSON from stdin and prints a status line.\n\n" +
				"  init    Run the interactive setup wizard\n", 0
		}
	}

	input, err := schema.Parse(stdin)
	if err != nil {
		return "", 0
	}

	cfg, err := config.LoadFile(config.ConfigPath())
	if err != nil {
		cfg = config.Default()
	}

	return render.Render(input, cfg), 0
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		if err := wizard.Run("", ""); err != nil {
			fmt.Fprintf(os.Stderr, "wizard error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	out, code := run(os.Args[1:], os.Stdin)
	fmt.Print(out)
	os.Exit(code)
}
