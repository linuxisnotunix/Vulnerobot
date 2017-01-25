package models

import "time"

//AnssiAVI represente a AnssiAVI stored in db
type AnssiAVI struct {
	ID             string `gorm:"primary_key"`
	Title          string
	Risk           string
	SystemAffected string
	Affect         []Component `gorm:"many2many:anssi_affect_components;"`
	Release        time.Time
	LastUpdate     time.Time
}
