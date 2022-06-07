package main

import (
	"log"
	"os"

	"github.com/ericsgagnon/qgenda/app"

	_ "github.com/lib/pq"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Println(err)
	}
	if err := a.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
