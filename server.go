// GoVPN-API | Release
package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// ApiData struct for using returned JSON from JSON-RPC API
type ApiData struct {
	Result struct {
		AuthTypeU32           int    `json:"AuthType_u32"`
		CreatedTimeDt         string `json:"CreatedTime_dt"`
		ExpireTimeDt          string `json:"ExpireTime_dt"`
		GroupNameStr          string `json:"GroupName_str"`
		HashedKeyBin          string `json:"HashedKey_bin"`
		HubNameStr            string `json:"HubName_str"`
		NameStr               string `json:"Name_str"`
		NoteUtf               string `json:"Note_utf"`
		NtLmSecureHashBin     string `json:"NtLmSecureHash_bin"`
		NumLoginU32           int    `json:"NumLogin_u32"`
		RealnameUtf           string `json:"Realname_utf"`
		RecvBroadcastBytesU64 int    `json:"Recv.BroadcastBytes_u64"`
		RecvBroadcastCountU64 int    `json:"Recv.BroadcastCount_u64"`
		RecvUnicastBytesU64   int    `json:"Recv.UnicastBytes_u64"`
		RecvUnicastCountU64   int    `json:"Recv.UnicastCount_u64"`
		SendBroadcastBytesU64 int    `json:"Send.BroadcastBytes_u64"`
		SendBroadcastCountU64 int    `json:"Send.BroadcastCount_u64"`
		SendUnicastBytesU64   int    `json:"Send.UnicastBytes_u64"`
		SendUnicastCountU64   int    `json:"Send.UnicastCount_u64"`
		UpdatedTimeDt         string `json:"UpdatedTime_dt"`
	} `json:"result"`
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
}

// EnumUser struct for using returned SON from JSON-RPC API
type EnumUser struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		HubNameStr string    `json:"HubName_str"`
		UserList   []ApiData `json:"UserList"`
	} `json:"result"`
}

// Define global variable for SoftEther API Authentication
var hubUser string       // SoftEther API username
var hubPass string       // SoftEther API password
var authEnc string       // Authorization header `Basic base64(username:password)`

func banner() {
	fmt.Printf("GoVPN-API / SoftEther Management\n")
	fmt.Printf("%s <options> <cmd>\n", os.Args[0])
	fmt.Println("---\tOptions\t---")
	fmt.Println("-p, --port\t Change service port of GoVPN-API")
	fmt.Println("---\tcmd\t---")
	fmt.Println("run\tRun the API & Web Server\t(Required)")
	fmt.Println("runapi\tRun the API only\t(Required)")
}

func main() {
	sport := "1337"
	if len(os.Args) <= 1 {
		banner()
		os.Exit(0)
	}
	for i, ch := range os.Args {
		if i == 0 {
			continue
		}
		if strings.Contains(ch, "-p") || strings.Contains(ch, "--port") {
			sport = os.Args[i+1]
		}else if strings.Contains(ch, "run") {
			// Serves files for Web-panel
			fs := http.FileServer(http.Dir("./files"))
			http.Handle("/api/", http.StripPrefix("/api/", fs))
		} else if strings.Contains(ch, "runapi") {
			continue
		}
	}
	fmt.Printf("running on localhost:%s\n", sport)

	// Handles REST-API Calls
	http.HandleFunc("/createUser", createUser)
	http.HandleFunc("/deleteUser", deleteUser)
	http.HandleFunc("/changePassword", changePassword)
	http.HandleFunc("/setExpireDate", setExpireDate)
	http.HandleFunc("/getUser", viewUser)
	http.HandleFunc("/listUsers", listUsers)
	http.HandleFunc("/init", setToken)
	port := fmt.Sprintf(":%s", sport)
	http.ListenAndServe(port, nil)
}

// Get current user count from the server
func getUserCount(serverip string) int {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("https://%s/api", serverip)
	mapA := map[string]interface{}{"HubName_str": "VPN"}
	mapB := map[string]interface{}{"jsonrpc": "2.0", "id": "rpc_call_id", "method": "EnumUser", "params": mapA}
	mapC, _ := json.Marshal(mapB)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(mapC))
	req.Header.Set("Authorization", authEnc)           // Set Authorization header with global authEnc
	req.Header.Set("Content-Type", "application/json") // Set content-type to json
	resp, _ := client.Do(req)                          // Execute request
	body, _ := ioutil.ReadAll(resp.Body)               // Read body to string
	defer resp.Body.Close()
	data := EnumUser{}
	re := bytes.NewReader([]byte(body))
	chatErr := json.NewDecoder(re).Decode(&data)
	if chatErr != nil {
		fmt.Println(chatErr)
	}
	count := len(data.Result.UserList)
	return count
}

