package database

import (
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //Access to sqlite driver
	"github.com/linuxisnotunix/Vulnerobot/modules/models"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
)

//Orm singleton database access.
var Orm *gorm.DB

func setLogger() {
	logger := log.New()
	logger.Level = log.DebugLevel
	logger.Out = os.Stdout
	//logger.SetFormatter(&log.TextFormatter{})
	//Orm.SetLogger(logger)
	Orm.SetLogger(logger)
}
func init() {
	//setLogger()
	var err error
	_, err = os.Stat(settings.DBPath)
	if err != nil {
		//Db don't exist
		if Orm, err = gorm.Open("sqlite3", settings.DBPath); err != nil {
			setLogger()
			log.Debug("Fail to access DB file(%s) but this is ok since the file don't exist: %v\n", settings.DBPath, err)
			Orm.CreateTable(models.CVE{}) //TODO testing
		}
	} else {
		//DB exist
		if Orm, err = gorm.Open("sqlite3", settings.DBPath); err != nil {
			setLogger()
			log.Warnf("Fail to access DB file(%s): %v\n", settings.DBPath, err)
		}
	}
	//Start debug mode
	if settings.AppVerbose {
		Orm.Debug()
	}
	//Orm.LogMode(settings.AppVerbose)

	defer Orm.Close()
}
