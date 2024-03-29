package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/ericsgagnon/qgenda/pkg/qgenda"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

//

func NewCommand(app *App) (*cli.App, error) {

	appEnvPrefix := strings.ReplaceAll(strings.ToUpper(app.Config.App.Name), "-", "_")
	cfgFile, err := ConfigFile(app.Config.App.Name, "")
	if err != nil {
		return nil, err
	}

	cmd := &cli.App{
		Name:    app.Config.App.Name,
		Usage:   "export data from qgenda rest api https://restapi.qgenda.com/",
		Version: app.Config.App.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"cfg", "c"},
				Value:   cfgFile,
				EnvVars: []string{fmt.Sprintf("%s_CONFIG", appEnvPrefix)},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: fmt.Sprintf("manage config for %s", app.Config.App.Name),
				Subcommands: []*cli.Command{
					{
						Name:    "example",
						Aliases: []string{"ex"},
						Usage:   "show an example config file",
						Action: func(c *cli.Context) error {
							// fmt.Println("You asked for it!")

							cfg := DefaultConfig(nil)
							cfgOut, err := yaml.Marshal(cfg)
							if err != nil {
								return err
							}
							fmt.Println(string(cfgOut))

							return nil
						},
					},
					{
						Name:    "check",
						Aliases: []string{"validate"},
						Usage:   "validate and display config that will be used",
						Action: func(c *cli.Context) error {
							cfg, err := LoadAndParseConfig(c.String("config"))
							if err != nil {
								fmt.Println(err)
								return err
							}
							cfgYAML, err := ConfigToYAML(*cfg)
							if err != nil {
								fmt.Println(err)
								return err
							}
							fmt.Printf("Config file at %s appears valid\nbelow is a yaml representation of the config that is parsed from it:\n%s\n", c.String("config"), cfgYAML)

							return nil
						},
					},
				},
			},
			{
				Name:  "run",
				Usage: "extract, process, and load qgenda into target db",
				Action: func(c *cli.Context) error {
					fmt.Println("Let's do this.")
					fmt.Println("load config file:  ", c.String("config"))
					ctx := context.Background()

					cfg, err := LoadAndParseConfig(c.String("config"))
					if err != nil {
						fmt.Println(err)
						return err
					}
					app.Config = cfg

					// new qgenda client
					fmt.Println("Config qgenda client")
					app.Client, err = qgenda.NewClient(&app.Config.Client)
					if err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Println("authorize qgenda client")
					if err := app.Client.Auth(); err != nil {
						fmt.Println(err)
						return err
					}

					// establish db connections
					pgcfg := cfg.DBClients["postgres"]
					// for k, v := range cfg.DBClients {
					// 	fmt.Printf("Name: %s\t%#v\n", k, v)

					// }
					db, err := NewDBClient(&pgcfg)
					if err != nil {
						return err
					}
					defer db.Close()
					app.DBClients["postgres"] = db
					if err := app.ExecSchedulePipeline(ctx); err != nil {
						return err
					}
					if err := app.ExecStaffPipeline(ctx); err != nil {
						return err
					}
					return nil
				},
			},
		},
		EnableBashCompletion: true,
	}

	app.Command = cmd

	return cmd, nil
}
