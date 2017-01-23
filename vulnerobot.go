package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/linuxisnotunix/Vulnerobot/cmd"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
	"github.com/urfave/cli"
)

var (
	//Name Name of app
	Name = "Vulnerobot"
	//Version version of app set by build flag
	Version = "testing"
	//Branch git branch of app set by build flag
	Branch string
	//Commit git commit of app set by build flag
	Commit string
	//BuildTime build time of app set by build flag
	BuildTime string
)

func init() {
	settings.AppName = Name
	settings.AppVersion = Version
	settings.AppBranch = Branch
	settings.AppCommit = Commit
	settings.AppBuildTime = BuildTime
}

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Usage = "Index CVE related to a list of progs"
	app.Version = Version
	cli.VersionPrinter = func(c *cli.Context) {
		if Branch != "" || Commit != "" || BuildTime != "" {
			fmt.Printf("%s == Version: %s - Branch: %s - Commit: %s - BuildTime: %s ==\n", c.App.Name, c.App.Version, Branch, Commit, BuildTime)
		} else {
			fmt.Printf("%s == Version: %s ==\n", c.App.Name, c.App.Version)
		}
	}
	app.Flags = append(app.Flags, cli.BoolFlag{
		Name:        "debug, d",
		Usage:       "Turns on verbose logging",
		EnvVar:      "DEBUG",
		Destination: &settings.AppVerbose,
	})
	app.EnableBashCompletion = true

	app.Before = setup

	app.Commands = []cli.Command{
		cmd.CmdCollect,
		cmd.CmdList,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Fail to run app with %s: %v", os.Args, err)
	}
}
func setup(c *cli.Context) error {
	if settings.AppVerbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return nil
}
