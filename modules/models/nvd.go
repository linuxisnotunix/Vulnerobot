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
	AffectedCPEs string
	Affect       []Component `gorm:"many2many:cve_affect_components;"`
	Release      time.Time
	LastUpdate   time.Time
}
