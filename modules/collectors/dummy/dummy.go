package dummy

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/linuxisnotunix/Vulnerobot/modules/models"
)

const id = "DUMMY"

//ModuleDUMMY dummy module do nothing
type ModuleDUMMY struct {
}

//New constructor for Module
func New(options map[string]string) models.Collector {
	log.WithFields(log.Fields{
		"id":      id,
		"options": options,
	}).Debug("Creating new Module")
	return &ModuleDUMMY{}
}

//ID Return the id of the module
func (m *ModuleDUMMY) ID() string {
	return id
}

//IsAvailable Return the availability of the module
func (m *ModuleDUMMY) IsAvailable() bool {
	return false //TODO
}

//Collect collect and parse data to put in database
func (m *ModuleDUMMY) Collect() error {
	return fmt.Errorf("Module '%s' does not implement Collect().", id) //TODO
}

//List display known CVE stored by this module in DB
func (m *ModuleDUMMY) List() error {
	return fmt.Errorf("Module '%s' does not implement List().", id) //TODO
}
