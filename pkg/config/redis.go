package config

import (
	"fmt"
	fiberRedis "github.com/gofiber/storage/redis/v3"
	"github.com/spf13/viper"
)

func InitFiberStorage() *fiberRedis.Storage {
	redisInstance := fiberRedis.New(fiberRedis.Config{
		Host:       viper.GetString("redis.host"),
		Port:       viper.GetInt("redis.port"),
		Username:   viper.GetString("redis.user"),
		Password:   viper.GetString("redis.pass"),
		Database:   viper.GetInt("redis.database"),
		ClientName: fmt.Sprintf("XYZ-%s", viper.GetString("app_env")),
	})
	return redisInstance
}

func GetRedisKey(format string, value ...any) string {
	k := "xyz"
	k += ":" + fmt.Sprintf(format, value...)
	return k
}
