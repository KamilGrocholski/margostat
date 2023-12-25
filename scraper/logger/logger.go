package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

func Init() {
	file, err := os.OpenFile("scraper.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}

	Logger = log.New(file, "SCRAPER: ", log.Ldate|log.Ltime|log.Lshortfile)
}
