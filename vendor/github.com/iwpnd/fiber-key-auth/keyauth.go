package keyauth

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func getKeysInEnv() []string {
	var keys []string

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		value := pair[1]

		if strings.HasPrefix(key, "API_KEY_") {
			keys = append(keys, value)
		}
	}

	return keys
}

func keyInKeys(key string, keys []string) bool {
	for _, k := range keys {
		if key == k {
			return true
		}
	}

	return false
}

func keyAuth(c *fiber.Ctx) error {
	key := c.Get("x-api-key")
	if key == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "no api key")
	}

	keys := getKeysInEnv()

	if keyInKeys(key, keys) {
		return c.Next()
	}

	return fiber.NewError(fiber.StatusUnauthorized, "invalid api key")
}

// New exports a keyauth middleware handler
func New() fiber.Handler {
	return keyAuth
}
