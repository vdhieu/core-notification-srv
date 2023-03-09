package notifsrv

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Neutronpay/core-notification-srv/config"
	"github.com/Neutronpay/core-notification-srv/controller"
	"github.com/Neutronpay/core-notification-srv/store"
	"github.com/Neutronpay/lib-go-common/comm/txncomm"
	"github.com/Neutronpay/lib-go-common/conn/redis"
	libparsers "github.com/Neutronpay/lib-go-common/dto/parsers"
	"github.com/Neutronpay/lib-go-common/logger"
	"github.com/Neutronpay/lib-go-common/route"
)

type templateSrv struct {
	name           string
	logger         logger.Logger
	httpServer     *http.Server
	txnCommHandler txncomm.TxnRmqCommHandler

	// 	execPool           queue.ThreadedExecutorPool   // we might need to do this sooner than later
}

// New
// creates a new Template service, please note the http server is NOT being used here until we change the restful endpoints
// to be created and managed here.
func New(cfg *config.Config, logger logger.Logger) (srv *templateSrv, err error) {

	srv = &templateSrv{
		name:   cfg.Base.Name,
		logger: logger,
	}

	// TODO: CUSOTMIZE THE CONTENT
	txnDtoTf, err := setupDtoTransformer()
	if nil != err {
		logger.Errorf(err, "error setting up command parser", err.Error())
		return nil, err
	}

	rd := redis.ConnectRedis(context.Background(), cfg.RedisConf)
	// database
	db := store.ConnDb(cfg)
	store := store.New(db, cfg)

	// TODO: THIS NEEDS TO BE CUSTOMIZED
	ctrler, err := controller.NewController(logger, store, rd, txnDtoTf)
	if err != nil {
		return nil, err
	}

	_, err = txncomm.NewTxnStatusUpdateListener(cfg.RmqConf, cfg.Base.Name, ctrler)
	if err != nil {
		return nil, err
	}

	/*sch, err := controller.NewReqHander(srv, logger, cfg.Base, cfg.RmqConf, ctrler, txnDtoTf)
	if err != nil {
		return nil, err
	}

	srv.txnCommHandler = sch
	*/
	// setting up http server   TODO: ALSO NEED TO BE CUSTOMIZED
	var jwtSecret = []byte(cfg.JwtSecret)
	router := route.NewRoutes(controller.NewHttpRouteHandler(*cfg, store, srv.Logger(), txnDtoTf), jwtSecret)

	srv.logger.Infof("Starting server: %#v", cfg.Base)

	srv.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Base.Port),
		Handler: router,
	}

	return srv, err
}

// TODO: CUSTOMIZE TO WHATEVER PARSHER YOU NEED
// setupDtoTransformer  this creates a new transformer that can handle all known transformations in the platform
func setupDtoTransformer() (t libparsers.TxnDtoTransformer, err error) {

	// here we just create the map of ccy marshallers... eventually we'll do better
	tt, err := libparsers.NewTxnDtoTransformer()
	if err != nil {
		return nil, err
	}

	// move all these to app (trx srv, etc)
	destMarshaller := &libparsers.TxnStatusMarshaller{}
	err = tt.Register(destMarshaller)
	if err != nil {
		return nil, err
	}

	return tt, nil
}

// NPService interface

// Name returns the name or id of the service
func (ps *templateSrv) Name() string {
	return ps.name
}

func (ps *templateSrv) Logger() logger.Logger {
	return ps.logger
}

func (ps *templateSrv) StartHttpListener() (err error) {
	return ps.httpServer.ListenAndServe()
}

func (ps *templateSrv) Shutdown() error {
	// TODO: add rmq shutdown

	//ps.execPool.Stop()

	ps.logger.Info("Server Shutting Down")
	if err := ps.httpServer.Shutdown(context.Background()); err != nil {
		ps.logger.Error(err, "failed to shutdown server")
	}

	ps.logger.Info("Server Exit")

	return nil
}
