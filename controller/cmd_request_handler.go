package controller

import (
	"fmt"
	"github.com/Neutronpay/lib-go-common/comm/txncomm"
	libconfig "github.com/Neutronpay/lib-go-common/config"
	"github.com/Neutronpay/lib-go-common/conn/rabbitmq"
	"github.com/Neutronpay/lib-go-common/dto/cmddto"
	"github.com/Neutronpay/lib-go-common/dto/cmddto/lnddto"
	"github.com/Neutronpay/lib-go-common/dto/cmddto/txndto"
	libparsers "github.com/Neutronpay/lib-go-common/dto/parsers"
	"github.com/Neutronpay/lib-go-common/logger"
	"github.com/Neutronpay/lib-go-common/server"
	"net/http"
	"time"
)

func NewReqHander(service server.NPService, logger logger.Logger, baseConfig libconfig.BaseConfig, rmqConfig libconfig.RabbitMqConn,
	ctrler *Controller, txnDtoTf libparsers.TxnDtoTransformer) (trh txncomm.TxnRmqCommHandler, err error) {

	sch := &srvCmdHandler{
		service:  service,
		logger:   logger,
		txnDtoTf: txnDtoTf,
	}

	err = sch.setupRabbitMqConsChans(rmqConfig)
	if nil != err {
		logger.Errorf(err, "error setting up rabbitmq for consuming transaction requests: %s config: %v", err.Error(), rmqConfig)
		return nil, err
	}

	txnStatusPub, err := txncomm.NewTxnStatusPublisher(baseConfig, rmqConfig)
	sch.txnStatusPub = txnStatusPub

	sch.ctrler = ctrler

	// TODO: this might need to be cut up so create publisher, create controller(publisher) then create
	// rabbitmq.MessageHandler(controller) to be more clean in steps
	sch.ctrler.SetTxnPublisher(sch)

	return trh, err
}

// rabbitmq.MessageHandler interface

// HandleMessage handles any incoming command requests
func (sch *srvCmdHandler) HandleMessage(reqCtx *cmddto.CmdPayloadContext) {
	// TODO: throw in a queue?
	sch.logger.Infof("HandleMessage processing message: %v", reqCtx)
	legReq, err := sch.parseCmdContext(reqCtx)
	if err != nil {
		err = fmt.Errorf("HandleMessage: error processing command: %s ", err.Error())
		sch.PublishCmdReqError(reqCtx, "", err)
	}

	// only so it won't complain about a non used variable in this code
	sch.logger.Infof("HandleMessage parsed leg: %v", legReq)
	var legResp txndto.LegResponse

	// TO BE REPLACED BY ACTUAL HANDLER CODE
	// the below is an example of handling of commands
	/*
		var legResp *paymedto.PayoutResponse
		switch reqCtx.Command() {
		case cmddto.CmdSendVndPaymePayout: // payout request
			l := legReq.(*paymedto.PayoutRequest)
			sch.logger.Infof("%v", l)
			legResp, err = sch.controller.CreatePayoutFromCmd(reqCtx, l)
		}

		if err != nil {
			sch.logger.Errorf(err, "HandleMessage Error processing %s %s, reason %s\n", reqCtx.Cmd, reqCtx.Payload, err.Error())
			sch.PublishCmdReqErrorWithReq(reqCtx, legReq, sch.service.Name(), err)
			return
		}

	*/
	resp := cmddto.NewCmdResponse(reqCtx.ReqId(), reqCtx.Cmd, http.StatusOK, sch.service.Name(), reqCtx.OrigInstance(), legResp, time.Now(), reqCtx.JwtPayload, reqCtx.JwtToken())
	sch.logger.Infof("HandleMessage Handle request done: %v", resp)
	err = sch.PublishStatusUpdate(resp)

}

// end rabbitmq.MessageHandler interface

type srvCmdHandler struct {
	service            server.NPService
	txnStatusPub       txncomm.TxnStatusPublisher
	txnRequestConsumer *rabbitmq.Client
	logger             logger.Logger
	txnDtoTf           libparsers.TxnDtoTransformer
	ctrler             *Controller
}

/*
This sets up the listeners for each channel.  The handler functions are the ones that need to deal with
sychronization and threading, so far up to here it's all the amop's thread still (or rather, the go func loop)
If this is not handed off to another thread this will potentially crash the main thread
*/
func (sch *srvCmdHandler) setupRabbitMqConsChans(c libconfig.RabbitMqConn) (err error) {

	chanName, err := cmddto.GetRmqChannelForCmd(cmddto.CmdSendVndPaymePayout)
	txnRequestConsumer, err := rabbitmq.CreateConsumerChannel(c, chanName, sch)
	if err != nil {
		sch.logger.Errorf(err, "error connecting to lnd txn consume channel: %s", chanName)
		return err
	}

	sch.txnRequestConsumer = txnRequestConsumer
	return nil
}

// txncomm.TxnStatusPublisher

// PublishStatusUpdate This seesentially is a delegate to send out the transaction response / status up
func (sch *srvCmdHandler) PublishStatusUpdate(resp *cmddto.CmdResponse) error {
	return sch.txnStatusPub.PublishStatusUpdate(resp)
}

func (sch *srvCmdHandler) PublishCmdReqError(reqCtx *cmddto.CmdPayloadContext, origSubId string, errIn error) (err error) {
	return sch.txnStatusPub.PublishCmdReqError(reqCtx, origSubId, errIn)
}

func (sch *srvCmdHandler) PublishCmdReqErrorWithReq(reqCtx *cmddto.CmdPayloadContext, req lnddto.LnRequest, subId string, errIn error) (err error) {
	return sch.txnStatusPub.PublishCmdReqErrorWithReq(reqCtx, req, subId, errIn)
}

// probably can also move to somewhere in lib common (parsers?)... this looks to be necessary all the time
func (sch *srvCmdHandler) parseCmdContext(reqCtx *cmddto.CmdPayloadContext) (legReq txndto.LegRequest, err error) {

	serializer, err := sch.txnDtoTf.GetCmdSerializer(reqCtx.Command())
	if err != nil {
		return nil, err
	}

	legReq, err = serializer.DeserializeLegRequest(reqCtx.Payload)
	return legReq, err

}
