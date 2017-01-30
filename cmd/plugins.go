package cmd

import (
	"github.com/linuxisnotunix/Vulnerobot/modules/collectors"
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
	cl := collectors.Init(map[string]interface{}{})
	return cl.Info() //TODO use a format output
}
