package main

import (
	"log"
	"net"
	"os"
)

func parsArg() string {
	path := os.Args

	return path[1] + " " + path[2] + " " + path[3] + " " + path[4]
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	path := parsArg()

	conn.Write([]byte(path))
}
