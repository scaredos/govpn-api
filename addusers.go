// Private file
// GoVPN-API / version 2.0 (Production)
// Make users (addusers.exe)
package main

import (
  "fmt"
  "strings"
  "io/ioutil"
  "net/http"
  "os"
)
var username string
var token string
var limit string
var mode string

func banner() {
  fmt.Println("addusers\t|\tgovpn-api admin tool")
  fmt.Println("-u, --user\t|\tUsername")
  fmt.Println("-t, --token\t|\tToken")
  fmt.Println("-l, --limit\t|\tLimit")
  os.Exit(0)
}

func main() {
  limit = "0"
  if len(os.Args) <= 1 {
    banner()
  }
  for i, ch := range os.Args {
    if i == 0 {continue}
    if strings.Contains(ch, "-u") || strings.Contains(ch, "--user") {
      username = os.Args[i+1]
      } else if strings.Contains(ch, "-t") || strings.Contains(ch, "--token") {
      token = os.Args[i+1]
    } else if strings.Contains(ch, "-l") || strings.Contains(ch, "--limit") {
      limit = os.Args[i+1]
    } else if strings.Contains(ch, "-d") {
      mode = "del"
    } else {mode = "add"}
  }
  client := &http.Client{}
  url := fmt.Sprintf("http://govpnapi.unknownvpn.net:2052/api/v1/addUser?username=%s&token=%s&key=JoshuaIsGay&limit=%s&mode=%s", username, token, limit, mode)
  req, _ := http.NewRequest("GET", url, nil)
  resp, err := client.Do(req)
  if err != nil {fmt.Println(err)}
  body, _ := ioutil.ReadAll(resp.Body)
  response := string(body)
  fmt.Println(response)
}
