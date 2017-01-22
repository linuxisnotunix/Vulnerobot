package main

import (
	"log"
	"os"

	"github.com/linuxisnotunix/Vulnerobot/cmd"
	"github.com/urfave/cli"
)

var (
	//Version version of app set by build flag
	Version = "testing"
	//Branch git branch of app set by build flag
	Branch string
	//Commit git commit of app set by build flag
	Commit string
	//BuildTime build time of app set by build flag
	BuildTime string
)

/*
func init() {
	settings.AppVersion = Version
	settings.AppBranch = Branch
	settings.AppCommit = Commit
	settings.AppBuildTime = BuildTime
}
*/

func main() {
	app := cli.NewApp()
	app.Name = "Vulnerobot"
	app.Usage = "Index CVE related to a list of progs"
	app.Version = Version
	if Branch != "" {
		app.Version += "-" + Branch
	}
	if Commit != "" {
		app.Version += "#" + Commit
	}
	if BuildTime != "" {
		app.Version += "@" + BuildTime
	}
	app.Commands = []cli.Command{
		cmd.CmdCollect,
		cmd.CmdList,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(4, "Fail to run app with %s: %v", os.Args, err)
	}
}
