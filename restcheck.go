package main

import (
	"flag"
	"fmt"
	"github.com/Ericwyn/RestCheck/commend"
)

var versionMsg = "RestCheck, v0.1, @Ericwyn"

var initFlag = flag.String("init", "null", "init commend")

var saveFlag = flag.Bool("save", false, "save the api check result")
var msgFlag = flag.Bool("msg", false, "show the check result commend msg")

var initApiFlag = flag.String("initapi", "", "init api check msg")
var check = flag.String("check", "", "check one api")
var checkAll = flag.Bool("checkall", false, "check all api")

var env = flag.String("env", "", "load environment config for check")

var version = flag.Bool("v", false, "show version msg")

func main() {
	flag.Parse()

	if *version {
		fmt.Println(versionMsg)
	} else if *initFlag != "null" {
		commend.InitProject(*initFlag)
		return
	} else if *msgFlag {
		commend.ShowProjectMsg()
		return
	} else if *initApiFlag != "" {

		return
	} else {
		if *checkAll {
			commend.CheckAllApi(*env, *saveFlag)
			return
		}
		if *check != "" {
			commend.CheckApi(*env, *check, *saveFlag)
			return
		}
	}

	fmt.Println("use -h to list the commands")
}
