package nvd

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gosuri/uiprogress"
	db "github.com/linuxisnotunix/Vulnerobot/modules/database"
	"github.com/linuxisnotunix/Vulnerobot/modules/models"
	feed "github.com/linuxisnotunix/Vulnerobot/modules/models/xsd/nvd.nist.gov/schema/nvd-cve-feed_2.0.xsd_go"
)

const (
	id                   = "NVD"
	yearXMLURLFormat     = `https://static.nvd.nist.gov/feeds/xml/cve/nvdcve-2.0-%d.xml.gz`
	yearXMLMETAURLFormat = `https://static.nvd.nist.gov/feeds/xml/cve/2.0/nvdcve-2.0-%s.meta`
	nvdURLFormat         = `https://web.nvd.nist.gov/view/vuln/detail?vulnId=%s`
	timeFormat           = time.RFC3339Nano
)

var (
	minYear = 2002
	maxYear = 2017 //Updated at run time
)

//ModuleNVD NVD/Mitre module cparse data from https://nvd.nist.gov/download.cfm
type ModuleNVD struct {
}

//New constructor for Module
func New(options map[string]string) models.Collector {
	log.WithFields(log.Fields{
		"id":      id,
		"options": options,
	}).Debug("Creating new Module")
	maxYear = time.Now().Year()
	return &ModuleNVD{}
}

//ID Return the id of the module
func (m *ModuleNVD) ID() string {
	return id
}

//IsAvailable Return the availability of the module
func (m *ModuleNVD) IsAvailable() bool {
	return true //TODO
}

//Collect collect and parse data to put in database
func (m *ModuleNVD) Collect(bar *uiprogress.Bar) error {
	neededList := listUpdatedList()
	if neededList.Size() > 0 {
		if bar != nil {
			bar.Total = neededList.Size()
		}
		tx := db.Orm().Begin() //Start sql session
		it := neededList.Iterator()
		for it.Next() {
			cveList, err := getList(it.Value().(int))
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Warnf("%s: Failed to get CVE List : %d", id, it.Value().(int))
			} else {
				//tx.Create(avi)
				for _, cve := range cveList {
					tx.Create(cve)
				}
				cveListMeta, err := getListMeta(strconv.Itoa(it.Value().(int)))
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Warnf("%s: Failed to get CVE List Meta : %d", id, it.Value().(int))
				} else {
					tx.Create(cveListMeta) //Store in DB for history
				}

			}
			if bar != nil {
				bar.Incr()
			}
		}
		tx.Commit() //Commit session
	} else {
		//Nothing to do
		log.Infof("%s: No new CVE to collect.", id)
	}
	return nil
	//return fmt.Errorf("Module '%s' does not implement Collect().", id)
}

