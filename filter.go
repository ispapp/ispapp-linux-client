package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
    "strings"
	"os"
    "sync"
	"time"
)

type Server struct {
	IP_HOST   string `json:"IP_HOST"`
	Reachable bool   `json:"reachable"`
    V3        bool   `json:"V3"`
    V2        bool   `json:"V2"`
    OPTIONS   string `json:"OPTIONS"`
    GB_S      string `json:"GB_S"`
    COUNTRY   string `json:"COUNTRY"`
    SITE      string `json:"SITE"`
    CONTINENT string `json:"CONTINENT"`
    PROVIDER  string `json:"PROVIDER"`
}
type FilteredServer struct {
	IP_HOST   string `json:"IP_HOST"`
    V3        bool   `json:"V3"`
    V2        bool   `json:"V2"`
    OPTIONS   string `json:"OPTIONS"`
    GB_S      string `json:"GB_S"`
    COUNTRY   string `json:"COUNTRY"`
    SITE      string `json:"SITE"`
    CONTINENT string `json:"CONTINENT"`
    PROVIDER  string `json:"PROVIDER"`
}

// IsIperfServerAlive checks if an iperf server is reachable on a given host and port.
func IsIperfServerAlive(host string, port string, timeout time.Duration) bool {
	address := net.JoinHostPort(host, port)

	// Try to establish a TCP connection within the timeout
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

func main() {
	inputFile := "./luci-app-ispapp/root/www/luci-static/resources/iperf"
	outputFile := "./luci-app-ispapp/root/www/luci-static/resources/iperf"

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	var servers []Server
	err = json.Unmarshal(data, &servers)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
    wg := sync.WaitGroup{}
	for i, server := range servers {
        wg.Add(1)
        go func(i int, server Server) {
            defer wg.Done()
            // server can be in this fomat: HOST -p PORT or HOST -p PORT-PORT or HOST where default PORT is 5201"
            hostparts := strings.Split(server.IP_HOST, " ")
            _server := hostparts[0]
            var port string
            if len(hostparts) > 1 {
                port = hostparts[1]
            } else {
                port = "5201"
            }
            if IsIperfServerAlive(_server, port, 5*time.Second) {
                fmt.Printf("\033[32m\u2713\033[0m %s\n", server.IP_HOST)
                servers[i].Reachable = true
            } else {
                fmt.Printf("\033[31m\u2717\033[0m %s\n", server.IP_HOST)
                servers[i].Reachable = false
            }
        }(i, server)
    }
    wg.Wait()
    // remove unreachable servers and Reachable key from struct 
    var filteredServers []FilteredServer
    for _, server := range servers {
        if server.Reachable {
            filteredServers = append(filteredServers, FilteredServer{
                IP_HOST:   server.IP_HOST,
                V3:        server.V3,
                V2:        server.V2,
                OPTIONS:   server.OPTIONS,
                GB_S:      server.GB_S,
                COUNTRY:   server.COUNTRY,
                SITE:      server.SITE,
                CONTINENT: server.CONTINENT,
                PROVIDER:  server.PROVIDER,
            })
        }
    }
	outputData, err := json.Marshal(filteredServers)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = ioutil.WriteFile(outputFile, outputData, 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}
}