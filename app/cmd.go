package app

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

//

func NewCLIApp(cfg *Config) (*cli.App, error) {
	appName := "qgenda-exporter"
	appEnvPrefix := strings.ReplaceAll(strings.ToUpper(appName), "-", "_")
	appVersion := "v0.2.0"
	cfgFile, err := ConfigFile(appName, "")
	if err != nil {
		return nil, err
	}
	// set default config file, based on os, etc.
	// configFile, err := configFile(appName, "")
	// if err != nil {
	// 	return &cli.App{}, err
	// }

	app := &cli.App{
		Name:    appName,
		Usage:   "export data from qgenda rest api https://restapi.qgenda.com/",
		Version: appVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"cfg", "c"},
				Value:   cfgFile,
				EnvVars: []string{fmt.Sprintf("%s_CONFIG", appEnvPrefix)},
			},
		},
		// Action: func(c *cli.Context) error {
		// 	fmt.Println("I'm Awesome")
		// 	fmt.Println("----------------------------------------------------------------------")
		// 	fmt.Println(configFile)
		// 	fmt.Println("----------------------------------------------------------------------")
		// 	fmt.Printf("test: %s\n", c.String("test"))
		// 	fmt.Println("----------------------------------------------------------------------")
		// 	// check config
		// 	if _, err := os.Stat(configFile); err != nil {
		// 		// fmt.Printf("Config File %s missing or inaccessible.\n", configFile)
		// 		fmt.Printf("Config File %s missing or inaccessible.\n", c.String("config"))
		// 		fmt.Printf("You can write an example one by:\n\t%s config example > %s\n", appName, c.String("config"))
		// 		return err
		// 	}

		// 	return nil
		// },
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: fmt.Sprintf("manage config for %s", appName),
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

							// exCfg := NewExampleConfig()

							// cfgOut, err := yaml.Marshal(exCfg)
							// if err != nil {
							// 	return err
							// }
							// fmt.Println(string(cfgOut))
							// fmt.Println("----------------------------------------------------------------------")
							return nil
						},
					},
					{
						Name:    "check",
						Aliases: []string{"validate"},
						Usage:   "validate and display config that will be used",
						Action: func(c *cli.Context) error {
							fmt.Println("checking...jk I'm not doing anything")
							fmt.Printf("TODO\n")
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

					fmt.Println("load config")
					// cfgFile
					// cfg, err := LoadAndParseConfig()
					// fmt.Println(c.String("config"))
					cfgFile, err := ConfigFile(appName, c.String("config"))
					if err != nil {
						return err
					}
					fmt.Println("Using Config File:  ", cfgFile)
					// fmt.Printf("%#v\n", c.App.Flags)
					// fmt.Println(appEnvPrefix)
					fmt.Println("iterate through data resources")

					fmt.Println("")
					// cfg := loadConfig()
					// parseConfig()
					// app := createApp(cfg)
					cfg, err := initConfig(c.String("config"))
					if err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Sprintf("%s", cfg)
					// fmt.Println(c.String("config"))
					// fmt.Printf("%v+\n", cfg)
					return nil
				},
			},
		},
		EnableBashCompletion: true,
	}

	return app, nil
}

// v.SetConfigName("config")
// it checks for the file, and writes an example config.yaml if none exists
// it returns the full config file name with path and an error
// if err := os.MkdirAll(path.Dir(cf), 0740); err != nil {
// 	return cf, err
// }

// // ConfigEndpoint is an attempt at a generic
// // way to config endpoint parameters
// type ConfigEndpoint struct {
// 	Name             string
// 	Kind             string
// 	Host             string
// 	Port             int
// 	Database         string
// 	User             string
// 	Password         string
// 	Arguments        url.Values
// 	ConnectionString string
// }

// NewConfigEndpoint returns an empty ConfigEndpoint
// to be populated later
// func NewConfigEndpoint() ConfigEndpoint {
// 	return ConfigEndpoint{
// 		Name:             "",
// 		Kind:             "",
// 		Host:             "",
// 		Database:         "",
// 		User:             "",
// 		Password:         "",
// 		Arguments:        url.Values{},
// 		ConnectionString: "",
// 	}
// }

// func NewCommand() (*cobra.Command, error) {

// 	var cfgFile string
// 	appName := "qgenda-exporter"
// 	appEnvPrefix := strings.ReplaceAll(strings.ToUpper(appName), "-", "_")
// 	appVersion := "v0.2.0"

// 	cmd := &cobra.Command{
// 		Use:   appName,
// 		Short: "export json data from qgenda rest api and import to postgres",
// 		Long: `export qgenda data as json, process it (minimal validation-type transformations),
// 		and load to postgres.`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Println("root command test for cobra.")
// 			fmt.Printf("%s %s\nenvironment variable prefix: %s\n", appName, appVersion, appEnvPrefix)
// 		},
// 	}

// 	// var cfgFile string
// 	cfgFile = "silly-nonsense"
// 	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", `config file, default is OS specific:
// 	On Unix systems, it returns $XDG_CONFIG_HOME as specified by
// 	https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
// 	non-empty, else $HOME/.config.
// 	On Darwin, it returns $HOME/Library/Application Support.
// 	On Windows, it returns %AppData%.
// 	On Plan 9, it returns $home/lib.
// 		`,
// 	)
// 	fmt.Println(cfgFile)
// 	fmt.Println("--------------------------------------")
// 	var cmdEcho = &cobra.Command{
// 		Use:   "echo [string to echo]",
// 		Short: "Echo anything to the screen",
// 		Long: `echo is for echoing anything back.
// 	Echo works a lot like print, except it has a child command.`,
// 		Args: cobra.MinimumNArgs(1),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Println("Echo: " + strings.Join(args, " "))
// 			fmt.Printf("Also, cfgFile: %s\n", cfgFile)
// 		},
// 	}
// 	cmd.AddCommand(cmdEcho)

// 	return cmd, nil
// }
