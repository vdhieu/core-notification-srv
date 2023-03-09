package store

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/Neutronpay/core-notification-srv/config"
	webhook "github.com/Neutronpay/core-notification-srv/store/webhook"
)

type TxFunc func(context context.Context, store Store) error

type Store interface {
	WithTransaction(ctx context.Context, txFunc TxFunc, opts ...*sql.TxOptions) error
	Webhook() webhook.Store
}

type pgStore struct {
	db            *gorm.DB
	inTransaction bool
	cfg           *config.Config
	webhookstore  webhook.Store
}

func New(db *gorm.DB, cfg *config.Config) Store {
	return &pgStore{
		db:           db,
		cfg:          cfg,
		webhookstore: webhook.New(db),
	}
}

func (s *pgStore) WithTransaction(ctx context.Context, txFunc TxFunc, opts ...*sql.TxOptions) error {
	if s.inTransaction {
		return errors.New("db txn nested in db txn")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		store := &pgStore{
			db:            tx,
			cfg:           s.cfg,
			inTransaction: true,
			webhookstore:  webhook.New(tx),
		}
		return txFunc(ctx, store)
	}, opts...)
}

func (s *pgStore) Webhook() webhook.Store {
	return s.webhookstore
}
