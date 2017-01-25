package anssi

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gosuri/uiprogress"
	db "github.com/linuxisnotunix/Vulnerobot/modules/database"
	"github.com/linuxisnotunix/Vulnerobot/modules/models"
)

const (
	id                 = "ANSSI"
	listURLFormat      = `http://cert.ssi.gouv.fr/site/%dindex.html`
	aviURLFormat       = `http://www.cert.ssi.gouv.fr/site/%s/index.html`
	aviDirectURLFormat = `http://www.cert.ssi.gouv.fr/site/%s/%s.html`
	aviRegex           = `^CERT(FR|A)-[0-9]{4}-AVI-[0-9]+$`
	aviYearRegex       = `^CERT(FR|A)-%d-AVI-[0-9]+$`
)

var (
	minYear = 2000 //Update at runtime if needed
	maxYear = 2017 //Updated at run time
)

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
	return &ModuleANSSI{}
}

//ID Return the id of the module
func (m *ModuleANSSI) ID() string {
	return id
}

//IsAvailable Return the availability of the module
func (m *ModuleANSSI) IsAvailable() bool {
	return true //TODO
}

func listNeededCVE(lastAVIInDB string) (*arraylist.List, error) {
	list := arraylist.New()
	log.WithFields(log.Fields{
		"minYear": minYear,
		"maxYear": maxYear,
	}).Debugf("Getting list of AVI to collect for : '%s'", id)

	if regexp.MustCompile(aviRegex).MatchString(lastAVIInDB) { //Match AVI format
		minYear = 2000 //TODO setup minYear from lastAVIInDB
	}
	//minYear = 2016 //Debug force to limit requests

	for y := minYear; y <= maxYear; y++ {
		url := fmt.Sprintf(listURLFormat, y)
		validYearAVI := regexp.MustCompile(fmt.Sprintf(aviYearRegex, y))
		log.WithFields(log.Fields{
			"year":    y,
			"yearURL": url,
		}).Debugf("Getting list from : '%s'", url)
		doc, err := goquery.NewDocument(url) //Get list of the year
		if err != nil {
			log.WithFields(log.Fields{
				"year":    y,
				"yearURL": url,
			}).Warnf("Failed to get list : %v", err)
			continue //Skip if we need blockin should use next line
			//return nil, fmt.Errorf("Module '%s' failed to get list of AVI to parse", id)
		}

		// Find all AVI link
		doc.Find(".corps a.mg, .corps a.ale").Each(func(i int, s *goquery.Selection) {
			avi := s.Text()
			if validYearAVI.MatchString(avi) {
				list.Add(avi)
			}
		})
	}
	log.WithFields(log.Fields{
		"size": list.Size(),
	}).Debugf("Finish collecting list of AVI")
	return list, nil
}

//Collect collect and parse data to put in database
func (m *ModuleANSSI) Collect(bar *uiprogress.Bar) error {

	o := db.Orm()
	neededAVI, err := listNeededCVE("") //TODO get last AVI from DB
	if err != nil {
		return err
	}
	bar.Total = neededAVI.Size()

	tx := o.Begin() //Start sql session
	it := neededAVI.Iterator()
	for it.Next() {
		avi, err := parseAVI(it.Value().(string))
		if err != nil {
			log.Warnf("Failed to get AVI : %s", it.Value().(string))
		} else {
			tx.Create(avi)
		}
		bar.Incr()
	}
	tx.Commit() //Commit session
	return nil
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
