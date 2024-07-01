package main

import (
	"forum/Api"
	"forum/db"
)

func main() {

	db.InitDB()

	forum.ServerInit()
}
