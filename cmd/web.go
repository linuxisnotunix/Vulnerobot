package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

// CmdWeb Start a web server to display result.
var CmdWeb = cli.Command{
	Name:        "web",
	Aliases:     []string{"w"},
	Usage:       "Start a web server to display result.",
	Description: ``,
	Action:      runWeb,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "data/configuration",
			Usage: "Application list to monitor",
		},
		cli.StringFlag{
			Name:  "port, p",
			Value: ":8080",
			Usage: "TCP port ot listen",
		}},
}

func runWeb(c *cli.Context) error {
	return fmt.Errorf("Command not implemented !")
}
