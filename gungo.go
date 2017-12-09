package main

import (
	"fmt"
	"log"
	"time"
)

var (
	gunDebug bool
)

func gunLog(msgFormat string, msg ...interface{}) {
	if gunDebug {
		fmt.Printf(msgFormat+"\n", msg...)
	}
}
func gunErr(msgFormat string, msg ...interface{}) {
	log.Printf(msgFormat, msg...)
}
func gunTimed(msgFormat string, startTime time.Time, msg ...interface{}) {
	if gunDebug {
		msg = append(msg, time.Since(startTime))
		log.Printf(msgFormat+"\t%s", msg...)
	}
}
