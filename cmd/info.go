package cmd

import (
	"os"

	"github.com/linuxisnotunix/Vulnerobot/modules/collectors"
	"github.com/urfave/cli"
)

// CmdInfo plugins availables.
var CmdInfo = cli.Command{
	Name:        "info",
	Aliases:     []string{"i"},
	Usage:       "Display global info like the of list plugins availables",
	Description: `Ask each modules to describe itself.`,
	Action:      runListDesc,
	Flags:       []cli.Flag{},
}

func runListDesc(c *cli.Context) error {
	cl := collectors.Init(map[string]interface{}{})
	return cl.Info(os.Stdout) //TODO use a format output
}
