package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
)

func generateSeed() ([]byte, error) {
	seed := make([]byte, 8) // сколько байт зерно?

	_, err := rand.Read(seed)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return seed, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var buffer [256]byte // есть ли смысл заранее, чтоыб клиент заранее передавал размер пакета инициализации?
	n, err := conn.Read(buffer[:])

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Пакет инициализации от клиента: ", buffer[:n])

	seed, err := generateSeed()

	if err != nil {
		log.Fatal(err)
	}

	// resp := fmt.Sprintf("1:%x", seed) // порядковый номер уникальный для клиента? или от чего он зависит
	// как пережовать порядковый номер?
	seed = append(seed, 5) // 5 порядкоый номер который бдует в конце массива
	fmt.Println(seed)
	conn.Write([]byte(seed))

}

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Println("Server is listening...")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)

		// buffer := make([]byte, 1024) // есть ли смысл клиенту передавать соклько байт в пакете инициализции?
		// conn.Read(buffer)
		// fmt.Println(buffer)
		// conn.Close()
	}
}
