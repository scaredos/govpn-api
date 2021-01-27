// GoVPN-API / version 2.0 (Production)
// Production (api.unknownvpn.net/)
package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"strings"
)

var users []string  // list of users
var tokens []string // list of tokens
var limit []string  // list of limit correlating to users/token position
var masterString string = "JoshuaIsGay"

func index(slice []string, item string) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}


func main() {
	userList, _ := readLines("users.txt")
	userMap := make(map[string]string)
	for i := range userList {
		username := strings.Split(userList[i], ":")[0]
		token := strings.Split(userList[i], ":")[1]
		userLimit := strings.Split(userList[i], ":")[2]
		users = append(users, username)
		tokens = append(tokens, token)
		limit = append(limit, userLimit)
	}
	for i := 0; i < len(users); i++ {
		userMap[users[i]] = tokens[i]
		//fmt.Printf("%s:%s\n", users[i], userMap[users[i]])
	}
	sport := "2052"
	fmt.Printf("running on localhost:%s\n", sport)

	// Handles REST-API Calls
	http.HandleFunc("/api/v1/userCount", userCount)
	http.HandleFunc("/api/v1/addUser", addUser)
	port := fmt.Sprintf(":%s", sport)
	http.ListenAndServe(port, nil)
}

// function: userCount
// description: Returns user management limit of the user
// on error, wrong username or token, returns 0 so client doesn't work
func userCount(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query()["token"]
	username := r.URL.Query()["username"]
	if len(token) == 0 || len(username) == 0 {
		fmt.Fprintf(w, "404\n")
		return
	} else if index(tokens, token[0]) != index(users, username[0]) {
		fmt.Fprintf(w, "0")
	} else if index(users, username[0]) == -1 {
		fmt.Fprintf(w, "0")
	} else {
		position := index(users, username[0])
		fmt.Fprintf(w, limit[position])
	}
}

// function: addUser
// description: Adds user to txt with username:token:limit
func addUser(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query()["key"]
	if len(key) == 0 || key[0] != masterString {
		fmt.Fprintf(w, "you do not have authorization to make users\n")
		return
	}
	username := r.URL.Query()["username"][0]
	token := r.URL.Query()["token"][0]
	mode := r.URL.Query()["mode"][0]
	if mode == "del" {
		position := index(users, username)
		if position == -1 {
			fmt.Fprintf(w, "doesn't exist")
			return
		}
		users[position] = users[len(users)-1]
		users[len(users)-1] = ""
		users = users[:len(users)-1]
		tokens[position] = tokens[len(tokens)-1]
		tokens[len(tokens)-1] = ""
		tokens = tokens[:len(tokens)-1]
		fileBytes, _ := ioutil.ReadFile("users.txt")
		lines := strings.Split(string(fileBytes), "\n")
		var newLines[]string
		for _, line := range lines {
			if strings.Contains(line, username) {
				continue
			} else {
				newLines = append(newLines, line)
			}
		}
		err := os.Remove("users.txt")
		if err != nil {fmt.Println(err)}
		file, _ := os.Create("users.txt")
		defer file.Close()
		for _, liner := range newLines {
			if !strings.Contains(liner, ":") {
				continue
			} else {
				_, err := file.WriteString(fmt.Sprintf("%s\n", liner))
				if err != nil {fmt.Println(err)}
			}	
		}
		fmt.Fprintf(w, "removed user")
	} else if mode == "add" {
		if index(users, username) != -1 {
			fmt.Fprintf(w, "user already exists\n")
			return
		}
		if index(tokens, token) != -1 {
		fmt.Fprintf(w, "token already exits\n")
		return
		}
		limits := r.URL.Query()["limit"][0]
		file, _ := os.OpenFile("users.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		_, err := file.WriteString(fmt.Sprintf("%s:%s:%s\n", username, token, limits))
		defer file.Close()
		if err != nil {fmt.Println(err)}
		fmt.Fprintf(w, fmt.Sprintf("success (%s:%s:%s)\n", username, token, limits))
		users = append(users, username)
		tokens = append(tokens, token)
		limit = append(limit, string(limits))
		return
	} else {
		fmt.Fprintf(w, "404")
		return
	}
}
