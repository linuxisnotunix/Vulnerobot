package cmd

import (
	"io/ioutil"
	"log"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
	"github.com/linuxisnotunix/Vulnerobot/modules/tools"
)

//ParsePluginFlag get flag plugin and return a formatted obj
func ParsePluginFlag() *hashset.Set {
	return tools.ParseFlagList(settings.PluginList)
}

//ParseConfigurationFlag get flag config and return a formatted obj
func ParseConfigurationFlag() *arraylist.List {
	data, err := ioutil.ReadFile(settings.ConfigPath)
	if err != nil {
		log.Fatalf("Fail to get config file : %v", err)
	}
	return tools.ParseConfiguration(string(data))
}

//TODO store in BD config ?
