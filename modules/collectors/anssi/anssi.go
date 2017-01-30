package anssi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gosuri/uiprogress"
	db "github.com/linuxisnotunix/Vulnerobot/modules/database"
	"github.com/linuxisnotunix/Vulnerobot/modules/models"
	"github.com/linuxisnotunix/Vulnerobot/modules/tools"
)

const (
	id                 = "ANSSI"
	listURLFormat      = `http://cert.ssi.gouv.fr/site/%dindex.html`
	aviURLFormat       = `http://www.cert.ssi.gouv.fr/site/%s/index.html`
	aviDirectURLFormat = `http://www.cert.ssi.gouv.fr/site/%s/%s.html`
	aviRegex           = `^CERT(FR|A)-([0-9]{4})-AVI-[0-9]+$`
	aviYearRegex       = `^CERT(FR|A)-%d-AVI-[0-9]+$`
	dateRegex          = `^(\d{2}) (.+) (\d{4})$`
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

func parseDate(date string) time.Time {
	validDate := regexp.MustCompile(dateRegex)
	if !validDate.MatchString(date) {
		return time.Time{}
	}
	matchsParu := validDate.FindAllStringSubmatch(date, -1)

	y, err := strconv.ParseInt(matchsParu[0][3], 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"date":       date,
			"matchsParu": matchsParu,
		}).Warnf("%s: Failed to extract year from date", id)
	}
	d, err := strconv.ParseInt(matchsParu[0][1], 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"date":       date,
			"matchsParu": matchsParu,
		}).Warnf("%s: Failed to extract day from date", id)
	}
	return time.Date(int(y), tools.GetFrMonth(matchsParu[0][2]), int(d), 0, 0, 0, 0, time.UTC)
}

