package cmd

import (
	"strings"

	"github.com/urfave/cli/v2"
)

func NewApp() (*cli.App, error) {
	appName := "cli-app"
	appEnvPrefix := strings.ReplaceAll(strings.ToUpper(appName), "-", "_")
	appVersion := "v0.1.0"

	configFile, err := config.ConfigFile(appName, "")
	if err != nil {
		return &cli.App{}, err
	}
	app := &cli.App{
		Name:    appName,
		Usage:   "app usage description",
		Version: appVersion,
	}

}
