package logger

import (
	"os"
	"strings"
	"time"
)

type LogFile struct {
	flog *os.File
	path string
	date string
}

func LogOpen(path string) (log *LogFile, err error) {
	log = new(LogFile)
	log.date = time.Now().Format(time.RFC3339)[0:10]
	log.path = path
	path += "log" + log.date + ".log"

	log.flog, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return
}
func (l *LogFile) Read(p []byte) (n int, err error) {
	n, err = l.flog.Read(p)
	return
}

func (l *LogFile) Write(p []byte) (n int, err error) {
	date := time.Now().Format(time.RFC3339)[0:10]
	n = 0
	if strings.Compare(l.date, date) != 0 {
		l.flog.Close()
		l.date = date
		path := l.path + "log" + l.date + ".log"
		l.flog, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return
		}
	}
	n, err = l.flog.Write(p)
	return
}

//Close закрытие
func (l *LogFile) Close() error {
	err := l.flog.Close()
	return err
}

// func (l *LogFile) KillOld(old time.Time) error {
// 	err := filepath.Walk(l.path, func(path string, info os.FileInfo, err error) error {
// 		if strings.Contains(path, ".log") {
// 			year, _ := strconv.ParseInt(path[3:6], 10, 32)
// 			month, _ := strconv.ParseInt(path[8:9], 10, 32)
// 			day, _ := strconv.ParseInt(path[10:11], 10, 32)
// 			date := time.Date(year, month, day, 0, 0, 0, 0, nil)
// 		}
// 		return err
// 	})
// 	return nil
// }
