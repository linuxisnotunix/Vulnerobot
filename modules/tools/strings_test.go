package tools

import (
	"testing"
	"time"
)

func TestGetFrMonthNov(t *testing.T) {
	m := "novembre"
	mo := GetFrMonth(m)
	if mo != time.November {
		t.Error("Expected month (", m, "), to match : ", time.November, " !")
	}
}

func TestGetFrMonthUnk(t *testing.T) {
	m := "unknown"
	mo := GetFrMonth(m)
	if mo != 0 {
		t.Error("Expected month (", m, "), to match : ", 0, " !")
	}
}