// Returns type ApiData for changePassword, setExpireDate, viewUser
func getUser(username string, serverip string) ApiData {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("https://%s/api", serverip)
	mapA := map[string]interface{}{"HubName_str": "VPN", "Name_str": username}
	mapB := map[string]interface{}{"jsonrpc": "2.0", "id": "rpc_call_id", "method": "GetUser", "params": mapA}
	mapC, _ := json.Marshal(mapB) // Create map of mapA and mapB
	req, _ := http.NewRequest("POST", url, bytes.NewReader(mapC))
	req.Header.Set("Authorization", authEnc)           // Set Authorization header with global authEnc
	req.Header.Set("Content-Type", "application/json") // Set content-type to json
	resp, _ := client.Do(req)                          // Execute request
	body, _ := ioutil.ReadAll(resp.Body)               // Read body to string
	defer resp.Body.Close()
	data := ApiData{}
	re := bytes.NewReader([]byte(body))
	chatErr := json.NewDecoder(re).Decode(&data)
	if chatErr != nil {
		fmt.Println(chatErr)
	}
	return data
}

// Set hubUser, hubPass, and base64 encode authEnc for Authorization header
func setToken(w http.ResponseWriter, r *http.Request) {
	hubUser = r.URL.Query()["hubuser"][0]
	hubPass = r.URL.Query()["hubpass"][0]
	authEnc = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", hubUser, hubPass))))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["username"][0]
	passw := r.URL.Query()["password"][0]
	password, err := base64.StdEncoding.DecodeString(passw)
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	}
	serverip := r.URL.Query()["sip"][0]
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
		url := fmt.Sprintf("https://%s/api", serverip)
		client := &http.Client{Transport: tr}
		currentTime := time.Now()
		futureTime := currentTime.AddDate(0, 1, 0)
		mapA := map[string]interface{}{"HubName_str": "VPN", "Name_str": username, "CreatedTime_dt": currentTime.Format("2006-01-2T15:04:05.000Z"), "AuthType_u32": 1, "Auth_Password_str": string(password), "ExpireTime_dt": futureTime.Format("2006-01-02T15:04:05.000Z")}
		mapB := map[string]interface{}{"jsonrpc": "2.0", "id": "rpc_call_id", "method": "CreateUser", "params": mapA}
		mapC, _ := json.Marshal(mapB) // Create map of mapA and mapB
		req, _ := http.NewRequest("POST", url, bytes.NewReader(mapC))
		req.Header.Set("Authorization", authEnc)           // Set Authorization header with global authEnc
		req.Header.Set("Content-Type", "application/json") // Set content-type to json
		resp, _ := client.Do(req)                          // Execute request
		body, _ := ioutil.ReadAll(resp.Body)               // Read body to string
		defer resp.Body.Close()
		if strings.Contains(string(body), "Error code 66") {
			fmt.Fprintf(w, fmt.Sprintf("{\"status\": \"fail\", \"error\": %s}", string(body)))
			//fmt.Fprintf(w, "error")
			return
		}
		//fmt.Fprintf(w, "success")
		fmt.Fprintf(w, "{\"status\": \"pass\"}")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["username"][0]
	serverip := r.URL.Query()["sip"][0]
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("https://%s/api", serverip)
	mapA := map[string]interface{}{"HubName_str": "VPN", "Name_str": username}
	mapB := map[string]interface{}{"jsonrpc": "2.0", "id": "rpc_call_id", "method": "DeleteUser", "params": mapA}
	mapC, _ := json.Marshal(mapB) // Create map of mapA and mapB
	req, _ := http.NewRequest("POST", url, bytes.NewReader(mapC))
	req.Header.Set("Authorization", authEnc)           // Set Authorization header with global authEnc
	req.Header.Set("Content-Type", "application/json") // Set content-type to json
	resp, _ := client.Do(req)                          // Execute request
	body, _ := ioutil.ReadAll(resp.Body)               // Read body to string
	defer resp.Body.Close()
	if strings.Contains(string(body), "Error code") {
		fmt.Fprintf(w, fmt.Sprintf("{\"status\": \"fail\", \"error\": %s}", string(body)))
		return
	}
	fmt.Fprintf(w, "{\"status\": \"pass\"}")
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["username"][0]
	passw := r.URL.Query()["password"][0]
	password, _ := base64.StdEncoding.DecodeString(passw)
	serverip := r.URL.Query()["sip"][0]
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	url := fmt.Sprintf("https://%s/api", serverip)
	client := &http.Client{Transport: tr}
	data := getUser(username, serverip)
	createdtime := data.Result.CreatedTimeDt
	expiredtime := data.Result.ExpireTimeDt
	mapA := map[string]interface{}{"HubName_str": "VPN", "Name_str": username, "CreatedTime_dt": createdtime, "ExpireTime_dt": expiredtime, "AuthType_u32": 1, "Auth_Password_str": string(password)}
	mapB := map[string]interface{}{"jsonrpc": "2.0", "id": "rpc_call_id", "method": "SetUser", "params": mapA}
	mapC, _ := json.Marshal(mapB) // Create map of mapA and mapB
	req, _ := http.NewRequest("POST", url, bytes.NewReader(mapC))
	req.Header.Set("Authorization", authEnc)           // Set Authorization header with global authEnc
	req.Header.Set("Content-Type", "application/json") // Set content-type to json
	resp, _ := client.Do(req)                          // Execute request
	body, _ := ioutil.ReadAll(resp.Body)               // Read body to string
	defer resp.Body.Close()
	if strings.Contains(string(body), "Error code 66") {
		fmt.Fprintf(w, string(body))
		//fmt.Fprintf(w, "error")
		return
	}
	//fmt.Fprintf(w, "success")
	fmt.Fprintf(w, string(body))
}

