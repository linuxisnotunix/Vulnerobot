package cmd

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/linuxisnotunix/Vulnerobot/modules/server"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
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
			Name:        "config, c",
			Value:       "data/configuration",
			Usage:       "Application list to monitor",
			Destination: &settings.ConfigPath,
		},
		cli.StringFlag{
			Name:        "port, p",
			Value:       ":8080",
			Destination: &settings.WebPort,
			Usage:       "TCP port ot listen",
		}},
}

func runWeb(c *cli.Context) error {

	http.HandleFunc("/public/", server.HandlePublic)
	http.HandleFunc("/api/", server.HandleAPI)

	log.Info("Server running on http://127.0.0.1" + settings.WebPort + "/public/index.html")

	http.ListenAndServe(settings.WebPort, nil)
	return nil
}
