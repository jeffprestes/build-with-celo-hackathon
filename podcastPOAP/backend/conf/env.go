package conf

import (
	"flag"
	"log"

	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"
)

// Cfg represents the pointer to configuration file
var Cfg *ini.File

// find configuration file
func init() {
	var err error
	var appIniLocation string
	flag.StringVar(&appIniLocation, "appIni", "conf/app.ini", "defines where the app initialization config location (app.ini) is")
	flag.Parse()
	Cfg, err = macaron.SetConfig(appIniLocation)
	if err != nil {
		if isDbConnParamsInEnvVariables() {
			log.Printf("[conf/Init] Error during app.ini reading. Error: %s\n", err.Error())
		} else {
			log.Fatalf("[conf/Init] Error during app.ini reading. Error: %s\n", err.Error())
		}
	}
}