func setExpireDate(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["username"][0]
	expdir := r.URL.Query()["expdate"][0]
	serverip := r.URL.Query()["sip"][0]
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	url := fmt.Sprintf("https://%s/api", serverip)
	client := &http.Client{Transport: tr}
	data := getUser(username, serverip)
	createdtime := data.Result.CreatedTimeDt
	securehash := data.Result.NtLmSecureHashBin
	hashedkey := data.Result.HashedKeyBin
	mapA := map[string]interface{}{"HubName_str": "VPN", "Name_str": username, "CreatedTime_dt": createdtime, "ExpireTime_dt": expdir, "AuthType_u32": 1, "NtLmSecureHash_bin": securehash, "HashedKey_bin": hashedkey}
	mapB := map[string]interface{}{"jsonrpc": "2.0", "id": "rpc_call_id", "method": "SetUser", "params": mapA}
	mapC, _ := json.Marshal(mapB) // Create map of mapA and mapB
	req, _ := http.NewRequest("POST", url, bytes.NewReader(mapC))
	req.Header.Set("Authorization", authEnc)           // Set Authorization header with global authEnc
	req.Header.Set("Content-Type", "application/json") // Set content-type to json
	resp, _ := client.Do(req)                          // Execute request
	body, _ := ioutil.ReadAll(resp.Body)               // Read body to string
	defer resp.Body.Close()
	fmt.Fprintf(w, string(body))
}

func viewUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["username"][0]
	serverip := r.URL.Query()["sip"][0]
	data := getUser(username, serverip)                        // getUser as ApiData
	m, _ := json.Marshal(map[string]interface{}{"data": data}) // JSON marshal ApiData in map
	fmt.Fprintf(w, fmt.Sprintf("%v", string(m)))
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	serverip := r.URL.Query()["sip"][0]
	url := fmt.Sprintf("https://%s/api", serverip)
	mapA := map[string]interface{}{"HubName_str": "VPN"}
	mapB := map[string]interface{}{"jsonrpc": "2.0", "id": "rpc_call_id", "method": "EnumUser", "params": mapA}
	mapC, _ := json.Marshal(mapB)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(mapC))
	req.Header.Set("Authorization", authEnc)           // Set Authorization header with global authEnc
	req.Header.Set("Content-Type", "application/json") // Set content-type to json
	resp, _ := client.Do(req)                          // Execute request
	body, _ := ioutil.ReadAll(resp.Body)               // Read body to string
	defer resp.Body.Close()
	fmt.Fprintf(w, string(body))
}
