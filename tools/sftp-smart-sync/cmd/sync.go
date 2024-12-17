package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ispapp/sftp-smart-sync/lib"
	"github.com/ispapp/sftp-smart-sync/utils"
	"github.com/urfave/cli/v2"
)

var SyncCmd = &cli.Command{
	Name:  "sync",
	Usage: "Run sync with config file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Usage:    "Config file path",
			Aliases:  []string{"c"},
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		configFile := c.String("config")
		config, err := utils.ReadConfig(configFile)
		if err != nil {
			log.Fatalf("Failed to read config file: %v\n", err)
		}

		for _, path := range config.SyncPaths {
			watch := lib.NewWatch(path.Local, path.Remote, 10*time.Second)
			go watch.WatchLocalFile()
			go watch.MonitorRemote()
		}

		// Block main thread
		select {}
	},
}

var LnCmd = &cli.Command{
	Name:  "ln",
	Usage: "Run sync with file pairs",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "file",
			Usage:    "Local and remote file pairs",
			Aliases:  []string{"f"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "server",
			Usage:    "Remote server address (user:password@host:port)",
			Aliases:  []string{"s"},
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		remote := c.String("server")
		client, err := lib.NewSFTPClient(remote)
		if err != nil {
			log.Fatalf("Failed to initialize SFTP client: %v\n", err)
		}
		defer client.Close()
		filePairs := c.StringSlice("file")
		for _, pair := range filePairs {
			paths := strings.Split(pair, ":")
			if len(paths) != 2 {
				log.Fatalf("Invalid file pair: %s\n", pair)
			}
			localPath := paths[0]
			remotePath := paths[1]
			watch := lib.NewWatch(localPath, remotePath, 10*time.Second)
			fmt.Printf("Local: %s, Remote: %s\n", localPath, remotePath)
			// go watch.WatchLocalFile()
			// go watch.MonitorRemote()
		}

		// Block main thread
		select {}
	},
}
