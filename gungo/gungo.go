package gungo

import (
	"fmt"
	"log"
	"time"
)

var (
	//GunDebug Set to True to get Debug Logs
	GunDebug bool
)

func gunLog(msgFormat string, msg ...interface{}) string {
	var logMsg = ""
	if GunDebug {
		logMsg = fmt.Sprintf(msgFormat, msg...)
		fmt.Println(logMsg)
	}
	return logMsg
}
func gunErr(msgFormat string, msg ...interface{}) string {
	var logMsg = ""
	if GunDebug {
		logMsg = fmt.Sprintf(msgFormat, msg...)
		log.Println(logMsg)
	}
	return logMsg
}
func gunTimed(msgFormat string, startTime time.Time, msg ...interface{}) string {
	var logMsg = ""
	if GunDebug {
		msg = append(msg, time.Since(startTime))
		logMsg = fmt.Sprintf(msgFormat+"\t%s", msg...)
		log.Println(logMsg)
	}
	return logMsg
}

func init() {
	gunLog("init")
}
