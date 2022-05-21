package main

import (
	"log"

	"github.com/ericsgagnon/qgenda/app"
)

func main() {

	// a, err := app.NewApp()
	// if err != nil {
	// 	log.Println(err)
	// }
	// a.Command.Execute()
	cmd, err := app.NewCommand()
	if err != nil {
		log.Println(err)
	}
	cmd.Execute()
}
