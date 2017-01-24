package anssi

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	db "github.com/linuxisnotunix/Vulnerobot/modules/database"
	"github.com/linuxisnotunix/Vulnerobot/modules/models"
)

const (
	id                 = "ANSSI"
	listURLFormat      = `http://cert.ssi.gouv.fr/site/%dindex.html`
	aviURLFormat       = `http://www.cert.ssi.gouv.fr/site/%s/index.html`
	aviDirectURLFormat = `http://www.cert.ssi.gouv.fr/site/%s/%s.html`
	aviRegex           = `^CERT(FR|A)-%d-AVI-[0-9]+$`
	minYear            = 2000
)

var maxYear = 2017 //Updated at run time

//TODO implement force and update since last view
//TODO implement all dield to store in DB (parseAVI)
//TODO implement match
//TODO implement system to view progess

//ModuleANSSI retrieve information form http://cert.ssi.gouv.fr/
type ModuleANSSI struct {
}

//New constructor for Module
func New(options map[string]string) models.Collector {
	log.WithFields(log.Fields{
		"id":      id,
		"options": options,
	}).Debug("Creating new Module")
	maxYear = time.Now().Year()
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
	log.WithFields(log.Fields{
		"minYear": minYear,
		"maxYear": maxYear,
	}).Infof("Start collect for : '%s'", id)

	o := db.Orm()
	parsedAVI := 0

	for y := minYear; y <= maxYear; y++ {
		tx := o.Begin()
		url := fmt.Sprintf(listURLFormat, y)
		validAVI := regexp.MustCompile(fmt.Sprintf(aviRegex, y))

		log.WithFields(log.Fields{
			"year":    y,
			"yearURL": url,
		}).Debugf("Getting list from : '%s'", url)
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.WithFields(log.Fields{
				"year":    y,
				"yearURL": url,
			}).Warnf("Faild to get list : %v", err)
		} else {
			// Find AVI link
			doc.Find(".corps a.mg, .corps a.ale").Each(func(i int, s *goquery.Selection) {
				title := s.Text()
				if validAVI.MatchString(title) {
					log.Debugf("AVI found %d: %s", i, title)
					avi, err := parseAVI(title)
					if err != nil {
						log.Warnf("Failed to get AVI : %s", title)
					} else {
						tx.Create(avi)
					}
					parsedAVI++
				} else {
					log.Debugf("Skipping id : %s", title)
				}
			})
		}
		tx.Commit()
	}

	log.WithFields(log.Fields{
		"parsedAVI": parsedAVI,
	}).Infof("Ended collect for : '%s'", id)
	return nil
	//return fmt.Errorf("Module '%s' does not implement Collect().", id) //TODO
}

func parseAVI(id string) (*models.AnssiAVI, error) {
	url := fmt.Sprintf(aviDirectURLFormat, id, id)
	log.WithFields(log.Fields{
		"id":  id,
		"url": url,
	}).Debugf("Getting AVI (%s) from : '%s'", id, url)
	doc, err := goquery.NewDocument(url)

	if err != nil {
		log.WithFields(log.Fields{
			"id":  id,
			"url": url,
		}).Warnf("Faild to get AVI : %v", err)
		return nil, err
	}

	// Find the values
	title := doc.Find("head>title").Text()
	headers := make(map[string]string)
	var lastTitle string
	doc.Find("td.corps > div[align='center'] > table  div[align='center'] > table tr[valign='baseline']").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("td:nth-child(1)").Text())

		value := s.Find("td:nth-child(2)").Text()
		if len(title) == 0 { //Keep previous title index but append value //TODO better
			title = lastTitle
			value += ", " + value
		}
		headers[title] = value
		lastTitle = title
	})

	contents := make(map[string]string)
	doc.Find("td.corps > ul").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			contents["Risques"] = s.Text()
		case 1:
			contents["Documentation"] = s.Text()
		}
	})
	doc.Find("td.corps > p").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 5:
			contents["Systèmes"] = s.Text()
		case 6:
			contents["Résumé"] = s.Text()
		}
	})

	log.WithFields(log.Fields{
		"id":      id,
		"url":     url,
		"title":   title,
		"headers": headers,
		//"contents": contents,
	}).Debug("AVI parsed")
	return &models.AnssiAVI{ID: id, Title: title}, nil

}

//List display known AVI stored by this module in DB
func (m *ModuleANSSI) List() error {
	return fmt.Errorf("Module '%s' does not implement List().", id) //TODO
}
