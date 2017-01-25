package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

// CmdPlugins plugins availables.
var CmdPlugins = cli.Command{
	Name:        "plugins",
	Aliases:     []string{"p"},
	Usage:       "List plugins availables",
	Description: `Ask each modules to describe itself.`,
	Action:      runListDesc,
	Flags:       []cli.Flag{},
}

func runListDesc(c *cli.Context) error {
	return fmt.Errorf("Command not implemented !")
}
