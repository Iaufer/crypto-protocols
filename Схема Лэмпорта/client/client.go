package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)


func hashPS(str, seed string, i int64) string {
	var result string

	for ; i > 0; i-- {
		h := md5.New()
		h.Write([]byte(str + seed))
		hash := h.Sum(nil)

		xorResult := make([]byte, 8)
		for j := 0; j < 4; j++ {
			xorResult[j] = hash[j] ^ hash[j+8]     // Первая часть ^ Третья часть
			xorResult[j+4] = hash[j+4] ^ hash[j+12] // Вторая часть ^ Четвертая часть
		}
		result = hex.EncodeToString(xorResult)
		str = result
	}
	fmt.Println("i = ", i)
	return result //64 бит
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not enough arg...")
		return
	}

	command := args[1]

	switch command {
	case "keyinit":
		if len(args) == 3 { //auth
			conn.Write([]byte("2 " + args[2]))

			buff := make([]byte, 256)
			n, _ := conn.Read(buff)

			if strings.Split(string(buff[:n]), " ")[0] == "unreg" {
				fmt.Println(string(buff[6:n]))
				return
			}

			passwd := ""
			fmt.Print("Enter password: ")
			fmt.Scanln(&passwd)

			num, _ := strconv.Atoi(strings.Split(string(buff), " ")[0])
			h := hashPS(passwd, strings.Split(string(buff[:n]), " ")[1], int64(num))
			// fmt.Println(num, h)
			// fmt.Println(h)

			conn.Write([]byte(h))
			n, _ = conn.Read(buff[:])
			fmt.Println(string(buff[:n]))
		}

		if len(args) == 5 {
			conn.Write([]byte("0 " + args[2] + " " + args[3] + " " + args[4]))
		}

		if len(args) == 6 {
			conn.Write([]byte("1 " + args[3] + " " + args[4] + " " + args[5]))

			buff := make([]byte, 256)
			n, _ := conn.Read(buff)

			if strings.Split(string(buff[:n]), " ")[0] == "unreg" {
				fmt.Println(string(buff[6:n]))
				return
			}

			passwd := ""
			fmt.Print("Enter password: ")
			fmt.Scanln(&passwd)

			num, _ := strconv.Atoi(strings.Split(string(buff), " ")[0])
			h := hashPS(passwd, strings.Split(string(buff[:n]), " ")[1], int64(num))

			conn.Write([]byte(h))
			n, _ = conn.Read(buff[:])
			fmt.Println(string(buff[:n]))
		}

	default:
		fmt.Println("Errors arg...")
	}
}
