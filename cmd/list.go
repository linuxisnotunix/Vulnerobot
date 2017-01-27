package cmd

import (
	"github.com/linuxisnotunix/Vulnerobot/modules/collectors"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
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
			Name:        "config, c",
			Value:       "data/configuration",
			Usage:       "Application list to monitor",
			Destination: &settings.ConfigPath,
		},
		cli.StringFlag{
			Name:  "plugins, p",
			Value: "all",
			Usage: "Plugins to load (ex : p1,p4,...)",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "csv",
			Usage: "Format to output (ex : csv or json)",
		},
		cli.StringFlag{
			Name:  "functions",
			Value: "all",
			Usage: "Functions to match from configuration (ex : f1,f5,...)",
		},
		cli.StringFlag{
			Name:  "components",
			Value: "all",
			Usage: "Components to match from configuration (ex : c1,c5,...)",
		},
	},
}

func runList(c *cli.Context) error {
	cl := collectors.Init(nil)
	return cl.List()
}
