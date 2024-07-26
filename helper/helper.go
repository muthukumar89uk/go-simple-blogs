package helper

import (
	//User-defined package(s)
	"blog/logs"

	//Third-party package(s)
	"github.com/joho/godotenv"
)

const (
	SecretKey = "secret"
	Host      = "localhost"
	Port      = 5435
	User      = "postgres"
	Password  = "password"
	Dbname    = "mitrah_blog"
)

func Config(file string) error {
	log := logs.Log()
	//Load the given file
	if err := godotenv.Load(file); err != nil {
		log.Error.Printf("Error : 'Error at loading %s file'", file)
		return err
	}
	return nil
}
