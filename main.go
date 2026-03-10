package main

import (
	"fmt"
	"os"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/render"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/wizard"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		if err := wizard.Run("", ""); err != nil {
			fmt.Fprintf(os.Stderr, "wizard error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	input, err := schema.Parse(os.Stdin)
	if err != nil {
		os.Exit(0)
	}

	cfg, err := config.LoadFile(config.ConfigPath())
	if err != nil {
		// Bad config: fall back to defaults silently
		cfg = config.Default()
	}

	fmt.Print(render.Render(input, cfg))
}
