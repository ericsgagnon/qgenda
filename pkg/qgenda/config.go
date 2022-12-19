package qgenda

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// config manages configuration for the app

// Config contains all app config parameters
type Config struct {
	ApiVersion string
	Kind       string
	Name       string
	App        AppConfig
	Endpoints  map[string]*Endpoint
}

// AppConfig hold user defined inputs for the overall application
type AppConfig struct {
	Log       zap.Config
	Cache     CacheConfig
	ItemLimit int
	Timeout   time.Duration
	Retries   int
	Timezone  string
	Tags      map[string]string
	Objects   map[string]RequestConfig
}

// Endpoint is an attempt at a generic
// way to config endpoint parameters
type Endpoint struct {
	Name             string
	Type             string
	Driver           string
	DataSourceName   string // only applicable if using DSN's
	ConnectionString string // will be parsed to url
	User             string
	Password         string   `yaml:"-"`
	url              *url.URL // let the program handle this
}

// ConfigFile is a convenience function that returns an os dependent
// default for /path/to/config/file or returns cf if not nil
func ConfigFile(an string, cf string) (string, error) {
	if cf == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return cf, err
		}
		cf = path.Join(userConfigDir, an, "config.yaml")
	}
	return cf, nil
}

// InitConfig reads in config file / Env vars
// and returns a config struct
func InitConfig(cf string) (*Config, error) {
	c := Config{}
	v := viper.New()
	v.SetConfigFile(cf)
	// v.AutomaticEnv() // read in environment variables that match - maybe later

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error using config file: %s\n", v.ConfigFileUsed())
		return &c, err
	}
	fmt.Printf("Using config file: %s\n", v.ConfigFileUsed())

	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("Error parsing config file: %s", err)
		return &c, err
	}
	return &c, nil
}

// NewExampleConfig creates an example config
func NewExampleConfig() Config {
	cfg := Config{
		Endpoints: map[string]*Endpoint{
			"src": {
				Name:             "src",
				Type:             "rest",
				ConnectionString: "https://api.example.com/api/v1",
				User:             "user",
				Password:         "password",
			},
			"dest": {
				Name:             "dest",
				Type:             "postgres",
				ConnectionString: "postgres://db.example.com:5432/dbname?ssl=require",
				User:             "user",
				Password:         "password",
			},
		},
	}
	for _, v := range cfg.Endpoints {
		v.ParseURL()
	}
	return cfg
}

func (c *Endpoint) ParseURL() error {
	u, err := url.Parse(c.ConnectionString)
	if err != nil {
		return err
	}
	c.url = u
	return nil
}

// NewEndpoint returns an empty ConfigEndpoint
// to be populated later
func NewEndpoint() *Endpoint {
	return &Endpoint{}
}

// // Configure handles all configuration for the app
// func (c *Config) Configure() error {

// 	// read and parse config file using viper
// 	vc := viper.New()
// 	vc.SetConfigFile(configFile)
// 	vc.SetEnvPrefix("")
// 	vc.AutomaticEnv()
// 	if err := vc.ReadInConfig(); err != nil {
// 		return err
// 	}

// 	//configure each element
// 	var appCfg AppConfig
// 	vcApp := vc.Sub("app")
// 	vcApp.SetEnvPrefix("")
// 	vcApp.AutomaticEnv()
// 	vcApp.Unmarshal(&appCfg)
// 	c.App = appCfg

// 	c.Source = DatabaseConfigure(vc, "source")
// 	// fmt.Printf("%+v\n", c.Source)
// 	c.Destination = DatabaseConfigure(vc, "destination")
// 	// fmt.Printf("%+v\n", c.Destination)
// 	return nil
// }

// // DatabaseConfigure takes values from the given key to configure the database
// func DatabaseConfigure(v *viper.Viper, dbKey string) DatabaseConfig {
// 	var dbConfig DatabaseConfig
// 	vdb := v.Sub(dbKey)
// 	vdb.SetEnvPrefix(vdb.GetString("environmentPrefix"))
// 	vdb.AutomaticEnv()
// 	vdb.Unmarshal(&dbConfig)
// 	u := &url.URL{
// 		Scheme:   dbConfig.Scheme,
// 		User:     url.UserPassword(dbConfig.Username, dbConfig.Password),
// 		Host:     fmt.Sprintf("%s:%s", dbConfig.Host, dbConfig.Port),
// 		Path:     dbConfig.Database,
// 		RawQuery: dbConfig.Query,
// 	}
// 	dbConfig.URL = *u
// 	return dbConfig
// }
