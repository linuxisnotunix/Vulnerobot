package models

//Collector interface for collector type
type Collector interface {
	ID() string
	IsAvailable() bool
	Collect() error
	List() error
}
