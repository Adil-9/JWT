package errorlog

import "log"

func ErrorPrint(str string, err error) {
	log.Printf("%s: %s", str, err.Error())
}

func ErrorPrintFatal(str string, err error) {
	log.Fatalf("%s: %s", str, err.Error())
}
