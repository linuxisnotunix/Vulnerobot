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
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/hashset"
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
	opts map[string]interface{}
}

//New constructor for Module
func New(options map[string]interface{}) models.Collector {
	log.WithFields(log.Fields{
		"id":      id,
		"options": options,
	}).Debug("Creating new Module")
	maxYear = time.Now().Year()
	return &ModuleNVD{opts: options}
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
	log.Infof("%s: Start Collect() ", id)

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
					//tx.Create(cve)
					tx.Save(cve)
				}
				cveListMeta, err := getListMeta(strconv.Itoa(it.Value().(int)))
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Warnf("%s: Failed to get CVE List Meta : %d", id, it.Value().(int))
				} else {
					//tx.Create(cveListMeta) //Store in DB for history
					tx.Save(cveListMeta) //Store in DB for history
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
		if bar != nil {
			bar.Incr()
		}
	}
	return nil
	//return fmt.Errorf("Module '%s' does not implement Collect().", id)
}

func parseList(year int, r io.Reader) ([]models.NvdCVE, error) {
	var l feed.TxsdNvd
	//listCPE := arraylist.New()
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
		var cpes []models.NvdCPE
		if entrie.XsdGoPkgHasElem_VulnerableSoftwareListsequencevulnerabilityTypeschema_VulnerableSoftwareList_TvulnerableSoftwareType_.VulnerableSoftwareList != nil && entrie.XsdGoPkgHasElem_VulnerableSoftwareListsequencevulnerabilityTypeschema_VulnerableSoftwareList_TvulnerableSoftwareType_.VulnerableSoftwareList.XsdGoPkgHasElems_ProductsequencevulnerableSoftwareTypeschema_Product_CpeLangTnamePattern_.Products != nil {
			cpes = make([]models.NvdCPE, len(entrie.XsdGoPkgHasElem_VulnerableSoftwareListsequencevulnerabilityTypeschema_VulnerableSoftwareList_TvulnerableSoftwareType_.VulnerableSoftwareList.XsdGoPkgHasElems_ProductsequencevulnerableSoftwareTypeschema_Product_CpeLangTnamePattern_.Products))
			for i, val := range entrie.XsdGoPkgHasElem_VulnerableSoftwareListsequencevulnerabilityTypeschema_VulnerableSoftwareList_TvulnerableSoftwareType_.VulnerableSoftwareList.XsdGoPkgHasElems_ProductsequencevulnerableSoftwareTypeschema_Product_CpeLangTnamePattern_.Products {
				cpes[i] = models.NvdCPE{
					ID: val.String(),
				}
				/*
					if !list.Contains(val.String()) {
						listCPE.Add(val.String())
					}
				*/
			}
		}
		results[i] = models.NvdCVE{
			ID:           entrie.XsdGoPkgHasElem_CveIdchoicesequencevulnerabilityTypeschema_CveId_CveTcveNamePatternType_.CveId.String(), //entrie.XsdGoPkgHasAttr_Id_TvulnerabilityIdType_.Id.String(),
			Summary:      entrie.XsdGoPkgHasElem_SummarysequencevulnerabilityTypeschema_Summary_XsdtString_.Summary.String(),
			Release:      release,
			AffectedCPEs: cpes,
			LastUpdate:   lastUpdt,
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
func hasBeenUpdated(old, new *models.NvdList, err error) bool {
	if err != nil {
		log.WithFields(log.Fields{
			"year":  old.Year,
			"error": err,
		}).Warnf("%s: Failed to get list information.", id)
		return true
	}
	log.WithFields(log.Fields{
		"year": old.Year,
	}).Debugf("%s: Check update of list via meta", id)

	return strings.Compare(old.Sha256, new.Sha256) != 0
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
			newList, err := getListMeta(listInDB.Year)
			if hasBeenUpdated(&listInDB, newList, err) {
				log.WithFields(log.Fields{
					"year":     y,
					"previous": listInDB,
					"new":      newList,
				}).Infof("%s: Adding year to list of CVE to collect since it has changed", id)
				list.Add(y) //List updated we need to parse it
			}
		} else {
			log.WithFields(log.Fields{
				"year": y,
			}).Infof("%s: Adding year to list of CVE to collect since we have no trace of it", id)
			list.Add(y) //Not found in DB so we need it
		}
	}
	return list
}

