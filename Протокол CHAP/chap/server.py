import hashlib
import os
import socket
import string
import random
import time


# Хранилище паролей
USER_PASSWORDS = {
    'user1': 'pass',
    'user2': '123'
}

letters = string.ascii_letters + string.digits

def hashChap(id, passwd, value):
    return hashlib.md5((id + passwd + value).encode()).hexdigest()

def getValue():
        r = random.randint(1,16)
        value = "".join(random.choice(letters) for i in range(r))
        return value

def genId():
    id = os.urandom(1).hex()
    str.encode(id, encoding='utf-8')
    #print(id)
    return id

def checkUser(name):
    for i in USER_PASSWORDS:
        if i == name:
            return True
    return False        

def state1(username):
    identifier = genId()
    value = getValue()

    lenV = len(value)
    lenAll = 1 + len(identifier) + len(str(lenV)) + len(value) + len(username)
    ch1 = "1|" + identifier + "|" + str(lenAll) + "|" + str(lenV) + "|" + value + "|" + username

    return identifier, value, ch1

def state2(challenge, identifier, value):
    ch = challenge
    challenge = challenge.split("|")
    if len(ch) != int(challenge[2]):
        return 1
    match challenge[0]:
        case '2':
            if checkUser(challenge[5]) != True:
                resultAuth = "User not found"
                conn.send(resultAuth.encode())
                conn.close()
                return 0
            idUser = challenge[1]
            hashUser = challenge[4]
            hashDB = hashChap(identifier, USER_PASSWORDS[challenge[5]], value)
            if idUser == identifier:
                if hashDB == hashUser:
                    #print("success 3")
                    return 3
                return 4
        case default:
            return 2

def server_program():
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    server_socket.bind(('localhost', 8081))
    server_socket.listen(1)
    conn, address = server_socket.accept()
    print("Server Listing...")



    username = conn.recv(1024).decode()
    print(f"Username: {username}")
    if checkUser(username) != True:
        resultAuth = "User not found"
        conn.send(resultAuth.encode())
        print("Not found user")
        return 0

    while True:
        rTime = random.randint(1, 15)
        identifier, value, ch1 = state1(username)
        state = 2
        while state == 2:
            conn.send(ch1.encode())

            challenge = conn.recv(1024).decode()
            if state2(challenge, identifier, value) == 3:
                print(identifier, " -- Auth Success")
                resultAuth = "3|" + identifier + "|" + str(len(identifier)+2)
                conn.send(resultAuth.encode())
                print(str(rTime) +"s -- Next Check\n\n")
                time.sleep(rTime)
                state = 3

            elif state2(challenge, identifier, value) == 4:
                print(identifier, " -- Auth Fail. Wrong pass")
                resultAuth = "4|" + identifier + "|" + str(len(identifier)+2)
                conn.send(resultAuth.encode())
                conn.close()
                state = 4
                return 0

            elif state2(challenge, identifier, value) == 2:
                print(identifier, " -- Auth Error. Wrong code")

            elif state2(challenge, identifier, value) == 1:
                print(identifier, " -- Auth Error. Wrong length")




    



    conn.close()

if __name__ == '__main__':
    server_program()