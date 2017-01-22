package cmd

import (
	log "github.com/Sirupsen/logrus"

	"github.com/urfave/cli"
)

// CmdCollect represents the available update sub-command.
var CmdCollect = cli.Command{
	Name:        "collect",
	Usage:       "Collect CVE from modules and add them to database",
	Description: `Ask each modules to update their database of known vulnerability based on application list.`,
	Action:      runCollect,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "data/configuration",
			Usage: "Application list to monitor",
		},
	},
}

func runCollect(c *cli.Context) error {
	log.Info("Command not implemented !")
	return nil
}
