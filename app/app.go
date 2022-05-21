package app

import (
	"github.com/ericsgagnon/qgenda/pkg/qgenda"
	"github.com/spf13/cobra"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type App struct {
	Config    any
	App       *cli.App
	Command   *cobra.Command
	Logger    *zap.Logger
	Client    *qgenda.Client
	Endpoints []*qgenda.Endpoint
}
