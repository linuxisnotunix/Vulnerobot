package models

//Component represent a Component and is version
type Component struct {
	ID           uint32 `gorm:"primary_key"`
	Name         string
	Version      string
	CPE          string
	IsAffectedBy []Function `gorm:"many2many:anssi_affect_components;"`
	InFunctions  []Function `gorm:"many2many:components_functions;"`
}

//Function represente a function
type Function struct {
	Name         string      `gorm:"primary_key"`
	UseComponent []Component `gorm:"many2many:components_functions;"`
}

//ResponseListEntry entry match in list response
type ResponseListEntry struct {
	Component map[string]string //TODO replace by struct
	Matchs    []ResponseListMatch
	Vulns     []ResponseListVuln
}

//ResponseListMatch  what entry match in list reponse
type ResponseListMatch struct {
	Type  string
	Value string
}

//ResponseListVuln what Vuln match in list reponse
type ResponseListVuln struct {
	Source string
	Value  map[string]string
}

//ResponseStatus represent the status of the app
type ResponseStatus struct {
	Plugins    []string
	Components []map[string]string
}
