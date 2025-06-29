package main

import (
	gear_go "github.com/misnaged/gear-go"
	"github.com/misnaged/gear-go/cmd/root"
	"github.com/misnaged/gear-go/cmd/serve"
	"github.com/misnaged/gear-go/pkg/logger"
	"os"
)

func main() {
	app, err := gear_go.NewGear()

	if err != nil {
		logger.Log().Error("An error occurred", err)
		os.Exit(1)
	}

	rootCmd := root.Cmd()
	rootCmd.AddCommand(serve.Cmd(app))

	if err := rootCmd.Execute(); err != nil {
		logger.Log().Error("An error occurred", err)
		os.Exit(1)
	}
}
