package cmd

import (
	"github.com/linuxisnotunix/Vulnerobot/modules/collectors"
	"github.com/urfave/cli"
)

// CmdList list collected vulnerability based on a list of apps.
var CmdList = cli.Command{
	Name:        "list",
	Aliases:     []string{"l"},
	Usage:       "List known CVE in database from a application list",
	Description: `Ask each modules to list known vulnerability in database based on application list.`,
	Action:      runList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "data/configuration",
			Usage: "Application list to monitor",
		},
	},
}

func runList(c *cli.Context) error {
	cl := collectors.Init(nil)
	return cl.List()
}
