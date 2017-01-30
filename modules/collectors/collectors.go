package collectors

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gosuri/uiprogress"

	"github.com/linuxisnotunix/Vulnerobot/modules/collectors/nvd"
	"github.com/linuxisnotunix/Vulnerobot/modules/models"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
	"github.com/willf/pad"
)

var (
	//listCollector = []func(map[string]string) models.Collector{anssi.New, dummy.New, nvd.New}
	listCollector = []func(map[string]interface{}) models.Collector{nvd.New}
)

//CollectorList list of collector
type CollectorList struct {
	options map[string]interface{}
	list    map[string]models.Collector
}

//Init a collector list and init them
func Init(options map[string]interface{}) *CollectorList {
	return &CollectorList{
		options: options,
		list:    getCollectors(options),
	}
}

//getCollectors Return collector list initalized
func getCollectors(options map[string]interface{}) map[string]models.Collector {
	l := make(map[string]models.Collector, len(listCollector)) //TODO only init collectors based on options args ?
	for _, builder := range listCollector {
		collector := builder(options)
		l[collector.ID()] = collector
	}
	return l
}

//Collect ask modules to collect and parse data to put in database
func (cl *CollectorList) Collect() error {
	var p *uiprogress.Progress
	if !settings.UIDontDisplayProgress {
		p = uiprogress.New()
		p.RefreshInterval = 50 * time.Millisecond
		p.Start()
		log.SetOutput(p.Bypass())
	}
	var wg sync.WaitGroup
	for id, collector := range cl.list {
		if err := executeCollectorCollect(p, &wg, id, collector); err != nil {
			//return err
			log.Fatalf("A uncatched error occured : %s", err.Error())
		}
	}
	wg.Wait()
	if !settings.UIDontDisplayProgress {
		time.Sleep(p.RefreshInterval) //Wait to UI finish
		log.SetOutput(os.Stdout)
		p.Stop()
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

//executeCollector start the collector
func executeCollectorCollect(p *uiprogress.Progress, wg *sync.WaitGroup, id string, collector models.Collector) error {
	wg.Add(1)
	if collector != nil {
		log.Info("Starting module ", id, " ...")
		var bar *uiprogress.Bar
		if p != nil {
			bar = p.AddBar(1).AppendCompleted().PrependElapsed()
			bar.PrependFunc(func(b *uiprogress.Bar) string {
				return pad.Right(id, 5, " ")
				//fmt.Sprintf("%s", id)
			})
			bar.AppendFunc(func(b *uiprogress.Bar) string {
				return fmt.Sprintf("(%d/%d)", b.Current(), b.Total)
			})
		}
		go func() {
			defer wg.Done()
			if err := collector.Collect(bar); err != nil {
				handleCollectorError(err)
			}
		}()
	} else {
		log.Debug("Skipping empty module ", id, " !")
	}
	return nil
}

//handleCollectorError handle some common errors
func handleCollectorError(err error) error {
	if err != nil {
		if strings.Contains(err.Error(), "does not implement") {
			log.Warnf("%s", err.Error())
			return nil
		}
		log.Fatalf("Unhandled error : %s", err.Error())
	}
	return err
}
