package main

import (
	"log"
	"netcat/netcat"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 2 {
		log.Fatal("[USAGE]: ./TCPChat $port")
	}

	name := "myServer"
	port := ":8989"

	if len(os.Args) > 1 {
		port = os.Args[1]
		if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}
	}

	netcat.Init(name, port)
}
