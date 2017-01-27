package cmd

import (
	"github.com/linuxisnotunix/Vulnerobot/modules/collectors"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
	"github.com/linuxisnotunix/Vulnerobot/modules/tools"

	"io/ioutil"

	log "github.com/Sirupsen/logrus"
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
			Name:  "plugins, p",
			Value: "all",
			Usage: "Plugins to load",
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
	data, err := ioutil.ReadFile(settings.ConfigPath)
	if err != nil {
		log.Fatalf("Fail to get config file : %v", err)
	}
	log.Info(string(data))
	tableauConfig := tools.ParseConfiguration(string(data))
	log.Info("Debug: Configuration : ", tableauConfig)
	cl := collectors.Init(nil)
	return cl.Collect()
	//return nil
}
