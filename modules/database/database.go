package database

import (
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //Access to sqlite driver
	"github.com/linuxisnotunix/Vulnerobot/modules/models"
	"github.com/linuxisnotunix/Vulnerobot/modules/settings"
)

//orm singleton database access.
var orm *gorm.DB
var logger *log.Logger
var once sync.Once

//setLogger Update config of db logger
func setLogger() {
	logger = log.New()
	logger.Level = log.GetLevel()
	logger.Hooks.Add(dbLogHook{})
	//logger.Out = os.Stdout
	//logger.Formatter = new(log.TextFormatter)
	//log.Debugf("Logger: %v\n", logger)
	orm.SetLogger(logger)
	orm.LogMode(settings.AppVerbose) //Only log if debug is active

	//Start debug mode
	if settings.AppVerbose {
		orm.Debug()
	}
}

//Orm Give access to insstancied orm
func Orm() *gorm.DB {
	once.Do(func() {
		setup()
	})
	return orm
}

//Setup construct database
func setup() {
	var err error
	_, err = os.Stat(settings.DBPath)
	if err != nil {
		//Db don't exist
		orm, err = gorm.Open("sqlite3", settings.DBPath)
		setLogger()
		if err != nil {
			log.Debugf("Fail to access DB file(%s) but this is ok since the file don't exist: %v\n", settings.DBPath, err)
		}
		//orm.CreateTable(models.CVE{})      //TODO testing
		orm.CreateTable(models.AnssiAVI{}) //TODO testing

	} else {
		//DB exist
		orm, err = gorm.Open("sqlite3", settings.DBPath)
		setLogger()
		if err != nil {
			log.Fatalf("Fail to access DB file(%s): %v\n", settings.DBPath, err)
		}
		log.Debugf("Db init ok: %v\n", settings.DBPath, err)
	}
	//defer orm.Close()
}
