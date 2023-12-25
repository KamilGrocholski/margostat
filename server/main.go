package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger

func main() {
	initLogger()

	http.HandleFunc("/", handleFile)
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	staticDirPath, ok := os.LookupEnv("STATIC_DIR_PATH")
	if ok == false {
		panic("no STATIC_DIR_PATH env")
	}
	world := r.URL.Path[len("/"):]
	if world == "" {
		world = "index"
	}

	logger.Printf("request for file: %s/%s.html", staticDirPath, world)

	filePath := fmt.Sprintf("%s/%s.html", staticDirPath, world)
	http.ServeFile(w, r, filePath)
}

func initLogger() {
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}

	logger = log.New(file, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)
}
