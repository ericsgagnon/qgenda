package app

import (
	"context"
	"fmt"
	"log"

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
		Version: "v0.4.0",
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

func (app *App) ExecDataPipelines() error {

	// for _
	return nil
}

func (app *App) ExecSchedulePipeline(ctx context.Context) error {

	// only one dbClient at present
	dbCfg := app.Config.DBClients["postgres"]
	db := app.DBClients["postgres"]

	rcCfg := app.Config.Data["schedule"]
	// rc, err := qgenda.GetPGStatus(ctx, db, qgenda.Staff{}, dbCfg.Schema, "", "")
	value := qgenda.Schedule{}
	result, err := qgenda.CreatePGTable(ctx, db, value, dbCfg.Schema, "")
	if err != nil {
		log.Println(result)
		return err
	}
	rc, err := value.GetPGStatus(ctx, db, dbCfg.Schema, "")
	if err != nil {
		return err
	}
	rc = rcCfg.Merge(rc)

	values, err := qgenda.GetSchedules(ctx, app.Client, rc)
	if err != nil {
		return err
	}

	if len(values) > 0 {
		if err := values.Process(); err != nil {
			return err
		}
	}

	preRowCount, err := qgenda.CountPGRows(ctx, db, values, dbCfg.Schema, "")
	if err != nil {
		return err
	}
	_, err = qgenda.BatchPutPG(ctx, db, 100, values, dbCfg.Schema, "")
	if err != nil {
		return err
	}
	// if result != nil {

	// 	// rowsAffected, _ := result.RowsAffected()
	// 	// log.Printf("%T Rows Inserted: %d", values, rowsAffected)

	// }
	rows, err := qgenda.CountPGRows(ctx, db, values, dbCfg.Schema, "")
	if err != nil {
		return err
	}
	log.Printf("Insert %T: preRowCount: %d\tpostRowCount: %d\trowInserted: %d", values, preRowCount, rows, rows-preRowCount)
	return nil
}

func (app *App) ExecStaffPipeline(ctx context.Context) error {
	fmt.Println("Starting Staffs")
	// only one dbClient at present
	dbCfg := app.Config.DBClients["postgres"]
	db := app.DBClients["postgres"]

	rcCfg := app.Config.Data["staff"]
	// rc, err := qgenda.GetPGStatus(ctx, db, qgenda.Staff{}, dbCfg.Schema, "", "")

	value := qgenda.Staff{}
	result, err := qgenda.CreatePGTable(ctx, db, value, dbCfg.Schema, "")
	if err != nil {
		log.Println(result)
		return err
	}

	rc, err := value.GetPGStatus(ctx, db, dbCfg.Schema, "")
	if err != nil {
		return err
	}
	rc = rcCfg.Merge(rc)

	values, err := qgenda.GetStaffs(ctx, app.Client, rc)
	if err != nil {
		return err
	}

	if len(values) > 0 {
		if err := values.Process(); err != nil {
			return err
		}
	}

	preRowCount, err := qgenda.CountPGRows(ctx, db, values, dbCfg.Schema, "")
	if err != nil {
		return err
	}

	_, err = qgenda.PutPG(ctx, db, values, dbCfg.Schema, "")
	if err != nil {
		return err
	}
	rows, err := qgenda.CountPGRows(ctx, db, values, dbCfg.Schema, "")
	if err != nil {
		return err
	}
	log.Printf("Insert %T: preRowCount: %d\tpostRowCount: %d\trowInserted: %d", values, preRowCount, rows, rows-preRowCount)

	// if result != nil {
	// 	rowsAffected, _ := result.RowsAffected()
	// 	log.Printf("%T Rows Inserted: %d", values, rowsAffected)
	// }

	// s := qgenda.Schedules{}
	// result, err := s.EPL(ctx,
	// 	app.Client,
	// 	&rqf,
	// 	app.DBClients["postgres"],
	// 	app.Config.DBClients["postgres"].Schema,
	// 	"schedule",
	// 	true)
	// if err != nil {
	// 	return err
	// }
	// if result != nil {
	// 	ra, _ := result.RowsAffected()
	// 	log.Printf("Schedule Rows Inserted: %d", ra)

	// }
	// data := []qgenda.Schedule{}
	// result, err := qgenda.GetProcessPutPGSchedules(ctx,
	// 	app.Client,
	// 	&rqf,
	// 	app.DBClients["postgres"],
	// 	app.Config.DBClients["postgres"].Schema,
	// 	"schedule",
	// 	true,
	// )
	// if err != nil {
	// 	log.Println(err)
	// }
	// if result != nil {
	// 	rowsAffected, _ := result.RowsAffected()
	// 	log.Printf("Schedule Rows Inserted: %d", rowsAffected)
	// }

	return nil
}

// func (app *App) ExecStaffMemberPipeline() error {

// 	ctx := context.Background()
// 	rqf := app.Config.Data["staffmember"]
// 	s := qgenda.Staffs{}
// 	result, err := s.EPL(ctx,
// 		app.Client,
// 		&rqf,
// 		app.DBClients["postgres"],
// 		app.Config.DBClients["postgres"].Schema,
// 		"staffmember",
// 		true)
// 	if err != nil {
// 		return err
// 	}
// 	if result != nil {
// 		ra, _ := result.RowsAffected()
// 		log.Printf("Schedule Rows Inserted: %d", ra)

// 	}
// 	// data := []qgenda.Schedule{}
// 	return nil
// }
