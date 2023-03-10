package store

import (
	"database/sql"
	"fmt"

	"github.com/Neutronpay/core-notification-srv/config"
	"github.com/Neutronpay/lib-go-common/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ConnDb(cfg *config.Config) *gorm.DB {
	ds := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBConf.User, cfg.DBConf.Pass,
		cfg.DBConf.Host, cfg.DBConf.Port, cfg.DBConf.Name, cfg.DBConf.SSLMode,
	)

	conn, err := sql.Open("postgres", ds)
	if err != nil {
		logger.L.Fatalf(err, "failed to open database connection")
	}

	db, err := gorm.Open(postgres.New(
		postgres.Config{Conn: conn}),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	if err != nil {
		logger.L.Fatalf(err, "failed to open database connection")
	}

	logger.L.Info("database connected")

	return db
}
