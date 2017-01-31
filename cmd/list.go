package cmd

import (
	"os"

	"github.com/linuxisnotunix/Vulnerobot/modules/collectors"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
	"github.com/linuxisnotunix/Vulnerobot/modules/tools"
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
			Name:        "plugins, p",
			Value:       "all",
			Usage:       "Plugins to load (all or separated by comma)",
			Destination: &settings.PluginList,
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "json",
			Usage: "Format to output (ex : csv or json)", //TODO
		},
		cli.StringFlag{
			Name:  "functions",
			Value: "all",
			Usage: "Functions to match from configuration (ex : f1,f5,...)", //TODO
		},
		cli.StringFlag{
			Name:  "components",
			Value: "all",
			Usage: "Components to match from configuration (ex : c1,c5,...)", //TODO
		},
	},
}

func runList(c *cli.Context) error {
	cl := collectors.Init(map[string]interface{}{
		"appList":       ParseConfigurationFlag(),
		"pluginList":    ParsePluginFlag(),
		"outputFormat":  c.String("format"),
		"functionList":  tools.ParseFlagList(c.String("functions")),
		"componentList": tools.ParseFlagList(c.String("components")),
	})

	return cl.List(os.Stdout)
}