func (m *ModuleNVD) elToMatch() *arraylist.List {
	log.WithFields(log.Fields{
		"options": m.opts,
	}).Debugf("%s: Start elToMatch()", id)
	l := arraylist.New()
	listComponents := hashmap.New()
	listMatchs := hashmap.New()

	if m.opts["functionList"] != nil && m.opts["functionList"].(*hashset.Set) != nil || m.opts["componentList"] != nil && m.opts["componentList"].(*hashset.Set) != nil {
		if m.opts["functionList"] != nil && m.opts["functionList"].(*hashset.Set) != nil {
			for _, fp := range m.opts["functionList"].(*hashset.Set).Values() {
				for _, app := range m.opts["appList"].(*arraylist.List).Values() {
					_, okf := app.(map[string]string)["Function"]
					if okf && strings.Contains(strings.ToLower(app.(map[string]string)["Function"]), strings.ToLower(fp.(string))) { //TODO better like split by comma ...
						c := app.(map[string]string)
						if _, ok := listComponents.Get(c["CPE"]); ok {
							lMatch, _ := listMatchs.Get(c["CPE"])
							lMatch.(*arraylist.List).Add(models.ResponseListMatch{
								Type:  "Function",
								Value: fp.(string),
							})
						} else {
							listComponents.Put(c["CPE"], c)
							lMatch := arraylist.New()
							lMatch.Add(models.ResponseListMatch{
								Type:  "Function",
								Value: fp.(string),
							})
							listMatchs.Put(c["CPE"], lMatch)
						}
					}
				}
			}
		}

		if m.opts["componentList"] != nil && m.opts["componentList"].(*hashset.Set) != nil {
			for _, cp := range m.opts["componentList"].(*hashset.Set).Values() {
				for _, app := range m.opts["appList"].(*arraylist.List).Values() {
					_, okv := app.(map[string]string)["Version"]
					_, okn := app.(map[string]string)["Name"]
					if okn && okv && strings.Contains(strings.ToLower(app.(map[string]string)["Name"]+":"+app.(map[string]string)["Version"]), strings.ToLower(cp.(string))) { //TODO better match for components
						c := app.(map[string]string)
						if _, ok := listComponents.Get(c["CPE"]); ok {
							lMatch, _ := listMatchs.Get(c["CPE"])
							lMatch.(*arraylist.List).Add(models.ResponseListMatch{
								Type:  "Component",
								Value: cp.(string),
							})
						} else {
							listComponents.Put(c["CPE"], c)
							lMatch := arraylist.New()
							lMatch.Add(models.ResponseListMatch{
								Type:  "Component",
								Value: cp.(string),
							})
							listMatchs.Put(c["CPE"], lMatch)
						}
					}
				}
			}
		}
		//Reconstruct
		for _, c := range listComponents.Values() {
			lMatch, _ := listMatchs.Get(c.(map[string]string)["CPE"])
			matchs := make([]models.ResponseListMatch, lMatch.(*arraylist.List).Size())
			for i, d := range lMatch.(*arraylist.List).Values() {
				matchs[i] = d.(models.ResponseListMatch)
			}
			l.Add(models.ResponseListEntry{
				Component: c.(map[string]string),
				Matchs:    matchs,
			})
		}
		return l
	}
	//By default Load all components by defaults
	for _, app := range m.opts["appList"].(*arraylist.List).Values() {
		l.Add(models.ResponseListEntry{
			Component: app.(map[string]string),
			Matchs:    []models.ResponseListMatch{},
		})
	}
	return l

}

//List display known CVE stored by this module in DB
func (m *ModuleNVD) List() (*arraylist.List, error) {
	l := arraylist.New()
	nbVulns := 0
	log.Infof("%s: Start List() ", id)
	els := m.elToMatch()
	for _, b := range els.Values() {
		el := b.(models.ResponseListEntry)
		cpes := make([]models.NvdCPE, 0)
		db.Orm().Where("id LIKE ?", el.Component["CPE"]+"%").Find(&cpes)
		//.Preload("NvdCPE", "id LIKE ?", el.Component["CPE"]+"%").Find(&vulns)
		resVulns := arraylist.New()
		for _, cpe := range cpes {
			log.WithFields(log.Fields{
				"cpe": el.Component["CPE"],
			}).Infof("%s: Found CPE : %s", id, cpe.ID)
			vulns := make([]models.NvdCVE, 0)
			db.Orm().Model(&cpe).Association("RelatedCVEs").Find(&vulns)
			for _, vuln := range vulns {
				log.WithFields(log.Fields{
					"cpe": el.Component["CPE"],
				}).Debugf("%s: Found Vuln : %s", id, vuln.ID)
				resVulns.Add(models.ResponseListVuln{
					Source: id,
					Value: map[string]string{
						"ID":      vuln.ID,
						"CPE":     cpe.ID,
						"Summary": vuln.Summary,
						"URL":     fmt.Sprintf(nvdURLFormat, vuln.ID),
					},
				})
				nbVulns++
			}
		}
		//Switch format (casting)
		vs := make([]models.ResponseListVuln, resVulns.Size())
		for i, v := range resVulns.Values() {
			vs[i] = v.(models.ResponseListVuln)
		}
		l.Add(models.ResponseListEntry{
			Component: el.Component,
			Matchs:    el.Matchs,
			Vulns:     vs,
		})
	}
	log.Infof("%s: Found %d vulns", id, nbVulns)
	return l, nil
	//return fmt.Errorf("Module '%s' does not implement List().", id)
}
