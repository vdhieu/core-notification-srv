package controller

import (
	"github.com/Neutronpay/lib-go-common/config"
	libparsers "github.com/Neutronpay/lib-go-common/dto/parsers"
	"github.com/Neutronpay/lib-go-common/logger"
	"github.com/Neutronpay/lib-go-common/route"
	"github.com/gin-gonic/gin"
)

type TxnHttpRouteHandler interface {
	route.RoutesHandler
}

// NewTxnHttpRouteHandler
// This creates a http handler to... handle the http requests.
func NewHttpRouteHandler(c config.BaseConfig, logger logger.Logger, dtoParser libparsers.TxnDtoTransformer) (handler TxnHttpRouteHandler) {
	handler = &txnHttpRouteHandler{
		logger:    logger,
		rootPath:  c.RootPath,
		dtoParser: dtoParser,
	}
	return
}

type txnHttpRouteHandler struct {
	logger    logger.Logger
	rootPath  string
	dtoParser libparsers.TxnDtoTransformer
}

// RoutesHandler interface

// GetBasePath
// the base path for
func (h *txnHttpRouteHandler) GetBasePath() string {
	return h.rootPath
}

// GetHealthChkFunc the handler for health check
func (h *txnHttpRouteHandler) GetHealthChkFunc() gin.HandlerFunc {
	return h.HealthHandler
}

// HealthHandler can make this generic
func (h *txnHttpRouteHandler) HealthHandler(ctx *gin.Context) {

}

// GetRouteDefs relative path to handler funcs
func (h *txnHttpRouteHandler) GetRouteDefs() map[string]route.RouteDef {

	// replace the below and the referred functions with what you want.  he idea is these functions
	// should be shared iwth those in the cmd_request_handler as much as possible
	routes := map[string]route.RouteDef{
		/*	"/": {
				HttpVerb: http.MethodPut,
				Secured:  false,
				HandlerF: h.NewTransactionRequestHandler,
			},

			// confirm transaction id
			"/:txn_id/confirm": {
				HttpVerb: http.MethodPut,
				Secured:  true,
				HandlerF: h.UserConfirmTransactionHandler,
			},
			"/:txnid": {
				HttpVerb: http.MethodGet,
				Secured:  false,
				HandlerF: h.GetTransactionDetailHandler,
			}, */
	}

	return routes
}

/*  EXAMPLES

// NewTransactionRequestHandler Creates a new transaction from the given request context, which is expected to
// contain a txndto.TxnRequestExternal object as its json body
func (h *txnHttpRouteHandler) NewTransactionRequestHandler(ctx *gin.Context) {

	// do auth stuff here
	jwtPayload, err := jwt.GetPayloadFromContext(ctx)
	if err != nil {
		h.logger.AddField("ctx", ctx.Request.Context()).Error(err, "UserConfirmTransactionHandler.GetPayloadFromContext")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no auth info provided"})
		return
	}

	reqId := ctx.Value(enum.RequestID)
	if nil == reqId {
		reqId = uuid.New()
	}

	// now processing the request

	var txnExt extdto.TxnRequestExternal
	err = ctx.BindJSON(&txnExt)

	h.logger.Infof("New transaction request received: %v", txnExt)

}

func (h *txnHttpRouteHandler) UserConfirmTransactionHandler(ctx *gin.Context) {
	jwtPayload, err := jwt.GetPayloadFromContext(ctx)
	if err != nil {
		h.logger.AddField("ctx", ctx.Request.Context()).Error(err, "UserConfirmTransactionHandler.GetPayloadFromContext")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no auth info provided"})
		return
	}

	reqId, err := uuid.Parse(ctx.GetString(enum.RequestID))
	if err != nil {
		reqId = uuid.New()
	}

	log := h.logger.Fields(logger.Fields{
		"func":       "UserConfirmTransactionHandler",
		"req_id":     reqId,
		"account_id": jwtPayload.AccountID,
	})

}

func (h *txnHttpRouteHandler) GetTransactionDetailHandler(ctx *gin.Context) {
	// get account id from request context
	accountCtx, err := jwt.GetPayloadFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	fmt.Println(accountCtx.AccountID)

	txnIdStr := ctx.Param("txnid")
	txnId, err := uuid.Parse(txnIdStr)
	if err != nil {
		h.logger.Warnf("GetTransactionDetailHandler Error while converting given txnid to uuid: %s Error: %s", txnIdStr, err.Error())
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error:": "Transaction not found"})
		return
	}

	txn, err := h.oms.GetTransaction(txnId)
	if err != nil {
		h.logger.Warnf("GetTransactionDetailHandler transaction for given id %s txnid not found", txnIdStr)
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error:": "Transaction not found"})
		return
	}

	// now we can respond with the now complete object
	// translate this back to external request
	reply, err := h.dtoParser.TranslateTxnEntityToExternal(txn)
	if err != nil {
		h.logger.Errorf(err, "GetTransactionDetailHandler Error while translating internal transaction back to external: %s %v", err.Error(), txn)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Transaction creation processed, please contact support if cannot retrieve information"})
		return
	}

	ctx.JSON(http.StatusOK, reply)
}
*/
