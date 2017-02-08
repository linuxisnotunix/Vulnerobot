package tools

import "testing"

func TestParseConfigurationEmpty(t *testing.T) {
	config := ""
	list := ParseConfiguration(config)
	if list.Size() != 0 {
		t.Error("Expected config list to be empty !")
	}
}

func TestParseConfigurationSample(t *testing.T) {
	config := `Ubuntu ; 12.04 ; FP1,FP3 ; cpe:/o:canonical:ubuntu_linux:12.04
	Ubuntu ; 16.10 ; FP1,FP2 ; cpe:/o:canonical:ubuntu_linux:16.10
	Microsoft Windows ; 7 ; FP2,FP3 ; cpe:/o:microsoft:windows_7
	Microsoft Windows ; 8.1 ; FP2 ; cpe:/o:microsoft:windows_8.1
	Mozilla Firefox ; 48 ; FP2 ; cpe:/a:mozilla:firefox:48
	Adobe Acrobat DC ; 15.006.30173 ; FP1,FP2 ; cpe:/a:adobe:acrobat_dc:15.006.30173
	Adobe Reader ; 11.0.04 ; FP1,FP3 ; cpe:/a:adobe:acrobat_reader:11.0.04`

	list := ParseConfiguration(config)
	if list.Size() != 7 {
		t.Error("Expected config list to be 7 rows !")
	}

	l1, exist := list.Get(0)
	if !exist {
		t.Error("Expected first line not found !")
	}
	if l1.(map[string]string)["Name"] != "Ubuntu" {
		t.Error("Expected first line name to be Ubuntu !")
	}
}

func TestParseFlagListDefault(t *testing.T) {
	v := ParseFlagList("")
	if v != nil {
		t.Error("Expected default parsed flag result to be null !")
	}
}

func TestParseFlagListAll(t *testing.T) {
	v := ParseFlagList("")
	if v != nil {
		t.Error("Expected 'all' parsed flag result to be null !")
	}
}

func TestParseFlagListNVD(t *testing.T) {
	v := ParseFlagList("nvd")
	if v == nil {
		t.Error("Expected 'nvd' parsed flag result to not be null !")
	}
	if v.Size() != 1 {
		t.Error("Expected 'nvd' parsed flag result size 1 !")
	}
}

func TestParseFlagListMultilple(t *testing.T) {
	v := ParseFlagList("nvd,anssi")
	if v == nil {
		t.Error("Expected 'nvd,anssi' parsed flag result to not be null !")
	}
	if v.Size() != 2 {
		t.Error("Expected 'nvd,anssi' parsed flag result size 2 !")
	}
}
