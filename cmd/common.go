package cmd

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
	"github.com/linuxisnotunix/Vulnerobot/modules/tools"
)

//ParsePluginFlag get flag plugin and return a formatted obj
func ParsePluginFlag() *hashset.Set {
	if settings.PluginList != "" && settings.PluginList == "all" {
		return nil
	}
	ret := strings.Split(settings.PluginList, ",")
	l := hashset.New()
	for _, val := range ret {
		l.Add(strings.ToLower(strings.TrimSpace(val)))
	}
	return l
}

//ParseConfigurationFlag get flag config and return a formatted obj
func ParseConfigurationFlag() *arraylist.List {
	data, err := ioutil.ReadFile(settings.ConfigPath)
	if err != nil {
		log.Fatalf("Fail to get config file : %v", err)
	}
	return tools.ParseConfiguration(string(data))
}
