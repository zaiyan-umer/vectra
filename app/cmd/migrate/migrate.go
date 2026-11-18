package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/zaiyan-umer/vectra/app/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
	}

	cfg := config.MustLoad()

	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("migrate.new: %v", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Steps(1); err != nil {
			log.Fatalf("m.steps(1): %v", err)
		}
		log.Println("Migrations applied")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("m.steps(-1): %v", err)
		}
		log.Println("Migration rolled back")
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}