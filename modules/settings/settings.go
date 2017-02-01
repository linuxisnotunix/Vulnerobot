package settings

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// settings
var (
	AppName               string
	AppVersion            string
	AppBranch             string
	AppCommit             string
	AppBuildTime          string
	AppVerbose            bool
	AppPath               string
	DBPath                string
	ConfigPath            string
	WebListen             string
	UIDontDisplayProgress bool
	PluginList            string
)

func init() {
	var err error
	if AppPath, err = execPath(); err != nil {
		log.Fatalf("Fail to get app path: %v\n", err)
	}
	// Note: we don't use path.Dir here because it does not handle case
	//	which path starts with two "/" in Windows: "//psf/Home/..."
	AppPath = strings.Replace(AppPath, "\\", "/", -1)
}

// execPath returns the executable path.
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

// WorkDir returns absolute path of work directory.
func WorkDir() string {
	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath
	}
	return AppPath[:i]
}
