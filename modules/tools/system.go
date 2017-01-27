package tools

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"strings"
)

// ParseConfiguration parses a csv string and return a array of configurations
func ParseConfiguration(data string) *arraylist.List {
	list := arraylist.New()

	lines := strings.Split(data, "\n")

	// initialisation d'un élement d'une liste
	for i := range lines {
		fields := strings.Split(lines[i], ";")
		if len(fields) < 4 {
			continue
		}

		tab := make(map[string]string)
		// TrimSpace delete beginning and finishing spaces
		tab["Name"] = strings.TrimSpace(fields[0])
		tab["Version"] = strings.TrimSpace(fields[1])
		tab["Function"] = strings.TrimSpace(fields[2])
		tab["CPE"] = strings.TrimSpace(fields[3])
		list.Add(tab)
	}
	return list
}
