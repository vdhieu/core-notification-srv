package main

import (
	"flag"
	"log"

	"github.com/Neutronpay/core-notification-srv/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	conn := flag.String("conn", "", "connection string to db")
	flag.Parse()
	//logger.GetTraceLogger().Info("Start migration")
	if *conn == "" {
		log.Fatalln("conn flag required")
	}
	db, err := gorm.Open(postgres.Open(*conn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	// Replace `&Product{}, &User{}` with the models of your application.
	err = db.AutoMigrate(
		&model.WebhookInfo{},
	)
	if err != nil {
		log.Fatalln(err)
	}
}
