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
	Endpoints []qgenda.Endpoint
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

func NewApp() (*App, error) {
	app := &App{
		Config:    DefaultConfig(NewAppConfig()),
		DBClients: map[string]*sqlx.DB{},
	}

	cmd, err := NewCommand(app)
	if err != nil {
		return nil, err
	}
	app.Command = cmd
	// dbclients := map[string]sqlx.DB{}
	// app = App{
	// 	// Config: DefaultConfig(NewAppConfig()),
	// 	Config:    cfg,
	// 	Command:   cmd,
	// 	Client:    &qgenda.Client{},
	// 	DBClients: dbclients,
	// 	// Endpoints: []qgenda.Endpoint{},
	// }

	return app, nil
}

func (app *App) Run(args []string) error {
	return app.Command.Run(args)
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

func (app *App) ExecDataPipelines() error {

	// for _
	return nil
}
