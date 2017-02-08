package tools

import "testing"

func TestIsTCPAccessibleOK(t *testing.T) {
	endpoint := "google.com:https"
	ok, err := IsTCPAccessible(endpoint)
	if !ok || err != nil {
		t.Error("Expected endpoint (", endpoint, "), to be accessible !")
	}
}

func TestIsTCPAccessibleNOK(t *testing.T) {
	endpoint := "some.fake.host:https"
	ok, _ := IsTCPAccessible(endpoint)
	if ok {
		t.Error("Expected endpoint (", endpoint, "), to not be accessible !")
	}
}
