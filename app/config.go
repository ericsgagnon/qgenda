package app

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/ericsgagnon/qgenda/pkg/qgenda"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config contains all app config parameters
type Config struct {
	App       AppConfig
	Cache     qgenda.CacheConfig
	Client    qgenda.ClientConfig
	Logger    LogConfig                        //zap.Config         `yaml:"-"`
	DBClients map[string]qgenda.DBClientConfig `yaml:"dbClients"`
	Data      map[string]qgenda.RequestConfig
	// DBClients map[string]url.URL `yaml:"dbClients"`
}

func DefaultConfig(ac *AppConfig) *Config {
	if ac == nil {
		ac = NewAppConfig()
	}
	cacheCfg, err := qgenda.NewCacheConfig(ac.Name)
	if err != nil {
		panic(err)
	}
	// sd, _ := time.Parse(time.RFC3339, "2006-01-01T00:00:00Z")
	// qcc := qgenda.ClientConfig{}
	return &Config{
		App:    *ac,
		Cache:  *cacheCfg,
		Client: *qgenda.DefaultClientConfig(),
		Logger: NewLogConfig(),
		DBClients: map[string]qgenda.DBClientConfig{
			"odbc": qgenda.ExampleDBClientConfig(),
		},
		Data: map[string]qgenda.RequestConfig{
			"schedule": *qgenda.DefaultScheduleRequestConfig(),
		},
	}
}

func LoadAndParseConfig(filename string) (*Config, error) {
	// read config file
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := Config{}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// NewExampleConfig creates an example config
func NewExampleConfig() Config {
	cfg := Config{
		// "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
		// Endpoints: map[string]ConfigEndpoint{
		// 	"postgres": {
		// 		Name:     "postgres",
		// 		Kind:     "odbc",
		// 		Host:     "db.example.com",
		// 		Port:     5432,
		// 		Database: "example",
		// 		User:     "user",
		// 		Password: "password",
		// 		Arguments: url.Values{
		// 			"ssl": []string{
		// 				"require",
		// 			}},
		// 	},
		// 	// "dest": {
		// 	// 	Name:     "dest",
		// 	// 	Kind:     "odbc",
		// 	// 	Host:     "db.example.com",
		// 	// 	Database: "example",
		// 	// 	User:     "user",
		// 	// 	Password: "password",
		// 	// 	Arguments: url.Values{
		// 	// 		"ssl": []string{
		// 	// 			"require",
		// 	// 		}},
		// 	// },
		// },
	}
	// cfg.Endpoints["src"].ConnectionString =
	return cfg
}

// ConfigFile is a convenience function that returns an os dependent
// default for /path/to/config/file or returns cf if not nil
func ConfigFile(app string, filename string) (string, error) {
	if filename == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return filename, err
		}
		filename = path.Join(userConfigDir, app, "config.yaml")
	}
	// os.IsExist()
	return filename, nil
}

// ExpandEnvVars substitutes environment variables of the form ${ENV_VAR_NAME}
// if you have characters that need to be escaped, they should be surrounded in
// quotes in the source string.
func ExpandEnvVars(s string) string {
	re := regexp.MustCompile(`\$\{.+\}`)

	envvars := map[string]string{}
	for _, m := range re.FindAllString(s, -1) {
		mre := regexp.MustCompile(`[${}]`)
		mtrimmed := mre.ReplaceAllString(m, "")
		// fmt.Printf("%s:\t%s\n", mtrimmed, os.Getenv(mtrimmed))
		envvars[m] = os.Getenv(mtrimmed)
	}

	for k, v := range envvars {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

func ConfigToYAML(cfg Config) (string, error) {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//	func (c Config) MarshalYAML() (any, error) {
//		return ConfigToYAML(c)
//	}
//
// initConfig reads in config file / Env vars
// and returns a Config
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
