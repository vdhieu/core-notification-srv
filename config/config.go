package config

import (
	"fmt"

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

var config *viper.Viper
var c Config

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) *Config {
	log := logger.NewLogrusLogger("init", env)
	var err error
	config = viper.New()
	config.SetConfigType("env")
	config.SetConfigName(fmt.Sprintf(env))
	config.AddConfigPath("env")

	config.AutomaticEnv()

	err = config.ReadInConfig()
	if err != nil {
		log.Errorf(err, "error on parsing configuration file, %s", err.Error())
	}
	err = config.Unmarshal(&c)
	if err != nil {
		log.Errorf(err, "error on parsing configuration file %s", err.Error())
	}
	// TODO: thing about a robust way to handle auto env in case of configuration file is not available
	c.Base.Env = env
	c.Base.Name = config.GetString("SRV_NAME")
	c.Base.Port = config.GetString("SRV_PORT")
	c.Base.BaseURL = config.GetString("SRV_BASE_URL")
	c.Base.RootPath = config.GetString("SRV_RESTFUL_ROOT_PATH")
	c.DBConf.Host = config.GetString("DB_HOST")
	c.DBConf.Port = config.GetString("DB_PORT")
	c.DBConf.User = config.GetString("DB_USER")
	c.DBConf.Pass = config.GetString("DB_PASS")
	c.DBConf.Name = config.GetString("DB_NAME")
	c.DBConf.SSLMode = config.GetString("DB_SSL_MODE")
	c.RmqConf.Host = config.GetString("PF_RMQ_HOST")
	c.RmqConf.Port = config.GetString("PF_RMQ_PORT")
	c.RmqConf.User = config.GetString("PF_RMQ_USER")
	c.RmqConf.Pass = config.GetString("PF_RMQ_PASS")
	c.RedisConf.Host = config.GetString("REDIS_HOST")
	c.RedisConf.Port = config.GetString("REDIS_PORT")
	c.RedisConf.Pass = config.GetString("REDIS_PASS")
	c.RedisConf.DB = config.GetInt("REDIS_DB")
	c.RedisConf.SSL = config.GetBool("REDIS_SSL")
	c.JwtSecret = config.GetString("JWT_SECRET")
	c.WebhookSecret = config.GetString("WEBHOOK_SECRET")

	return &c
}

func GetConfig() *Config {
	return &c
}
