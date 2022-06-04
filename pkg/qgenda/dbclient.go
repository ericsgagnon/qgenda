package qgenda

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/jmoiron/sqlx"
)

type DBClient interface {
	Open(cfg *DBClientConfig) (*sqlx.DB, error)
	Ping(ctx context.Context) (bool, error)
	CreateTable(ctx context.Context) (sql.Result, error)
	DropTable(ctx context.Context) (sql.Result, error)
	InsertRows(ctx context.Context) (sql.Result, error)
	QueryConstraints(ctx context.Context) error
}

type DBClientConfig struct {
	Name               string // descriptive name - used only for logs and reference
	Type               string // descriptive type - used only for logs and reference
	Driver             string // driver name - will be passed to sqlx.DB
	DataSourceName     string // only applicable if using DSN's
	ConnectionString   string // will be parsed to url
	ExpandEnvVars      bool   // whether or not to interpolate env vars of the form ${ENV_VAR} in connection string and dsn
	ExpandFileContents bool   // whether or not to interpolate file contents of the form {file:/path/to/file} in connection string and dsn
	// User             string   // prefer to reference env var or file contents by ${ENV_VAR_NAME} or {file:/path/to/file}
	// Password         string   // prefer to reference env var or file contents by ${ENV_VAR_NAME} or {file:/path/to/file}
	url *url.URL // let the program handle this
}

func ExampleDBClientConfig() DBClientConfig {
	cfg := DBClientConfig{
		Name:               "database",
		Type:               "odbc",
		Driver:             "odbc",
		ConnectionString:   "${DB_CONN_SCHEME}://${DB_USER}:${DB_PASSWORD}@${DB_HOSTNAME}:${DB_PORT}/${DB_DATABASE}?${DB_ARGUMENTS}",
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
		s = ExpandFileContents(s)
	}
	return s
}

// OpenDBConnection doesn't technically 'open' a real connection, it follows
// the go default of creating a DB struct that manages connections as needed
func OpenDBClient(cfg *DBClientConfig) (DBClient, error) {
	return nil, nil
}

// // ClientConfig is an attempt at a generic
// // way to config client parameters
// type DBClientConfig struct {
// 	Name             string
// 	Kind             string
// 	Scheme           string
// 	Host             string
// 	Port             int
// 	Database         string
// 	User             string
// 	Password         string
// 	Arguments        url.Values
// 	ConnectionString string
// 	DataSourceName   string   // only applicable if using DSN's
// 	url              *url.URL // let the program handle this
// }

// // ParseDBClientConfig first attempts to parse the connection string, if that is empty,
// // it parses a url as [Scheme]://[User]@[Password]:[Host]:[Port]/[Database]?[Arguments]
// // It will leave envvars and file references as-is until it passes the connection string
// // to the DB driver.
// func ParseDBClientConfig(cfg DBClientConfig) (DBClientConfig, error) {
// 	srcCfg := cfg
// 	// prefer connection string first
// 	var u *url.URL
// 	var err error
// 	if cfg.ConnectionString != "" && cfg.ConnectionString != "-" {
// 		u, err = url.Parse(cfg.ConnectionString)
// 		if err != nil {
// 			return srcCfg, err
// 		}
// 	} else {
// 		u = &url.URL{
// 			Scheme:   cfg.Scheme,
// 			Host:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
// 			Path:     cfg.Database,
// 			RawQuery: cfg.Arguments.Encode(),
// 		}
// 	}

// 	// populate cfg.url.UserInfo, if necessary
// 	user := u.User.Username()
// 	if user == "" {
// 		user = cfg.User
// 	}
// 	password, _ := u.User.Password()
// 	if password == "" {
// 		password = cfg.Password
// 	}
// 	u.User = url.UserPassword(user, password)
// 	cfg.url = u
// 	return cfg, nil
// }
