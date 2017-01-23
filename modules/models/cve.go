package models

import "github.com/jinzhu/gorm"

//https://cve.mitre.org/cve/

//CVE represente a CVE stored in db
type CVE struct {
	gorm.Model
	Reference   string
	Description string
}
