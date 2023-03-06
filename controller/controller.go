package controller

import (
	"github.com/Neutronpay/lib-go-common/comm/txncomm"
	"github.com/Neutronpay/lib-go-common/logger"
	"github.com/Neutronpay/lib-go-common/store"
	"github.com/go-redis/redis/v8"
)

type Controller struct {
	logger       logger.Logger
	txnStatusPub txncomm.TxnStatusPublisher
	//
	dbStore     *store.Store
	redisClient *redis.Client
	/*  Any service specific controllers here, e.g.:
		InvoiceController *InvoiceController
	    PayoutController  *PayoutController
	*/
}

func NewController(logger logger.Logger, db *store.Store, redisClient *redis.Client) (*Controller, error) {
	return &Controller{
		logger:      logger,
		dbStore:     db,
		redisClient: redisClient,
	}, nil
}

func (c *Controller) SetTxnPublisher(txnStatusPub txncomm.TxnStatusPublisher) {
	c.txnStatusPub = txnStatusPub
}

/*
  PUT HANDLER LOGIC HERE FOR THE SERVICE FOR EACH CMD TO BE HANDLED



*/
