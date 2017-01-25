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
