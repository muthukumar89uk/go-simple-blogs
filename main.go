package main

import (
	//User-defined package
	"blog/drivers"
	"blog/router"
)

func main() {
	//Establishing a DB-connection
	Db := drivers.DbConnection()

	//Routing all the handlers
	router.Router(Db)
}
