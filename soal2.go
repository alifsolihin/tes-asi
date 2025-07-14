package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

var ctx = context.Background()

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	RealName string `json:"realname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func hashSHA1(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func loginHandler(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		key := "login_" + req.Username
		raw, err := redisClient.Get(ctx, key).Result()

		if err == redis.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		var user User
		if err := json.Unmarshal([]byte(raw), &user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user data",
			})
		}

		if hashSHA1(req.Password) != user.Password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		return c.JSON(fiber.Map{
			"message":  "Login successful",
			"username": req.Username,
			"name":     user.RealName,
			"email":    user.Email,
		})
	}
}

func main() {
	app := fiber.New()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	app.Post("/login", loginHandler(redisClient))

	log.Fatal(app.Listen(":3000"))
}

