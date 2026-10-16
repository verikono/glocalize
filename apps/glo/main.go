package main

import (
	"context"
	"log"
	"os"

	"github.com/verikono/glocalize/apps/glo/cmd"
)

func main() {
	app := cmd.NewApp()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
