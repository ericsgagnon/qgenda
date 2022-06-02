package app

import (
	"github.com/ericsgagnon/qgenda/pkg/qgenda"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type App struct {
	Config    *Config
	Command   *cli.App
	Logger    *zap.Logger
	Client    *qgenda.Client
	DBClients map[string]*sqlx.DB
	Endpoints []*qgenda.Endpoint
}

type AppConfig struct {
	Name    string
	Version string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		Name:    "qgenda-exporter",
		Version: "v0.2.0",
	}
}

func NewApp() *App {

	app := App{
		Config: DefaultConfig(NewAppConfig()),
	}

	return &app
}

// Config contains all app config parameters
// type Config struct {
// 	ApiVersion  string
// 	Name        string
// 	Cache       qgenda.CacheConfig
// 	Client      qgenda.ClientConfig
// 	Logger      zap.Config
// 	DBClients   map[string]url.URL `yaml:"dbClients"`
// 	DataObjects map[string]qgenda.RequestQueryFields
// }
