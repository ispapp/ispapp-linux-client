package lib

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

// Watch represents a file watcher
type Watch struct {
	LocalPath  string
	RemotePath string
	Interval   time.Duration
}

// Remote represents a remote server
type Remote struct {
	Host     string
	Port     string
	Username string
	Key      string
	Password string
}

// Monitor struct that Sync all pairs
type Monitor struct {
	Pairs  []Watch
	Remote Remote
	scp    *scp.Client
	ssh    *ssh.Client
	m      sync.Mutex
}

func (m *Monitor) SyncAll() {
	for _, pair := range m.Pairs {
		go pair.WatchLocalFile(m.scp, m.m)
		go pair.WatchRemoteFile(m.scp, m.m)
	}
}

// NewSCPClient initializes a new SCP client
func (m *Monitor) NewSCPClient() (*scp.Client, *ssh.Client, error) {
	// Create a new SSH client
	var clientConfig ssh.ClientConfig
	var err error
	if m.Remote.Password != "" {
		clientConfig, err = auth.PrivateKey(m.Remote.Username, m.Remote.Key, ssh.InsecureIgnoreHostKey())
		if err != nil {
			return nil, nil, err
		}
	} else if m.Remote.Key != "" {
		clientConfig, err = auth.PasswordKey(m.Remote.Username, m.Remote.Password, ssh.InsecureIgnoreHostKey())
		if err != nil {
			return nil, nil, err
		}
	} else {
		return nil, nil, errors.New("no password or key provided")
	}

	client := scp.NewClient(m.Remote.Host+":"+m.Remote.Port, &clientConfig)
	if client.Connect() != nil {
		return nil, nil, err
	}
	sshClien := client.SSHClient()
	return &client, sshClien, nil
}

func (m *Monitor) CheckPairs() error {
	for _, pair := range m.Pairs {
		// Check if local file exists
		if _, err := os.Stat(pair.LocalPath); os.IsNotExist(err) {
			log.Fatalf("Local file does not exist: %s\n", pair.LocalPath)
			return err
		}
		// Check if remote file exists
		if ok, payload, err := m.ssh.SendRequest("stat", true, []byte(pair.RemotePath)); err != nil || !ok || len(payload) == 0 {
			log.Fatalf("Remote file does not exist: %s\n", pair.RemotePath)
			return err
		}
	}
}

// Initializes a new Monitor
func NewMonitor(pairs []Watch, remote Remote) (*Monitor, error) {
	m := Monitor{
		Pairs:  pairs,
		Remote: remote,
	}
	var err error
	m.scp, m.ssh, err = m.NewSCPClient()
	if err != nil {
		return nil, err
	}
	if err := m.CheckPairs(); err != nil {
		return nil, err
	}
	return &m, nil
}

// WatchLocalFile watches the local file
func (w *Watch) WatchLocalFile(client *scp.Client, m sync.Mutex) {
	for {
		time.Sleep(w.Interval)
		m.Lock()
		err := client.CopyFile(w.LocalPath, w.RemotePath)
		if err != nil {
			panic(err)
		}
		m.Unlock()
	}
}

// WatchLocalFile watches the local file
func (w *Watch) WatchRemoteFile(client *scp.Client, m sync.Mutex) {
	for {
		// time.Sleep(w.Interval)
		// m.Lock()
		// // err := client.CopyFile(w.LocalPath, w.RemotePath)
		// if err != nil {
		// 	panic(err)
		// }
		// m.Unlock()
	}
}
