package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/labstack/echo"
)

func main() {

	log.Println("setting up config file")
	setupConfigFile()

	log.Println("downloading database(s)")
	err := updateDB()
	if err != nil {
		log.Fatalln("database download failed")
	}

	server := echo.New()

	scheduler := gocron.NewScheduler(time.UTC)

	_, err = scheduler.Every(1).Day().At("12:00").Do(updateDB)

	if err != nil {
		log.Fatalf("error when scheduling database update job: %v", err)
	}
	// scheduler.WaitForScheduleAll()

	server.GET("/location/:ip", ip2location)
	server.GET("/update-db", func(ctx echo.Context) error { return updateDB() })

	scheduler.StartAsync()

	server.Start(":8080")

	scheduler.Stop()
}
