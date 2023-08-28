package main

import (
	"PayWalletEngine/internal/db"
	"fmt"
	"log"
)

// Run - is going to be responsible for / the instantiation and startup of our / go application
func Run() error {
	fmt.Println("starting up the application...")
	dsn, err := db.LoadConfig()
	if err != nil {
		log.Println("LoadConfig Error")
		return err
	}
	database, err := db.NewDatabase(dsn)
	if err != nil {
		log.Println("Database Connection Failure")
		return err
	}
	if err := database.HealthCheck(); err != nil {
		return err
	}
	log.Println("Successfully connected to the database")

	if err := database.MigrateDB(); err != nil {
		log.Println("failed to setup database migrations")
		return err
	}
	return nil

}
func main() {
	fmt.Println("GO REST API Course")
	if err := Run(); err != nil {
		log.Println(err)
	}

}
