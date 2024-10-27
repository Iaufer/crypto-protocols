import hashlib
import socket
import string
import random


def hashChap(id, passwd, value):
    return hashlib.md5((id + passwd + value).encode()).hexdigest()

def lenCh(ch):
    lenAll = len(ch)
    lenAll = lenAll + len(str(lenAll)) - 1
    ch1 = ch.split("|")
    ch1[2] = lenAll
    ch2 = "|".join([str(item) for item in ch1])
    return ch2


def client_program():
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect(('localhost', 8081))



    username = input("Enter username: ")
    client_socket.send(username.encode())

    while True:
        challenge = client_socket.recv(1024).decode()
        if challenge == "User not found":
            print("User not found")
            client_socket.close()
            return 0
        
        challenge = challenge.split("|")
        match challenge[0]:
            case '1':
                print("challange = ", challenge)
                password = input("Enter password: ")
                identifier = challenge[1]
                val = challenge[4]
                username = challenge[5]
                hashUser = hashChap(identifier, password, val)
                # print(hash)

                lenHash = len(hashUser)
                ch2 = "2|" + identifier + "|" + " " + "|" + str(lenHash) + "|" + hashUser + "|" + username
                ch2 =lenCh(ch2)
                if random.choice([True, False]):
                    ch2 = "9"+ch2[1:]
                    print("Wrong Code", ch2, "\n")

                if random.choice([True, False]):
                    ch2 = ch2[:-3]
                    print("Wrong length", ch2, "\n")
                client_socket.send(ch2.encode())
            case '3':
                print("Auth Success\n\n")

            case "4":
                print("Auth Error. Wrong password\n\n")
                client_socket.close()
                return 0


    client_socket.close()

if __name__ == '__main__':
    client_program()