/*
type xmlList struct {
	CVEList []xmlCVE `xml:"http://scap.nist.gov/schema/feed/vulnerability/2.0 entry"`
}

//xmlCVE represente a xmlCVE stored in xml
type xmlCVE struct {
	ID string `xml:"vuln:cve-id"`
	//Summary    string
	Release    string `xml:"vuln:published-datetime"`
	LastUpdate string `xml:"vuln:last-modified-datetime"`
}
*/
func parseList(year int, r io.Reader) ([]models.NvdCVE, error) {
	var l feed.TxsdNvd

	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("Failed to read list data. (read-list) %v", err)
	}
	err = xml.Unmarshal(bs, &l)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode list data. (parse-list) %v", err)
	}
	results := make([]models.NvdCVE, len(l.XsdGoPkgHasElems_Entry.Entries))

	for i := 0; i < len(l.XsdGoPkgHasElems_Entry.Entries); i++ {
		entrie := l.XsdGoPkgHasElems_Entry.Entries[i]

		release, err := time.Parse(timeFormat, entrie.XsdGoPkgHasElem_PublishedDatetimesequencevulnerabilityTypeschema_PublishedDatetime_XsdtDateTime_.PublishedDatetime.String())
		if err != nil {
			log.WithFields(log.Fields{
				"timeFormat":   timeFormat,
				"release_text": entrie.XsdGoPkgHasElem_PublishedDatetimesequencevulnerabilityTypeschema_PublishedDatetime_XsdtDateTime_.PublishedDatetime.String(),
				"error":        err,
			}).Errorf("%s: Failed to decode Release time format. %v", id, err)
		}
		lastUpdt, err := time.Parse(timeFormat, entrie.XsdGoPkgHasElem_LastModifiedDatetimesequencevulnerabilityTypeschema_LastModifiedDatetime_XsdtDateTime_.LastModifiedDatetime.String())
		if err != nil {
			log.WithFields(log.Fields{
				"timeFormat":    timeFormat,
				"lastUpdt_text": entrie.XsdGoPkgHasElem_LastModifiedDatetimesequencevulnerabilityTypeschema_LastModifiedDatetime_XsdtDateTime_.LastModifiedDatetime.String(),
				"error":         err,
			}).Errorf("%s: Failed to decode LastUpdate time format %v", id, err)
		}

		results[i] = models.NvdCVE{
			ID:         entrie.XsdGoPkgHasElem_CveIdchoicesequencevulnerabilityTypeschema_CveId_CveTcveNamePatternType_.CveId.String(), //entrie.XsdGoPkgHasAttr_Id_TvulnerabilityIdType_.Id.String(),
			Summary:    entrie.XsdGoPkgHasElem_SummarysequencevulnerabilityTypeschema_Summary_XsdtString_.Summary.String(),
			Release:    release,
			LastUpdate: lastUpdt,
		}
	}

	log.WithFields(log.Fields{
		"year":    year,
		"count":   len(results),
		"xml":     l,
		"results": results,
	}).Debugf("%s: Finish to get CVE List content", id)
	return results, nil
}
func getList(year int) ([]models.NvdCVE, error) {
	response, err := http.Get(fmt.Sprintf(yearXMLURLFormat, year))
	if err != nil {
		return nil, fmt.Errorf("Failed to get list information. (get-list) %v", err)
	}
	defer response.Body.Close()
	respUnZip, err := gzip.NewReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to get list information. (gunzip-list) %v", err)
	}
	return parseList(year, respUnZip)
}
func parseListMeta(year string, data string) (*models.NvdList, error) {
	m := &models.NvdList{
		Year: year,
	}
	lines := strings.Split(data, "\r\n")
	if len(lines) != 5 {
		return nil, fmt.Errorf("Failed to parse meta file (lines)!")
	}
	for i := range lines {
		fields := strings.SplitN(lines[i], ":", 2)
		if len(fields) != 2 {
			return nil, fmt.Errorf("Failed to parse meta file (fields)!")
		}
		switch fields[0] {
		case "lastModifiedDate":
			m.LastModifiedDate = fields[1]
		case "sha256":
			m.Sha256 = fields[1]
		}
	}
	return m, nil
}
func getListMeta(year string) (*models.NvdList, error) {
	response, err := http.Get(fmt.Sprintf(yearXMLMETAURLFormat, year))
	if err != nil {
		return nil, fmt.Errorf("Failed to get list information. (get-meta) %v", err)
	}
	defer response.Body.Close()
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to get list information. (read-meta) %v", err)
	}
	return parseListMeta(year, string(bs))
}
func hasBeenUpdated(l models.NvdList) bool {
	log.WithFields(log.Fields{
		"year": l.Year,
	}).Debugf("%s: Check update of list via meta", id)
	newL, err := getListMeta(l.Year)
	if err != nil {
		log.WithFields(log.Fields{
			"year":  l.Year,
			"error": err,
		}).Warnf("%s: Failed to get list information.", id)
		return true
	}
	return strings.Compare(l.Sha256, newL.Sha256) != 0
}
func listUpdatedList() *arraylist.List {
	log.WithFields(log.Fields{
		"minYear": minYear,
		"maxYear": maxYear,
	}).Infof("%s: Getting list of CVE to collect", id)
	list := arraylist.New()
	for y := minYear; y <= maxYear; y++ {
		var listInDB models.NvdList
		db.Orm().Where("year = ?", y).First(&listInDB)
		if listInDB.Year == strconv.Itoa(y) {
			//Check if object has been updated
			if hasBeenUpdated(listInDB) {
				list.Add(y) //List updated we need to parse it
			}
		} else {
			list.Add(y) //Not found in DB so we need it
		}
	}
	return list
}

//List display known CVE stored by this module in DB
func (m *ModuleNVD) List() error {
	return fmt.Errorf("Module '%s' does not implement List().", id)
}
