package config

import (
	"authServer/errors"

	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

func init() {
	errors.ErrorConfig("Error loading the env file...", godotenv.Load())
	configLogrus()

}

func configLogrus() {
	FILE := os.Getenv("LOGFILE")
	// Log to file
	if FILE != "" {

		// no colors for file output
		log.SetFormatter(&log.TextFormatter{
			ForceColors:   false,
			DisableColors: true,
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})

		f, err := os.OpenFile(FILE, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

		if err != nil {
			log.Fatalf("unable to create log file from config, %v", err)
		}

		log.SetOutput(f)
	} else {

		log.SetFormatter(&log.TextFormatter{
			ForceColors:   true,
			DisableColors: false,
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	}

	level, err := log.ParseLevel(os.Getenv("LOGLEVEL"))

	if err != nil {
		log.Fatalf("unable to set log level from config, %v", err)
	}
	log.SetLevel(level)

	// log File and line
	log.SetReportCaller(true)
}
