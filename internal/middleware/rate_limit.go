package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
	"xyz/pkg/config"
	"xyz/pkg/helper"
	"xyz/pkg/response"
)

func RateLimit(fiberStorage fiber.Storage) func(c *fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:        500,
		Expiration: 15 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			key := config.GetRedisKey("rate_limiter:ip:%s", helper.GetIP(c))
			return key
		},
		Storage:           fiberStorage,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(ctx *fiber.Ctx) error {
			return response.ErrorRateLimit()
		},
		Next: func(c *fiber.Ctx) bool {
			// bypass if production
			/*ip := helper.GetIP(c)
			appEnv := viper.GetString("app_env")
			if appEnv != "production" && ip == "127.0.0.1" {
				return true
			}*/
			return false
		},
	})
}
