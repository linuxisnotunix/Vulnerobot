package cmd

import (
	"github.com/linuxisnotunix/Vulnerobot/modules/collectors"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"

	"github.com/urfave/cli"
)

// CmdCollect represents the available update sub-command.
var CmdCollect = cli.Command{
	Name:        "collect",
	Aliases:     []string{"c"},
	Usage:       "Collect CVE from modules and add them to database",
	Description: `Ask each modules to update their database of known vulnerability based on application list.`,
	Action:      runCollect,
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
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Force reload of data",
		},
		cli.BoolFlag{
			Name:        "no-progress",
			Usage:       "Don't display progress bar",
			Destination: &settings.UIDontDisplayProgress,
		},
	},
}

func runCollect(c *cli.Context) error {
	cl := collectors.Init(map[string]interface{}{
		"appList":      ParseConfigurationFlag(),
		"pluginList":   ParsePluginFlag(),
		"forceRefresh": c.Bool("force"),
	})
	return cl.Collect()
	//return nil
}
