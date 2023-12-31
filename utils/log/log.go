package log

import (
	"fmt"
	"github.com/Ericwyn/GoTools/date"
	"time"
)

type logTag string

var timeFormat = "[MMdd-HH:mm:ss]"

var info logTag = "INFO"
var debug logTag = "DBUG"
var err logTag = "ERR "

func I(msg ...interface{}) {
	printLog(info, msg...)
}

func D(msg ...interface{}) {
	printLog(debug, msg...)
}

func E(msg ...interface{}) {
	printLog(err, msg...)
}

func printLog(tag logTag, msg ...interface{}) {
	fmt.Println(date.Format(time.Now(), timeFormat), "[RCHECK-"+tag+"]", fmt.Sprint(msg...))
}
