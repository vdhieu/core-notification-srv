package config

import (
	"sync"

	libconfig "github.com/Neutronpay/lib-go-common/config"
	logger "github.com/Neutronpay/lib-go-common/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Base          libconfig.BaseConfig   `json:"base"`
	DBConf        libconfig.DBConn       `json:"dbconf"`
	RmqConf       libconfig.RabbitMqConn `json:"rmqconf"`
	RedisConf     libconfig.RedisConn    `json:"redisconf"`
	JwtSecret     string                 `json:"jwtSecret"`
	WebhookSecret string                 `json:"webhookSecret"`
}

var cfg *viper.Viper
var c Config
var singleton sync.Once

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func GetConfig(env string) *Config {
	singleton.Do(func() {
		log := logger.NewLogrusLogger("init_config", env)
		var err error
		cfg = viper.New()
		cfg.SetConfigType("env")
		cfg.SetConfigName(env)
		cfg.AddConfigPath("env")

		cfg.AutomaticEnv()

		err = cfg.ReadInConfig()
		if err != nil {
			log.Errorf(err, "error on parsing configuration file, %s", err.Error())
		}
		err = cfg.Unmarshal(&c)
		if err != nil {
			log.Errorf(err, "error on parsing configuration file %s", err.Error())
		}

		c.Base.Env = env
		c.Base.Name = cfg.GetString("NAME")
		c.Base.Port = cfg.GetString("PORT")
		c.Base.BaseURL = cfg.GetString("BASE_URL")

		c.JwtSecret = cfg.GetString("JWT_SECRET")

		c.DBConf = libconfig.GenDBConfig(cfg)
		c.RedisConf = libconfig.GenRedisConfig(cfg)
		c.RmqConf = libconfig.GenRabbitMqConfig(cfg)
	})

	return &c
}
