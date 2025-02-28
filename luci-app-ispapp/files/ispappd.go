package main

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"encoding/json"
	"net/http"
	"time"
	"strconv"
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
	httpClient          = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
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
	Data   json.RawMessage `json:"data,omitempty"`
}

// RpcdRequest represents the JSON structure for ubus calls
type RpcdRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

// RpcdResponse represents the JSON structure for ubus responses
type RpcdResponse struct {
	Result json.RawMessage `json:"result"`
}

func main() {
	log.Println("Starting ispappd...")
	for running {
		if !authenticate() {
			log.Println("Authentication failed, retrying...")
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
		
		switch event.Type {
		case "terminal":
			handleTerminal(event)
		case "getupdate":
			handleGetUpdate(event)
		case "getconfig":
			handleGetConfig(event)
		case "speedtest":
			handleSpeedTest(event)
		default:
			log.Printf("Unknown event type: %s", event.Type)
		}
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

func handleGetUpdate(event Event) {
	log.Printf("Handling getupdate event: %s", event.UUID)
	result, err := callRpcd("ispapp", "getupdate")
	response := Response{
		Type: "getupdate",
		UUID: event.UUID,
	}
	
	if err != nil {
		response.Stderr = fmt.Sprintf("Error calling getupdate: %v", err)
	} else {
		response.Data = result
	}
	
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending getupdate response: %v", err)
	}
}

func handleGetConfig(event Event) {
	log.Printf("Handling getconfig event: %s", event.UUID)
	result, err := callRpcd("ispapp", "getconfig")
	response := Response{
		Type: "getconfig",
		UUID: event.UUID,
	}
	
	if err != nil {
		response.Stderr = fmt.Sprintf("Error calling getconfig: %v", err)
	} else {
		response.Data = result
	}
	
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending getconfig response: %v", err)
	}
}

func handleSpeedTest(event Event) {
	log.Printf("Handling speedtest event: %s", event.UUID)
	
	// Notify client that speedtest has started
	startResponse := Response{
		Type: "speedtest_status",
		UUID: event.UUID,
		Stdout: "Starting speedtest...",
	}
	conn.WriteJSON(startResponse)
	
	// Call the speedtest function from ispapp rpcd - simple direct ubus call
	result, err := callRpcd("ispapp", "speedtest")
	
	response := Response{
		Type: "speedtest",
		UUID: event.UUID,
	}
	
	if err != nil {
		response.Stderr = fmt.Sprintf("Error running speedtest: %v", err)
	} else {
		response.Data = result
	}
	
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending speedtest response: %v", err)
	}
}

// callRpcd calls the rpcd service using ubus
func callRpcd(service, method string, params ...interface{}) (json.RawMessage, error) {
	// Build the ubus call command
	args := []string{"call", service, method}
	
	// Add parameters if provided
	var jsonParams string
	if len(params) > 0 && params[0] != nil {
		paramBytes, err := json.Marshal(params[0])
		if err != nil {
			return nil, fmt.Errorf("error marshalling params: %v", err)
		}
		jsonParams = string(paramBytes)
		args = append(args, jsonParams)
	}
	
	// Execute the ubus command
	log.Printf("Executing: ubus %s", strings.Join(args, " "))
	cmd := exec.Command("ubus", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error calling ubus: %v, output: %s", err, output)
	}
	
	// Return the raw JSON output
	return output, nil
}

func updateThread() {
	for running {
		time.Sleep(updateInterval)
		if conn == nil {
			continue
		}
		
		// Get update data from rpcd - simple direct ubus call
		updateData, err := callRpcd("ispapp", "getupdate")
		if err != nil {
			log.Printf("Error getting update data: %v", err)
			continue
		}
		
		event := Event{
			Type: "update",
			Data: updateData,
		}
		if err := conn.WriteJSON(event); err != nil {
			log.Printf("Error sending update: %v", err)
		}
	}
}

func configThread() {
	for running {
		time.Sleep(configInterval)
		if conn == nil {
			continue
		}
		
		// Get config data from rpcd - simple direct ubus call
		configData, err := callRpcd("ispapp", "getconfig")
		if err != nil {
			log.Printf("Error getting config data: %v", err)
			continue
		}
		
		event := Event{
			Type: "config",
			Data: configData,
		}
		if err := conn.WriteJSON(event); err != nil {
			log.Printf("Error sending config: %v", err)
		}
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
