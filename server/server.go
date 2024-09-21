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

func hashPS(str, seed string, i int64) string {
	var result string

	for ; i > 0; i-- {
		h := md5.New()
		h.Write([]byte(str + seed))
		hash := h.Sum(nil)[:8]
		result = hex.EncodeToString(hash)
		str = result
	}
	fmt.Println("i = ", i)
	return result //64 бит
}

func addUser(user []string) error {
	file, err := os.OpenFile("data.csv", os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	seed := generateSeed()

	i, _ := strconv.ParseInt(user[2], 10, 64)

	file.Write([]byte(user[0] + " " + hashPS(user[1], seed, i) + " " + decIter(user[2]) + " " + seed + "\n"))

	return nil
}

func updateHashCount(username, hash string) error {
	file, err := os.Open("data.csv")
	if err != nil {
		return err
	}

	defer file.Close()

	maps := make(map[string][]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), " ")
		maps[row[0]] = row[1:]
	}

	value, exist := maps[username]

	if !exist {
		return errors.New("El not found")
	}

	value[0] = hash
	num, _ := strconv.Atoi(value[1])
	if num <= 0 {
		fmt.Println("Вам нужно пройти регистрацию еще раз")
		delete(maps, username)
	} else {
		value[1] = strconv.Itoa(num - 1)
	}

	// for key, value := range maps {
	// 	fmt.Println(key, value, len(value))
	// }

	file, err = os.Create("data.csv")
	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(file)

	for key, value := range maps {
		_, err := writer.WriteString(key + " " + value[0] + " " + value[1] + " " + value[2])
		if err != nil {
			return err
		}

	}

	return writer.Flush()
}

func getInfo(name string) (string, string, string, error) {
	file, err := os.OpenFile("data.csv", os.O_RDONLY, 0666)
	if err != nil {
		return "", "", "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), " ")
		if row[0] == strings.Split(name, " ")[1] {
			return row[1], row[2], row[3], nil
		}
	}

	return "", "", "", errors.New("User not found")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 100)
	n, _ := conn.Read(buffer)
	s := string(buffer[:n])
	fmt.Println(s[0])
	if s[0] == '0' {
	}
	switch s[0] {
	case '0': // reg
		fmt.Println("kekekek")

		err := addUser(strings.Split(s[2:], " "))
		if err != nil {
			log.Fatal(err)
		}
	case '1': // reset

	case '2': // auth
		hash, count, seed, err := getInfo(s)
		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte(count + " " + seed))

		buffer := make([]byte, 20)
		n, _ := conn.Read(buffer)

		newHash := hashPS(string(buffer[:n]), seed, 1)
		fmt.Println(hash, newHash)
		// fmt.Println(string(buffer), n)
		if hash == newHash {
			err = updateHashCount("artyom", string(buffer[:n]))
			fmt.Println(string(buffer))

			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Println("Auth Success")
		} else {
			fmt.Println("Auth unSuccess")
		}

		// fmt.Println(strings.Split(s[2:], " "))

	default:
		fmt.Println("Errors something")
	}

	// user := strings.Split(string(buffer[:n]), " ")

	// err := checkUser(user[1])

	// if err != nil {
	// 	fmt.Println("User not found")
	// 	err := addUser(user)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	fmt.Println("User found")
	// }

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
