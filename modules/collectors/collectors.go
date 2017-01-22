package collectors

import (
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/linuxisnotunix/Vulnerobot/modules/collectors/anssi"
	"github.com/linuxisnotunix/Vulnerobot/modules/collectors/dummy"
	"github.com/linuxisnotunix/Vulnerobot/modules/models"
)

var (
	listCollector = []func(map[string]string) models.Collector{anssi.New, dummy.New}
)

//CollectorList list of collector
type CollectorList struct {
	list map[string]models.Collector
}

//Init a collector list and init them
func Init(options map[string]string) *CollectorList {
	return &CollectorList{
		list: getCollectors(options),
	}
}

//getCollectors Return collector list initalized
func getCollectors(options map[string]string) map[string]models.Collector {
	l := make(map[string]models.Collector, len(listCollector)) //TODO only init collectors based on options args ?
	for _, builder := range listCollector {
		collector := builder(options)
		l[collector.ID()] = collector
	}
	return l
}

//Collect ask modules to collect and parse data to put in database
func (cl *CollectorList) Collect() error {
	for id, collector := range cl.list {
		if collector != nil {
			if err := collector.Collect(); handleCollectorError(err) != nil {
				return err
			}
		} else {
			log.Debug("Skipping empty module ", id, " !")
		}
	}
	return nil
}

//List ask module to display known CVE stored by them in DB
func (cl *CollectorList) List() error {
	for id, collector := range cl.list {
		if collector != nil {
			if err := collector.List(); handleCollectorError(err) != nil {
				return err
			}
		} else {
			log.Debug("Skipping empty module ", id, " !")
		}
	}
	return nil
}

//handleCollectorError handle some common errors
func handleCollectorError(err error) error {
	if err != nil {
		if strings.Contains(err.Error(), "does not implement") {
			log.Warn(err)
			return nil
		}
	}
	return err
}