//New constructor for Module
func New(options map[string]interface{}) models.Collector {
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

func getLastKnownAVI() string {
	var lastAVI models.AnssiAVI
	db.Orm().Last(&lastAVI)
	log.WithFields(log.Fields{
		"ID": lastAVI.ID,
		//"lastAVI": lastAVI,
	}).Debugf("%s: Getting last AVI in DB", id)
	return lastAVI.ID
}
func listNeededAVI(lastAVIInDB string) (*arraylist.List, error) {
	list := arraylist.New()

	aviMatch := regexp.MustCompile(aviRegex)
	hasStartAvi := aviMatch.MatchString(lastAVIInDB)
	if hasStartAvi { //Match AVI format
		matchs := aviMatch.FindAllStringSubmatch(lastAVIInDB, -1)
		y, err := strconv.ParseInt(matchs[0][2], 10, 64)
		if err != nil {
			log.WithFields(log.Fields{
				"lastAVIInDB": lastAVIInDB,
				"matchs":      matchs,
			}).Warnf("%s: Failed to extract year from lastAVIInDB", id)
		}
		minYear = int(y)
	}
	//minYear = 2016 //Debug force to limit requests
	log.WithFields(log.Fields{
		"minYear": minYear,
		"maxYear": maxYear,
	}).Debugf("%s: Getting list of AVI to collect", id)

	for y := minYear; y <= maxYear; y++ {
		url := fmt.Sprintf(listURLFormat, y)
		validYearAVI := regexp.MustCompile(fmt.Sprintf(aviYearRegex, y))
		log.WithFields(log.Fields{
			"year":    y,
			"yearURL": url,
		}).Debugf("%s: Getting list from : '%s'", id, url)
		doc, err := goquery.NewDocument(url) //Get list of the year
		if err != nil {
			log.WithFields(log.Fields{
				"year":    y,
				"yearURL": url,
			}).Warnf("%s: Failed to get list : %v", id, err)
			continue //Skip if we need blockin should use next line
			//return nil, fmt.Errorf("Module '%s' failed to get list of AVI to parse", id)
		}

		// Find all AVI link
		doc.Find(".corps a.mg, .corps a.ale").Each(func(i int, s *goquery.Selection) {
			avi := s.Text()
			if validYearAVI.MatchString(avi) { //TODO check if newer than lastAVIInDB
				if hasStartAvi && avi <= lastAVIInDB {
					return //Skipping
				}
				list.Add(avi)
			}
		})
	}
	log.WithFields(log.Fields{
		"size": list.Size(),
	}).Debugf("%s: Finish collecting list of AVI", id)
	return list, nil
}

//Collect collect and parse data to put in database
func (m *ModuleANSSI) Collect(bar *uiprogress.Bar) error {
	log.Infof("%s: Start Collect() ", id)

	neededAVI, err := listNeededAVI(getLastKnownAVI())
	if err != nil {
		return err
	}
	if neededAVI.Size() > 0 {
		if bar != nil {
			bar.Total = neededAVI.Size()
		}

		tx := db.Orm().Begin() //Start sql session
		it := neededAVI.Iterator()
		for it.Next() {
			avi, err := parseAVI(it.Value().(string))
			if err != nil {
				log.Warnf("Failed to get AVI : %s", it.Value().(string))
			} else {
				tx.Save(avi)
			}
			if bar != nil {
				bar.Incr()
			}
		}
		tx.Commit() //Commit session
	} else {
		//Nothing to do
		log.Infof("%s: No new AVI to collect.", id)
		if bar != nil {
			bar.Incr()
		}
	}
	return nil
}

func parseAVI(AVIid string) (*models.AnssiAVI, error) {
	url := fmt.Sprintf(aviDirectURLFormat, AVIid, AVIid)
	log.WithFields(log.Fields{
		"AVIid": AVIid,
		"url":   url,
	}).Debugf("%s: Getting AVI (%s) from : '%s'", id, AVIid, url)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.WithFields(log.Fields{
			"AVIid": AVIid,
			"url":   url,
		}).Warnf("%s: Faild to get AVI : %v", id, err)
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
	s := doc.Find("td.corps > h1 > a, td.corps > ul, td.corps > p, td.corps > dl")
	//.Each(func(i int, s *goquery.Selection) {
	for i := 0; i < s.Length(); i++ {
		h := s.Eq(i)
		if h.Is("a") {
			if h.Is("#SECTION00010000000000000000") {
				continue //Skip "Gestion du document"
			}
			/*
				h1 := h.Text()
				if strings.Contains(h1, "-") && len(h1) > 4 {
					h1 = h1[4:]
				}
				if (strings.Contains(h1, "1") || strings.Contains(h1, "2") || strings.Contains(h1, "3") || strings.Contains(h1, "4") || strings.Contains(h1, "5") || strings.Contains(h1, "6") || strings.Contains(h1, "7") || strings.Contains(h1, "8")) && len(h1) > 4 {
					h1 = h1[4:]
				}*/
			/*
				ref, exists := h.Attr("id")
				if !exists {
					log.WithFields(log.Fields{
						"AVIid": AVIid,
						"title": h.Text(),
					}).Warnf("%s: Failed to parse AVI header title", id)
					continue
				}
			*/
			h1 := h.Text()
			if strings.Contains(h1, "Système") {
				h1 = "Systèmes"
			}
			if strings.Contains(h1, "sumé") {
				h1 = "Résumé"
			}
			if strings.Contains(h1, "Risque") {
				h1 = "Risque"
			}
			if strings.Contains(h1, "escription") || strings.Contains(h1, "esccription") {
				h1 = "Description"
			}
			if strings.Contains(h1, "Contournement") || strings.Contains(h1, "solution provisoire") {
				h1 = "Contournement"
			}
			if strings.Contains(h1, "Solution") {
				h1 = "Solution"
			}
			if strings.Contains(h1, "Documentation") || strings.Contains(h1, "document") {
				h1 = "Documentation"
			}
			if strings.Contains(h1, "Vulnérabilité") {
				h1 = "Vulnérabilités"
			}
			if strings.Contains(h1, "Détection") {
				h1 = "Détection"
			}
			if strings.Contains(h1, "Recommandation") {
				h1 = "Recommandation"
			}
			if strings.Contains(h1, "Remerciements") {
				h1 = "Remerciements"
			}

			switch h1 {
			case "Systèmes":
				h1 = "Systèmes"
			case "Résumé":
				h1 = "Résumé"
			case "Risque":
				h1 = "Risque"
			case "description":
			case "Desccription":
			case "Description":
				h1 = "Description"
			case "Recommendation":
				h1 = "Recommendation"
			case "Vulnérabilités":
				h1 = "Vulnérabilités"
			case "Détection":
				h1 = "Détection"
			case "Contournement":
				h1 = "Contournement"
			case "Solution":
				h1 = "Solution"
			case "Documentation":
			case "Gestion détaillée du document":
				h1 = "Documentation"
			case "Recommandation":
				h1 = "Recommandation"
			case "Remerciements":
				h1 = "Remerciements"

			default:
				if len(strings.TrimSpace(h1)) > 0 {
					log.Warnf("Un-catched H1 format : %s", h1)
				}
			}
			n := 1
			content := ""
			for i+n < s.Length() {
				c := s.Eq(i + n)
				if !c.Is("a") {
					content += c.Text()
				} else {
					break
				}
				n++
			}
			//log.Info("ID : ", ref)
			//log.Info("Content", content)
			contents[h1] = content
		}
	}
	//})
	//log.Info(contents)
	log.WithFields(log.Fields{
		"AVIid":   AVIid,
		"url":     url,
		"title":   title,
		"headers": headers,
		//"contents": contents,
	}).Debug("AVI parsed")
	return &models.AnssiAVI{ID: AVIid, Title: title, Risk: contents["Risque"], SystemAffected: contents["Systèmes"], Release: parseDate(headers["Date de la première version"]), LastUpdate: parseDate(headers["Date de la dernière version"])}, nil
}

//List display known AVI stored by this module in DB
func (m *ModuleANSSI) List() error {
	return fmt.Errorf("Module '%s' does not implement List().", id) //TODO
}
