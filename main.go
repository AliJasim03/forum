package main

import (
    "forum/db"
    "forum/server"
)

func main() {

    db.InitDB()

    server.ServerInit()
}

