package app

import (
	"github.com/exiledavatar/gotoolkit/meta"
	"github.com/jmoiron/sqlx"
)

type DBClientConfig struct {
	Name               string // descriptive name - used only for logs and reference
	Type               string // descriptive type - used only for logs and reference
	Driver             string // driver name - will be passed to sqlx.DB
	DataSourceName     string // only applicable if using DSN's
	ConnectionString   string // will be parsed to url
	Schema             string // schema to use for this client
	ExpandEnvVars      bool   // whether or not to interpolate env vars of the form ${ENV_VAR} in connection string and dsn
	ExpandFileContents bool   // whether or not to interpolate file contents of the form {file:/path/to/file} in connection string and dsn
	// User             string   // prefer to reference env var or file contents by ${ENV_VAR_NAME} or {file:/path/to/file}
	// Password         string   // prefer to reference env var or file contents by ${ENV_VAR_NAME} or {file:/path/to/file}
	// url *url.URL // let the program handle this
}

func ExampleDBClientConfig() DBClientConfig {
	cfg := DBClientConfig{
		Name:               "database",
		Type:               "odbc",
		Driver:             "odbc",
		ConnectionString:   "${DB_CONN_SCHEME}://${DB_USER}:${DB_PASSWORD}@${DB_HOSTNAME}:${DB_PORT}/${DB_DATABASE}?${DB_ARGUMENTS}",
		Schema:             "qgenda",
		ExpandEnvVars:      true,
		ExpandFileContents: true,
		// User:             "${DB_USER}",
		// Password:         "${DB_PASSWORD}",
	}
	return cfg
}

func (cfg DBClientConfig) String() string {
	s := cfg.ConnectionString
	if cfg.ExpandEnvVars {
		s = ExpandEnvVars(s)
	}
	if cfg.ExpandFileContents {
		s = meta.ExpandFileContents(s)
	}
	return s
}

func NewDBClient(cfg *DBClientConfig) (*sqlx.DB, error) {
	connString := ExpandEnvVars(cfg.ConnectionString)
	// fmt.Printf("Driver: %s\t ConnString: %s\n", cfg.Driver, connString)
	return sqlx.Open(cfg.Driver, connString)
}
