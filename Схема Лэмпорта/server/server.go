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
		hash := h.Sum(nil)

		xorResult := make([]byte, 8)
		for j := 0; j < 4; j++ {
			xorResult[j] = hash[j] ^ hash[j+8]     // Первая часть ^ Третья часть
			xorResult[j+4] = hash[j+4] ^ hash[j+12] // Вторая часть ^ Четвертая часть
		}
		result = hex.EncodeToString(xorResult)
		str = result
	}
	return result //64 бит
}

func addUser(user []string) error {
	file, err := os.OpenFile("data.csv", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	seed := generateSeed()

	i, _ := strconv.ParseInt(user[2], 10, 64)

	_, err = file.Write([]byte(user[0] + " " + hashPS(user[1], seed, i) + " " + decIter(user[2]) + " " + seed + "\n"))
	if err != nil {
		return err
	}
	fmt.Println(user[0]," -- Пользователь добавлен")

	return nil
}

func updateHashCount(username, hash string, conn net.Conn) error {
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
	if num == 1 {
		conn.Write([]byte("Вход выполнен успешно! N.B: Чтобы войти в следующий раз - зарегистрируйтесь"))
		delete(maps, username)
	} else {
		value[1] = strconv.Itoa(num - 1)
	}

	file, err = os.Create("data.csv")
	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(file)

	for key, value := range maps {
		_, err := writer.WriteString(key + " " + value[0] + " " + value[1] + " " + value[2] + "\n")
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
	reset := false
	defer conn.Close()

	buffer := make([]byte, 100)
	n, _ := conn.Read(buffer)

	s := string(buffer[:n])
	name := strings.Split(string(s), " ")[1]
	switch s[0] {
	case '0': // reg
		err1 := checkUser(name)

		if err1 == nil {
			fmt.Println("unreg" + " " + "Пользователь с таким именем уже зарегистрирован!!!")
		} else {
			err := addUser(strings.Split(s[2:], " "))
			if err != nil {
				log.Fatal(err)
			}
		}
	case '1': // reset
		reset = true
		fallthrough
	case '2': // auth
		hash, count, seed, err := getInfo(s)
		if err != nil {
			// log.Fatal(err)
		}

		err = checkUser(name)

		if err != nil {
			conn.Write([]byte("unreg" + " " + "Пользователь с таким именем не зарегистрирован!!!"))

		} else {

			conn.Write([]byte(count + " " + seed))
		}

		buffer := make([]byte, 20)
		n, _ := conn.Read(buffer)
		newHash := hashPS(string(buffer[:n]), seed, 1)
		if hash == newHash {

			err = updateHashCount(name, string(buffer[:n]), conn)

			if err != nil {
				log.Fatal(err)
				return
			}
			conn.Write([]byte("Auth Success"))
			fmt.Println(name, " -- Пользователь успешно авторизован")

			if reset {
				deleteUser(name)
				err := addUser(strings.Split(s[2:], " "))
				fmt.Println(name, " -- Пользователь успешно сбросил пароль")
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			conn.Write([]byte("Auth Unsuccess"))
			fmt.Println(name, " -- Неуспешная попытка аутентификации")
		}
	default:
		fmt.Println("Errors something")
	}
}

func deleteUser(name string) error {
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

	_, exist := maps[name]

	if !exist {
		return errors.New("El not found")
	}

	delete(maps, name)

	file, err = os.Create("data.csv")
	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(file)

	for key, value := range maps {
		_, err := writer.WriteString(key + " " + value[0] + " " + value[1] + " " + value[2] + "\n")
		if err != nil {
			return err
		}

	}
	fmt.Println(name, " -- Пользователь успешно удален")
	return writer.Flush()
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
