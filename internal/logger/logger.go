package logger

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	logFileName = "scdebug.log"
)

func SetupLogger() (bool, *log.Logger) {

	var (
		debugset bool = false
		lf       *os.File
		err      error
	)

	if len(os.Getenv("DEBUG")) > 0 {
		lf, err = tea.LogToFile(logFileName, "debug")
		if err != nil {
			log.Fatalf("FATAL: %s\n", err)
		}
		debugset = true
	} else {
		lf, err = os.OpenFile("/dev/null", os.O_WRONLY|os.O_APPEND, 0000)
		if err != nil {
			log.Fatalf("Open /dev/null for logging ended up fatal!\n")
		}
	}

	l := log.New(lf, "SC: ", log.Lshortfile|log.Lmicroseconds)
	l.Printf("Log file: %s\n", lf.Name())

	return debugset, l
}
