package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/m1/go-generate-password/generator"
)

func generateSeed() string {
	config := generator.Config{
		Length:                     16,
		IncludeSymbols:             false,
		IncludeNumbers:             true,
		IncludeLowercaseLetters:    true,
		IncludeUppercaseLetters:    false,
		ExcludeSimilarCharacters:   false,
		ExcludeAmbiguousCharacters: false,
	}

	g, _ := generator.New(&config)

	seed, _ := g.Generate()

	return *seed
}

func checkUser(username string) error {
	file, err := os.OpenFile("data.csv", os.O_RDONLY, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), " ")
		fmt.Println(row[0], username, len(row[0]), len(username))
		if row[0] == username {
			return nil
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return errors.New("User not found")
}

func decIter(num string) string {
	n, err := strconv.ParseInt(num, 10, 64)

	if err != nil {
		log.Fatal(err)
	}
	return strconv.Itoa(int(n) - 1)
}

func hashPS(str string, i int64) string {
	h := md5.New()
	for ; i > 0; i-- {
		h.Write([]byte(str))
	}
	fmt.Println("i = ", i)
	return hex.EncodeToString(h.Sum(nil)[:8])
}

func addUser(user []string) error {
	file, err := os.OpenFile("data.csv", os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	seed := generateSeed()

	i, _ := strconv.ParseInt(user[3], 10, 64)

	file.Write([]byte(user[1] + " " + hashPS(user[2]+seed, i) + " " + decIter(user[3]) + " " + seed + "\n"))

	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 100)
	n, _ := conn.Read(buffer)

	user := strings.Split(string(buffer[:n]), " ")

	err := checkUser(user[1])

	if err != nil {
		fmt.Println("User not found")
		err := addUser(user)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("User found")
	}

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
	}

}
