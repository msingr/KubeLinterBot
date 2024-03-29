// Package main reads config files, and contains the hook-receiving server.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aMMokschaf/KubeLinterBot/internal/callkubelinter"
	"github.com/aMMokschaf/KubeLinterBot/internal/config"
	"github.com/aMMokschaf/KubeLinterBot/internal/server"
)

var (
	logger = log.New(os.Stdout, "", 0)
)

// type kubeLinterBot struct {
// 	logger *log.Logger

// }

// func NewKubeLinterBot(logger *log.Logger) *kubeLinterBot {
// 	return &kubeLinterBot{
// 		logger: logger,
// 	}
// }

// func (b *kubeLinterBot) initialize() error {

// }

// mainCmd sets up a logger, a webHookServer, prints the address and port, starts the server
func mainCmd() error {
	// TODO argument for config file
	// TODO check if cfg-file exists
	cfg, err := config.OptionParser()
	if err != nil {
		return fmt.Errorf("Error reading configuration file: %w", err)
	}

	err = callkubelinter.CheckForKubeLinterBinary()
	if err != nil {
		return fmt.Errorf("checking for kube-linter binary: %w", err)
	}

	// TODO: implement check if token is actually valid, not just "empty"
	if cfg.User.AccessToken == "" {
		cfg, err = config.OptionParser()
		if err != nil {
			return fmt.Errorf("Could not read configuration-file: %w. Please copy the file './samples/kube-linter-bot-configuration.yaml' to kube-linter-bots directory.", err)
		}
	}
	webHookServ := server.SetupServer(logger, *cfg)
	logger.Printf("KubeLinterBot is listening on http://localhost%s\n", webHookServ.Addr)
	webHookServ.ListenAndServe()
	return nil
}

func main() {
	if err := mainCmd(); err != nil {
		logger.Fatalln(err)
	}
}
