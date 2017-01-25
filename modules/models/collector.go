package models

import "github.com/gosuri/uiprogress"

//Collector interface for collector type
type Collector interface {
	ID() string
	IsAvailable() bool
	Collect(bar *uiprogress.Bar) error
	List() error
}
