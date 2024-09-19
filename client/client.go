package main

import (
	"fmt"
	"log"
	"net"
)

const (
	pkgInit = "initialization package"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	_, err = conn.Write([]byte(pkgInit))

	if err != nil {
		log.Fatal(err)
	}

	var resp [256]byte

	n, err := conn.Read(resp[:])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Порядковый номер: ", resp[n-1])
	fmt.Println("Зерно от сервера: ", resp[:n-1])

	fmt.Println(n)
}
