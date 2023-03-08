package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Neutronpay/core-notification-srv/dto/extdto"
	"github.com/Neutronpay/core-notification-srv/entity"
	"github.com/Neutronpay/lib-go-common/comm/txncomm"
	"github.com/Neutronpay/lib-go-common/conn/rabbitmq"
	"github.com/Neutronpay/lib-go-common/dto/cmddto"
	"github.com/Neutronpay/lib-go-common/dto/cmddto/txndto"
	libparsers "github.com/Neutronpay/lib-go-common/dto/parsers"
	"github.com/Neutronpay/lib-go-common/logger"
	"github.com/Neutronpay/lib-go-common/signature"
	"github.com/Neutronpay/lib-go-common/store"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"io"

	"io/ioutil"
	"net/http"
)

type Controller struct {
	logger       logger.Logger
	txnStatusPub txncomm.TxnStatusPublisher
	//
	dbStore     *store.Store
	redisClient *redis.Client

	dtoParser libparsers.TxnDtoTransformer
	/*  Any service specific controllers here, e.g.:
		InvoiceController *InvoiceController
	    PayoutController  *PayoutController
	*/
	rabbitmq.MessageRespHandler
}

func NewController(logger logger.Logger, db *store.Store, redisClient *redis.Client, dtoParser libparsers.TxnDtoTransformer) (*Controller, error) {
	return &Controller{
		logger:      logger,
		dbStore:     db,
		redisClient: redisClient,
		dtoParser:   dtoParser,
	}, nil
}

func (c *Controller) SetTxnPublisher(txnStatusPub txncomm.TxnStatusPublisher) {
	c.txnStatusPub = txnStatusPub
}

// HandleMessage
// this is where the meat of the handling of transaction status updates is done
// TODO: we should really have a thread pool to handle this so it won't hold the rabbitmq consumption thread back
func (c *Controller) HandleMessage(cmdResp *cmddto.CmdResponse) {
	// TODO: handle the transaction responses here
	curLogger := c.logger.AddField("func", "HandleMessage")

	// get the pl serializer
	cls, err := c.dtoParser.GetCmdSerializer(cmdResp.Command())
	if err != nil {
		if cls != nil { // something is really wrong, if it's nil that means it's not a txn response msg and we can ignore
			curLogger.Errorf(err, "issue getting parser: %s", err.Error())
		}
		return
	}

	pl, err := cls.DeserializeLegResponse(cmdResp.CmdPayload())
	if err != nil {
		curLogger.Errorf(err, "error deserializing txn response: %s", err.Error())
		return
	}

	// making sure this is a status update
	// yes, type detection... and yes, single switch case but I didn't want to use reflection here, this is simpler
	switch pl.(type) {
	case *txndto.StateUpdate:
		txnStatus := pl.(*txndto.StateUpdate)
		c.sendExternalStatusCallback(txnStatus, curLogger)
		return
	default:
		errMsg := fmt.Errorf("HandleMessage: wrong type: %T, expecting *txndto.StateUpdate", pl)
		c.logger.Errorf(errMsg, errMsg.Error())
		return
	}

}

func (c *Controller) sendExternalStatusCallback(txnStatus *txndto.StateUpdate, curLogger logger.Logger) {
	// check status is in final state
	if !txnStatus.State.IsFinal() {
		curLogger.Infof("Received non final txn %s, status %s, no op", txnStatus.TxnId(), txnStatus.Status())
		return // mo op
	}

	// get account id of the transaction itself
	// look up web hook if any
	record, err2 := entity.GetWebhookRecordForAccount(txnStatus.AccountId)
	if record == nil {
		if err2 != nil { // record not found or error with database.. how to differentiate?
			curLogger.Errorf(err2, "Error getting webhook record txn %s, accountId %s, no op", txnStatus.TxnId(), txnStatus.AccountId)
		}
		return
	}

	// generate the body
	whBodySer, err2 := c.generateCallbackBody(txnStatus)
	if err2 != nil {
		curLogger.Errorf(err2, "Error generating callback body: %s", err2.Error())
		return
	}

	// create signature from web hook
	sig, err2 := signature.ComputeSimpleSignature(whBodySer, record.Secret())
	if err2 != nil { // record not found or error with database.. how to differentiate?
		curLogger.Errorf(err2, "Error creating signature txn %s, %s, %s", txnStatus.TxnId(), whBodySer, record.Secret())
		return
	}

	// send updated status

	err2 = c.sendCallback(txnStatus.TxnId(), whBodySer, sig, record)
	if err2 != nil {
		curLogger.Errorf(err2, "Error sending callback: %s", err2.Error())
		return
	}
	return
}

func (c *Controller) generateCallbackBody(txnStatus *txndto.StateUpdate) (whBodySer []byte, err error) {
	whBody := extdto.WebhookCallbackResp{
		TxnId:     txnStatus.TxnId().String(),
		ExtRefId:  txnStatus.ExtRefId(),
		TxnState:  txnStatus.State,
		Msg:       txnStatus.Message(),
		UpdatedAt: txnStatus.CreateAtMs,
	}
	whBodySer, err = json.Marshal(whBody)
	if err != nil { // record not found or error with database.. how to differentiate?
		return nil, fmt.Errorf("error serializing response txn %s, %v, %s", txnStatus.TxnId(), whBody, err.Error())
	}
	return whBodySer, err
}

func (c *Controller) sendCallback(txnId uuid.UUID, whBodySer []byte, sig string, record entity.WebhookRecord) (err error) {
	client := &http.Client{}
	bodyReader := bytes.NewReader(whBodySer)
	req, err := http.NewRequest(http.MethodPut, record.UrlStr(), bodyReader)

	if err != nil {
		return fmt.Errorf("error creating request txn %s, %s: %s", txnId, record.UrlStr(), err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Neutronpay-Signature", sig)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request txn %s, %s: %s", txnId, record.UrlStr(), err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// nothing really consequential but helps remove linting
		}
	}(res.Body)

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response txn %s, %v: %s", txnId, res, err.Error())
	}

	if res.StatusCode != http.StatusOK {
		// doing a hard convert to string for raw instead of handing it off to the logger because who knows, maybe
		// logrus has the same bug as java log4j in some future
		return fmt.Errorf("error sending request txn %s, %s, %s: %s", txnId, record.UrlStr(), string(raw), err.Error())
	}

	return

}

/*
  PUT HANDLER LOGIC HERE FOR THE SERVICE FOR EACH CMD TO BE HANDLED
*/
