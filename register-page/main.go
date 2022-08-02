package main

import (
	"web-server/server"
)

func main() {
	server.DebugPrint = true

	server.Start()
}
