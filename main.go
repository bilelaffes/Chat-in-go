package main

import (
	"Chat/Client"
	"Chat/Server"
	"os"
)

func main() {

	argument := os.Args[1:]
	if argument[0] == "client" {
		Client.InitClient()
	}

	if argument[0] == "server" {
		Server.InitServer()
	}
}
