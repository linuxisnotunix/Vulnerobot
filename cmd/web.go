package cmd

import (
	"net/http"
	"strings"

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
			Name:        "listen, l",
			Value:       "127.0.0.1:8080",
			Destination: &settings.WebListen,
			Usage:       "Address and port to listen (ex: 127.0.0.1:8080 or 127.0.0.1:4242 or :8080)",
		}},
}

func runWeb(c *cli.Context) error {

	http.HandleFunc("/public/", server.HandlePublic)
	http.HandleFunc("/api/", server.HandleAPI)

	log.Info("Server running on http://localhost:" + strings.Split(settings.WebListen, ":")[1] + "/public/index.html")

	http.ListenAndServe(settings.WebListen, nil)
	return nil
}
