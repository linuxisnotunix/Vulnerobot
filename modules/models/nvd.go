package models

import "time"

//NvdList represente a NvdCVE stored in db
type NvdList struct {
	Year             string `gorm:"primary_key"`
	LastModifiedDate string
	Sha256           string
}

//NvdCVE represente a NvdCVE stored in db
type NvdCVE struct {
	ID           string `gorm:"primary_key"`
	Summary      string
	AffectedCPEs []NvdCPE `gorm:"many2many:nvdcve_affect_nvdcpe;"`
	//Affect       []Component `gorm:"many2many:cve_affect_components;"`
	Release    time.Time
	LastUpdate time.Time
}

//NvdCPE represente a NvdCPE stored in db
type NvdCPE struct {
	RelatedCVEs []NvdCVE `gorm:"many2many:nvdcve_affect_nvdcpe;"`
	ID          string   `gorm:"primary_key"`
}
