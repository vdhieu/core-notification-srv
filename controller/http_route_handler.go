package controller

import (
	"net/http"
	"time"

	"github.com/Neutronpay/core-notification-srv/config"
	"github.com/Neutronpay/core-notification-srv/dto/webhook"
	"github.com/Neutronpay/core-notification-srv/model"
	"github.com/Neutronpay/core-notification-srv/store"
	libparsers "github.com/Neutronpay/lib-go-common/dto/parsers"
	"github.com/Neutronpay/lib-go-common/jwt"
	"github.com/Neutronpay/lib-go-common/logger"
	"github.com/Neutronpay/lib-go-common/route"
	"github.com/gin-gonic/gin"
)

type TxnHttpRouteHandler interface {
	route.RoutesHandler
}

// NewTxnHttpRouteHandler
// This creates a http handler to... handle the http requests.
func NewHttpRouteHandler(c config.Config, store store.Store, logger logger.Logger, dtoParser libparsers.TxnDtoTransformer) (handler TxnHttpRouteHandler) {
	handler = &notiHttpRouteHandler{
		logger:    logger,
		cfg:       c,
		store:     store,
		rootPath:  c.Base.RootPath,
		dtoParser: dtoParser,
	}
	return
}

type notiHttpRouteHandler struct {
	store     store.Store
	cfg       config.Config
	logger    logger.Logger
	rootPath  string
	dtoParser libparsers.TxnDtoTransformer
}

// RoutesHandler interface

// GetBasePath
// the base path for
func (h *notiHttpRouteHandler) GetBasePath() string {
	return h.rootPath
}

// GetHealthChkFunc the handler for health check
func (h *notiHttpRouteHandler) GetHealthChkFunc() gin.HandlerFunc {
	return h.HealthHandler
}

// HealthHandler can make this generic
func (h *notiHttpRouteHandler) HealthHandler(ctx *gin.Context) {

}

// GetRouteDefs relative path to handler funcs
func (h *notiHttpRouteHandler) GetRouteDefs() []route.RouteDef {
	routes := []route.RouteDef{
		{
			Path:     "/webhook",
			HttpVerb: http.MethodPost,
			Secured:  true,
			HandlerF: h.CreateWebHookHandler,
		},
		{
			Path:     "/webhook",
			HttpVerb: http.MethodGet,
			Secured:  true,
			HandlerF: h.GetWebHookHandler,
		},
		{
			Path:     "/webhook/:id",
			HttpVerb: http.MethodPut,
			Secured:  true,
			HandlerF: h.UpdateWebHookHandler,
		},
		{
			Path:     "/webhook/:id",
			HttpVerb: http.MethodDelete,
			Secured:  true,
			HandlerF: h.DeleteWebHookHandler,
		},
	}

	return routes
}

func (h *notiHttpRouteHandler) CreateWebHookHandler(ctx *gin.Context) {
	log := h.logger.Fields(logger.Fields{
		"func": "notiHttpRouteHandler.CreateWebHookHandler",
	})

	jwtPayload, err := jwt.GetPayloadFromContext(ctx)
	if err != nil {
		log.Error(err, "jwt.GetPayloadFromContext")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no auth info provided"})
		return
	}

	var webhookModel webhook.CreateWebHookReq
	if err = ctx.ShouldBindJSON(&webhookModel); err != nil {
		log.Error(err, "notiHttpRouteHandler.CreateWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
		return
	}

	newWebhook := model.WebhookInfo{
		AccountID:   jwtPayload.AccountID,
		CallbackURL: webhookModel.CallbackURL,
		Secret:      webhookModel.Secret,
		CreatedAt:   time.Now(),
	}
	res, err := h.store.Webhook().Create(ctx, newWebhook, h.cfg.WebhookSecret)
	if err != nil {
		log.Error(err, "notiHttpRouteHandler.CreateWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create webhook"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *notiHttpRouteHandler) GetWebHookHandler(ctx *gin.Context) {
	log := h.logger.Fields(logger.Fields{
		"func": "notiHttpRouteHandler.GetWebHookHandler",
	})
	jwtPayload, err := jwt.GetPayloadFromContext(ctx)
	if err != nil {
		log.Error(err, "notiHttpRouteHandler.GetWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no auth info provided"})
		return
	}

	res, err := h.store.Webhook().GetAllByAccountID(ctx, jwtPayload.AccountID, h.cfg.WebhookSecret)
	if err != nil {
		log.Error(err, "notiHttpRouteHandler.GetWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get webhooks"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *notiHttpRouteHandler) UpdateWebHookHandler(ctx *gin.Context) {
	log := h.logger.Fields(logger.Fields{
		"func": "notiHttpRouteHandler.UpdateWebHookHandler",
	})

	jwtPayload, err := jwt.GetPayloadFromContext(ctx)
	if err != nil {
		log.Error(err, "notiHttpRouteHandler.UpdateWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no auth info provided"})
		return
	}

	webhookID := ctx.Param("id")
	if webhookID == "" {
		log.Error(err, "notiHttpRouteHandler.UpdateWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no webhook id provided"})
		return
	}

	var webhookModel webhook.UpdateWebHookReq
	if err = ctx.ShouldBindJSON(&webhookModel); err != nil {
		log.Error(err, "notiHttpRouteHandler.UpdateWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
		return
	}

	// get webhook info from db
	webhook, err := h.store.Webhook().GetOneByID(ctx, webhookID)
	if err != nil {
		log.Error(err, "notiHttpRouteHandler.UpdateWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get webhook"})
		return
	}

	// check if the webhook belongs to the user
	if webhook.AccountID != jwtPayload.AccountID {
		log.Error(err, "wrong account id")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// update webhook info
	webhook.CallbackURL = webhookModel.CallbackURL
	webhook.Secret = webhookModel.Secret

	if err = h.store.Webhook().Update(ctx, *webhook, h.cfg.WebhookSecret); err != nil {
		log.Error(err, "notiHttpRouteHandler.UpdateWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update webhook"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (h *notiHttpRouteHandler) DeleteWebHookHandler(ctx *gin.Context) {
	log := h.logger.Fields(logger.Fields{
		"func": "notiHttpRouteHandler.UpdateWebHookHandler",
	})

	jwtPayload, err := jwt.GetPayloadFromContext(ctx)
	if err != nil {
		log.Error(err, "notiHttpRouteHandler.DeleteWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no auth info provided"})
		return
	}

	webhookID := ctx.Param("id")
	if webhookID == "" {
		log.Error(err, "notiHttpRouteHandler.DeleteWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no webhook id provided"})
		return
	}

	// get webhook info from db
	webhook, err := h.store.Webhook().GetOneByID(ctx, webhookID)
	if err != nil {
		log.Error(err, "notiHttpRouteHandler.DeleteWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get webhook"})
		return
	}

	// check if the webhook belongs to the user
	if webhook.AccountID != jwtPayload.AccountID {
		log.Error(err, "notiHttpRouteHandler.DeleteWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err = h.store.Webhook().Delete(ctx, webhook); err != nil {
		log.Error(err, "notiHttpRouteHandler.DeleteWebHookHandler")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to delete webhook"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
