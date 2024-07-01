package main

import (
	api "forum/Api"
	backend "forum/db"
)

func main() {

	db := backend.OpenConnection()

	server := api.New(db)
	server.Init()
}
