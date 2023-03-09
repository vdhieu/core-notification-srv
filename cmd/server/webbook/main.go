package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/Neutronpay/core-notification-srv/config"
	"github.com/Neutronpay/core-notification-srv/notifsrv"
	"github.com/Neutronpay/lib-go-common/logger"
)

func main() {
	c := config.GetConfig(os.Getenv("ENV"))
	log := logger.NewLogrusLogger("core_notification", c.Base.Env)
	srv, err := notifsrv.New(c, log)

	if err != nil {
		panic(err)
	}

	println(srv.Name())

	go func() {
		if err := srv.StartHttpListener(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatal(err, "failed to listen and serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	srv.Shutdown()

}
