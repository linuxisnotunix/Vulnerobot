package anssi

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/linuxisnotunix/Vulnerobot/modules/models"
)

const id = "ANSSI"
const listURLFormat = `http://cert.ssi.gouv.fr/site/%dindex.html`

//ModuleANSSI retrieve information form http://cert.ssi.gouv.fr/
type ModuleANSSI struct {
}

//New constructor for Module
func New(options map[string]string) models.Collector {
	log.WithFields(log.Fields{
		"id":      id,
		"options": options,
	}).Debug("Creating new Module")
	/*
		log.WithFields(log.Fields{
			"Orm": db.Orm(),
		}).Debug("DB Orm debug")
	*/
	return &ModuleANSSI{}
}

//ID Return the id of the module
func (m *ModuleANSSI) ID() string {
	return id
}

//IsAvailable Return the availability of the module
func (m *ModuleANSSI) IsAvailable() bool {
	return false //TODO
}

//Collect collect and parse data to put in database
func (m *ModuleANSSI) Collect() error {
	return fmt.Errorf("Module '%s' does not implement Collect().", id) //TODO
}

//List display known CVE stored by this module in DB
func (m *ModuleANSSI) List() error {
	return fmt.Errorf("Module '%s' does not implement List().", id) //TODO
}
