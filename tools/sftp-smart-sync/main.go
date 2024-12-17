package main

import (
	"log"
	"os"

	cmd "github.com/ispapp/sftp-smart-sync/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:           "sftp-smart-sync",
		Usage:          "Synchronize files between local and remote paths",
		HelpName:       "sftp-smart-sync",
		Version:        "0.1.0",
		Description:    "Synchronize files between local and remote paths using SCP",
		DefaultCommand: "sync",
		Commands: []*cli.Command{
			cmd.SyncCmd,
			cmd.LnCmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
