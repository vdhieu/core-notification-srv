package webhook

import (
	"context"
	"fmt"

	"github.com/Neutronpay/core-notification-srv/dto/webhook"
	"github.com/Neutronpay/core-notification-srv/model"
	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func New(db *gorm.DB) Store {
	return &store{
		db: db,
	}
}

func (s store) Create(ctx context.Context, webhookModel model.WebhookInfo, WebhookSecret string) (model.WebhookInfo, error) {
	return webhookModel, s.db.WithContext(ctx).Table("webhook_infos").Exec(`
		INSERT INTO webhook_infos (account_id, created_at, callback_url, secret) VALUES (
			@account_id,
			@created_at,
			@callback_url,
			pgp_sym_encrypt(@secret, @secret_key)
		);`, map[string]interface{}{
		"account_id":   webhookModel.AccountID,
		"created_at":   webhookModel.CreatedAt,
		"callback_url": webhookModel.CallbackURL,
		"secret_key":   WebhookSecret,
		"secret":       webhookModel.Secret,
	}).First(&webhookModel).Error
}

func (s store) Update(ctx context.Context, webhookModel model.WebhookInfo, WebhookSecret string) error {
	return s.db.WithContext(ctx).Exec(`
		UPDATE
			webhook_infos
	 	SET	
			callback_url = @callback_url,
			secret = pgp_sym_encrypt(@secret, @secret_key)
		WHERE
			id = @id;
		`, map[string]interface{}{
		"id":           webhookModel.ID,
		"callback_url": webhookModel.CallbackURL,
		"secret_key":   WebhookSecret,
		"secret":       webhookModel.Secret,
	}).Error
}

func (s store) Delete(ctx context.Context, webhook *model.WebhookInfo) error {
	return s.db.Table("webhook_infos").Delete(webhook).Error
}

func (s store) GetAllByAccountID(ctx context.Context, accountID string, secret string) (webhooks []*webhook.GetWebhookRes, err error) {
	db := s.db.WithContext(ctx).Table("webhook_infos").
		Select("id", fmt.Sprintf("PGP_SYM_DECRYPT(secret::bytea, '%s') as secret", secret), "created_at", "callback_url").
		Where("account_id = ?", accountID)

	return webhooks, db.Find(&webhooks).Error
}

func (s store) GetOneByID(ctx context.Context, id string) (webhook *model.WebhookInfo, err error) {
	return webhook, s.db.WithContext(ctx).Table("webhook_infos").Where("id = ?", id).First(&webhook).Error
}
