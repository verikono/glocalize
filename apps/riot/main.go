package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"github.com/verikono/glocalize/apps/riot/cmd"
)

func main() {
	_ = godotenv.Load()

	app := cmd.NewApp()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
