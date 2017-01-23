package database

import log "github.com/Sirupsen/logrus"

type dbLogHook struct {
}

func (h dbLogHook) Levels() []log.Level {
	return log.AllLevels
}
func (h dbLogHook) Fire(entry *log.Entry) error {
	log.Info("Custom DB Hoook fired! ", entry)
	//entry.Level = log.DebugLevel //Force debug level
	return nil
}
