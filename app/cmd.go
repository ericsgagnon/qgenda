package app

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func NewCommand() (*cobra.Command, error) {

	var cfgFile string
	appName := "qgenda-exporter"
	appEnvPrefix := strings.ReplaceAll(strings.ToUpper(appName), "-", "_")
	appVersion := "v0.2.0"

	cmd := &cobra.Command{
		Use:   appName,
		Short: "export json data from qgenda rest api and import to postgres",
		Long: `export qgenda data as json, process it (minimal validation-type transformations), 
		and load to postgres.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("root command test for cobra.")
			fmt.Printf("%s %s\nenvironment variable prefix: %s\n", appName, appVersion, appEnvPrefix)
		},
	}

	// var cfgFile string
	cfgFile = "silly-nonsense"
	cmd.PersistentFlags().StringVar(&cfgFile, "config-file", "", `config file, default is OS specific: 
	On Unix systems, it returns $XDG_CONFIG_HOME as specified by
	https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
	non-empty, else $HOME/.config.
	On Darwin, it returns $HOME/Library/Application Support.
	On Windows, it returns %AppData%.
	On Plan 9, it returns $home/lib.
		`,
	)
	fmt.Println(cfgFile)
	fmt.Println("--------------------------------------")
	var cmdEcho = &cobra.Command{
		Use:   "echo [string to echo]",
		Short: "Echo anything to the screen",
		Long: `echo is for echoing anything back.
	Echo works a lot like print, except it has a child command.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Echo: " + strings.Join(args, " "))
			fmt.Printf("Also, cfgFile: %s\n", cfgFile)
		},
	}
	cmd.AddCommand(cmdEcho)

	return cmd, nil
}

func NewApp() (*cli.App, error) {
	appName := "qgenda-exporter"
	appEnvPrefix := strings.ReplaceAll(strings.ToUpper(appName), "-", "_")
	appVersion := "v0.2.0"

	// set default config file, based on os, etc.
	configFile, err := configFile(appName, "")
	if err != nil {
		return &cli.App{}, err
	}

	app := &cli.App{
		Name:    appName,
		Usage:   "export data from qgenda rest api https://restapi.qgenda.com/",
		Version: appVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-file",
				Aliases:     []string{"f"},
				Value:       configFile,
				Destination: &configFile,
				EnvVars:     []string{fmt.Sprintf("%s_CONFIG_FILE", appEnvPrefix)},
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("I'm Awesome")
			fmt.Println("----------------------------------------------------------------------")
			fmt.Println(configFile)
			fmt.Println("----------------------------------------------------------------------")
			// check config
			if _, err := os.Stat(configFile); err != nil {
				// fmt.Printf("Config File %s missing or inaccessible.\n", configFile)
				fmt.Printf("Config File %s missing or inaccessible.\n", c.String("config-file"))
				fmt.Printf("You can write an example one by:\n\t%s config example > %s\n", appName, c.String("config-file"))
				return err
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "Create or check config file",
				Subcommands: []*cli.Command{
					{
						Name:  "example",
						Usage: "show an example config file",
						Action: func(c *cli.Context) error {
							// fmt.Println("You asked for it!")
							exCfg := NewExampleConfig()

							cfgOut, err := yaml.Marshal(exCfg)
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
							fmt.Println("checking...jk I'm not doing anything")
							fmt.Printf("TODO\n")
							return nil
						},
					},
				},
			},
			{
				Name:  "run",
				Usage: "Get er done",
				Action: func(c *cli.Context) error {
					fmt.Println("Let's do this.")
					// cfg := loadConfig()
					// parseConfig()
					// app := createApp(cfg)
					cfg, err := initConfig(c.String("config-file"))
					if err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Sprintf("%s", cfg)
					// fmt.Println(c.String("config-file"))
					// fmt.Printf("%v+\n", cfg)
					return nil
				},
			},
		},
		EnableBashCompletion: true,
	}

	return app, nil
}

// configFile is a convenience function that returns an os dependent
// default for /path/to/config/file or returns cf if not nil
func configFile(app string, filename string) (string, error) {
	if filename == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return filename, err
		}
		filename = path.Join(userConfigDir, app, "config.yaml")
	}
	return filename, nil
}

// v.SetConfigName("config")
// it checks for the file, and writes an example config.yaml if none exists
// it returns the full config file name with path and an error
// if err := os.MkdirAll(path.Dir(cf), 0740); err != nil {
// 	return cf, err
// }

// initConfig reads in config file / Env vars
// and returns a config struct
func initConfig(cf string) (Config, error) {
	c := Config{}
	v := viper.New()
	v.SetConfigFile(cf)
	// v.AutomaticEnv() // read in environment variables that match - maybe later

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error using config file: %s\n", v.ConfigFileUsed())
		return c, err
	}
	fmt.Printf("Using config file: %s\n", v.ConfigFileUsed())

	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("Error parsing config file: %s", err)
		return c, err
	}
	return c, nil
}

// Config contains all app config parameters
type Config struct {
	ApiVersion string
	Name       string
	Endpoints  map[string]ConfigEndpoint
}

// NewExampleConfig creates an example config
func NewExampleConfig() Config {

	cfg := Config{

		Endpoints: map[string]ConfigEndpoint{
			"src": {
				Name:     "src",
				Kind:     "odbc",
				Host:     "db.example.com",
				Database: "example",
				User:     "user",
				Password: "password",
				Arguments: url.Values{
					"ssl": []string{
						"require",
					}},
			},
			"dest": {
				Name:     "dest",
				Kind:     "odbc",
				Host:     "db.example.com",
				Database: "example",
				User:     "user",
				Password: "password",
				Arguments: url.Values{
					"ssl": []string{
						"require",
					}},
			},
		},
	}
	// cfg.Endpoints["src"].ConnectionString =
	return cfg
}

// ConfigEndpoint is an attempt at a generic
// way to config endpoint parameters
type ConfigEndpoint struct {
	Name             string
	Kind             string
	Host             string
	Database         string
	User             string
	Password         string
	Arguments        url.Values
	ConnectionString string
}

// NewConfigEndpoint returns an empty ConfigEndpoint
// to be populated later
func NewConfigEndpoint() ConfigEndpoint {
	return ConfigEndpoint{
		Name:             "",
		Kind:             "",
		Host:             "",
		Database:         "",
		User:             "",
		Password:         "",
		Arguments:        url.Values{},
		ConnectionString: "",
	}
}
