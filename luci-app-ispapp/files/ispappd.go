package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/digineo/go-uci"
	"github.com/gorilla/websocket"
)

var (
	serverURL           = "wss://prv.cloud.ispapp.co" // Replace with your server URL
	updateInterval      = 10 * time.Second
	configInterval      = 15 * time.Second
	healthCheckInterval = 20 * time.Second
	reconnectInterval   = 30 * time.Second
	running             = true
	conn                *websocket.Conn
	accessToken         = ""
	refreshToken        = ""
	configFilePath      = "/etc/config/ispapp"
	uciTree             = uci.NewTree(configFilePath)
)

type Event struct {
	Type string `json:"type"`
	UUID string `json:"uuid"`
	Data json.RawMessage `json:"data"`
}

type Response struct {
	Type   string `json:"type"`
	UUID   string `json:"uuid"`
	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`
}

func main() {
	for running {
		if !authenticate() {
			time.Sleep(reconnectInterval)
			continue
		}
		connectToServer()
		if conn != nil {
			go updateThread()
			go configThread()
			go healthCheckThread()
			handleMessages()
		}
		time.Sleep(reconnectInterval)
	}
}

func authenticate() bool {
	// Load tokens from config
	accessToken = loadConfig("accessToken")
	refreshToken = loadConfig("refreshToken")

	if accessToken != "" {
		return true
	}

	if refreshToken != "" {
		return refreshAccessToken()
	}

	return login()
}

func loadConfig(key string) string {
	values, ok := uciTree.Get("ispapp", "@settings[0]", key)
	if !ok || len(values) == 0 {
		log.Printf("Failed to load config key: %s", key)
		return ""
	}
	return values[0]
}

func refreshAccessToken() bool {
	host := loadConfig("Domain")
	port := loadConfig("ListenerPort")
	uri := fmt.Sprintf("https://%s:%s/auth/refresh?refreshToken=%s", host, port, refreshToken)
	resp, err := http.Get(uri)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Failed to refresh access token: %v", err)
		return false
	}
	defer resp.Body.Close()

	var tokens map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		log.Printf("Failed to decode token response: %v", err)
		return false
	}

	accessToken = tokens["accessToken"]
	refreshToken = tokens["refreshToken"]
	saveConfig("accessToken", accessToken)
	saveConfig("refreshToken", refreshToken)
	return true
}

func login() bool {
	host := loadConfig("Domain")
	port := loadConfig("ListenerPort")
	key := loadConfig("Key")
	login := loadConfig("login")

	if login == "" || login == "00000000-0000-0000-0000-000000000000" {
		login = gatherUniqueID()
		saveConfig("login", login)
	}

	uri := fmt.Sprintf("https://%s:%s/auth/login?login=%s&key=%s", host, port, login, key)
	resp, err := http.Get(uri)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Failed to login: %v", err)
		return false
	}
	defer resp.Body.Close()

	var tokens map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		log.Printf("Failed to decode login response: %v", err)
		return false
	}

	accessToken = tokens["accessToken"]
	refreshToken = tokens["refreshToken"]
	saveConfig("accessToken", accessToken)
	saveConfig("refreshToken", refreshToken)
	return true
}

func saveConfig(key, value string) {
	if sectionExists := uciTree.Set("ispapp", "@settings[0]", key, value); !sectionExists {
		_ = uciTree.AddSection("ispapp", "@settings[0]", "settings")
		_ = uciTree.Set("ispapp", "@settings[0]", key, value)
	}
	if err := uciTree.Commit(); err != nil {
		log.Printf("Failed to save config: %v", err)
	}
}

func connectToServer() {
	var err error
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+accessToken)
	conn, _, err = websocket.DefaultDialer.Dial(serverURL, headers)
	if err != nil {
		log.Printf("Failed to connect to server: %v", err)
		conn = nil
		return
	}
	log.Println("Connected to server")
}

func handleMessages() {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			conn.Close()
			conn = nil
			return
		}
		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}
		if event.Type == "terminal" {
			handleTerminal(event)
		}
		// if event.Type == "wireless" {
		// 	handleWireless(event)
		// }
	}
}

func handleTerminal(event Event) {
	var command []string
	if err := json.Unmarshal(event.Data, &command); err != nil {
		log.Printf("Error unmarshalling command: %v", err)
		return
	}
	cmd := exec.Command(command[0], command[1:]...)
	stdout, err := cmd.Output()
	response := Response{
		Type:   "terminal",
		UUID:   event.UUID,
		Stdout: string(stdout),
	}
	if err != nil {
		response.Stderr = err.Error()
	}
	conn.WriteJSON(response)
}

func updateThread() {
	for running {
		time.Sleep(updateInterval)
		if conn == nil {
			continue
		}
		event := Event{Type: "update"}
		conn.WriteJSON(event)
	}
}

func configThread() {
	for running {
		time.Sleep(configInterval)
		if conn == nil {
			continue
		}
		event := Event{Type: "config"}
		conn.WriteJSON(event)
	}
}

func healthCheckThread() {
	for running {
		time.Sleep(healthCheckInterval)
		if conn == nil {
			continue
		}
		event := Event{Type: "healthcheck"}
		conn.WriteJSON(event)
	}
}

func gatherUniqueID() string {
	// Function to get the MAC address using `ifconfig`
	getMacAddress := func() string {
		output, err := exec.Command("/sbin/ifconfig", "-a").Output()
		if err != nil {
			return "00:00:00:00:00:00"
		}
		macAddress := ""
		for _, line := range strings.Split(string(output), "\n") {
			if strings.Contains(line, "HWaddr") {
				macAddress = strings.TrimSpace(strings.Split(line, "HWaddr")[1])
				break
			}
		}
		if macAddress == "" {
			return "00:00:00:00:00:00"
		}
		return macAddress
	}

	// Gather necessary values
	macAddress := getMacAddress()
	cpuCores, _ := exec.Command("sh", "-c", "cat /proc/cpuinfo | grep 'processor' | wc -l").Output()
	ramSize := 0
	ramPartitions, _ := exec.Command("sh", "-c", "cat /proc/partitions | grep ram | wc -l").Output()
	if ramPartitionsInt := stringToInt(string(ramPartitions)); ramPartitionsInt > 0 {
		ramSize = ramPartitionsInt * 4096 // Assuming each ram partition is 4096 blocks
	}
	mtdSize := 0
	mtdPartitions, _ := exec.Command("sh", "-c", "cat /proc/partitions | grep mtdblock | wc -l").Output()
	if mtdPartitionsInt := stringToInt(string(mtdPartitions)); mtdPartitionsInt > 0 {
		mtdData, _ := exec.Command("sh", "-c", "cat /proc/partitions | grep mtdblock").Output()
		for _, line := range strings.Split(string(mtdData), "\n") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				mtdSize += stringToInt(parts[2])
			}
		}
	}
	bogomips, _ := exec.Command("sh", "-c", "cat /proc/cpuinfo | grep 'BogoMIPS' | head -n 1 | awk '{print $3}'").Output()
	numPorts, _ := exec.Command("sh", "-c", "ls /sys/class/net/ | grep 'eth' | wc -l").Output()

	// Combine all values into a string
	uniqueString := fmt.Sprintf("%s%d%d%d%s%d", macAddress, stringToInt(string(cpuCores)), ramSize, mtdSize, strings.TrimSpace(string(bogomips)), stringToInt(string(numPorts)))

	// Generate UUID format from the unique string
	hash := md5.Sum([]byte(uniqueString))
	hashString := hex.EncodeToString(hash[:])
	res := fmt.Sprintf("%s-%s-4%s-%s-%s", hashString[:8], hashString[8:12], hashString[12:16], hashString[16:20], hashString[20:])
	return res
}

func stringToInt(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}
