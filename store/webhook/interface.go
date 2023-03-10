package webhook

import (
	"context"

	"github.com/Neutronpay/core-notification-srv/dto/webhook"
	"github.com/Neutronpay/core-notification-srv/model"
)

type Store interface {
	Create(ctx context.Context, webhookModel model.WebhookInfo, WebhookSecret string) (res model.WebhookInfo, err error)
	Update(ctx context.Context, webhookModel model.WebhookInfo, WebhookSecret string) (err error)
	Delete(ctx context.Context, webhook *model.WebhookInfo) (err error)
	GetAllByAccountID(ctx context.Context, accountID string, secret string) (webhooks []*webhook.GetWebhookRes, err error)
	GetOneByID(ctx context.Context, id string) (webhook *model.WebhookInfo, err error)
}
