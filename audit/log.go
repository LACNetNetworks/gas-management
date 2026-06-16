package audit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// GeneralLogger exported
var GeneralLogger *log.Logger

// ErrorLogger exported
var ErrorLogger *log.Logger

func init() {
	absPath, err := filepath.Abs("./log")
	if err != nil {
		fmt.Println("Error reading given path:", err)
	}

	_, err = os.Stat(absPath)
    if os.IsNotExist(err) {
    	errDir := os.MkdirAll(absPath, 0750)
        if errDir != nil {
        	panic(err)
        }
    }

	generalLog, err := os.OpenFile(absPath+"/idbServiceLog.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	GeneralLogger = log.New(generalLog, "General Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(generalLog, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
}