package main

import (
	"flag"
	"fmt"
	"github.com/Neutronpay/core-notification-srv/config"
	srv "github.com/Neutronpay/core-notification-srv/srv"
	"net/http"
	"os"
	"os/signal"

	"github.com/Neutronpay/lib-go-common/logger"
)

func main() {

	environment := flag.String("e", "development", "")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
	}

	flag.Parse()
	// Initialize
	cfg := config.Init(*environment)
	log := logger.NewLogrusLogger(cfg.Base.Name, *environment)

	curSrv, err := srv.NewSrv(cfg, log)

	if err != nil {
		panic(err)
	}

	println(curSrv.Name())

	go func() {
		if err := curSrv.StartHttpListener(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatal(err, "failed to listen and serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	curSrv.Shutdown()

}
