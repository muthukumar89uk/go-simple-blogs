package drivers

import (
	//User-defined packages
	"blog/helper"
	"blog/logs"
	"blog/repository"

	//Inbuild packages
	"fmt"
	"os"

	//Third-party packages
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConnection() *gorm.DB {
	log := logs.Log()

	//Loading a '.env' file
	if err := helper.Config(".env"); err != nil {
		log.Error.Println("Error : 'Error at loading '.env' file'")
	}

	Host := os.Getenv("HOST")
	Port := os.Getenv("PORT")
	User := os.Getenv("USER")
	Password := os.Getenv("PASSWORD")
	Dbname := os.Getenv("DBNAME")

	//create a connection to ginframework database
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, Dbname)
	Db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Info.Printf("Message : Established a successful connection to '%s' database!!!\n", Dbname)

	//Table creation
	repository.TableCreation(Db)
	return Db
}
