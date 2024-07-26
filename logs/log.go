package logs

import (
	//user defined package
	"blog/models"

	//inbuild package(s)
	"log"
	"os"
)

// Create a custom log
func Log() (logger models.Logs) {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	file, _ := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logger.Info = log.New(file, "[INFO:] ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	logger.Error = log.New(file, "[ERROR:]", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	return
}
