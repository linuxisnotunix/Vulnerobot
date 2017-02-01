package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/linuxisnotunix/Vulnerobot/modules/collectors"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
	"github.com/linuxisnotunix/Vulnerobot/modules/tools"
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
	data, err := ioutil.ReadFile(settings.ConfigPath)
	if err != nil {
		log.Fatalf("Fail to get config file : %v", err)
	}
	cl := collectors.Init(map[string]interface{}{
		"appList": tools.ParseConfiguration(string(data)),
	})
	return cl.Info(os.Stdout) //TODO use a format output
}
