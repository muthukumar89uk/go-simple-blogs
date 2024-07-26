package middleware

import (
	//User-defined packages
	"blog/helper"
	"blog/logs"
	"blog/models"
	"blog/repository"

	//Inbuild packages
	"errors"
	"os"
	"strconv"
	"time"

	//Third-party packages
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Create a JWT token with the needed claims
func CreateToken(user models.User, c *fiber.Ctx) (string, error) {
	log := logs.Log()

	if err := helper.Config(".env"); err != nil {
		log.Error.Println("Error : 'Error at loading '.env' file'")
	}

	exp := time.Now().Add(time.Hour * 24).Unix()
	userId := strconv.Itoa(int(user.UserId))
	roleId := strconv.Itoa(int(user.RoleId))

	claims := jwt.StandardClaims{
		ExpiresAt: exp,
		Id:        userId,
		IssuedAt:  time.Now().Unix(),
		Subject:   roleId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Token and claims validation
func AuthMiddleware(db *gorm.DB) fiber.Handler {
	log := logs.Log()

	if err := helper.Config(".env"); err != nil {
		log.Error.Println("Error : 'Error at loading '.env' file'")
	}
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		//To check the token is empty or not
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Error": "token is empty",
			})
		}

		for index, char := range tokenString {
			if char == ' ' {
				tokenString = tokenString[index+1:]
			}
		}

		claims := jwt.StandardClaims{}

		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"Error": "Invalid token signature",
				})
			} else if claims.ExpiresAt < time.Now().Unix() {
				repository.DeleteToken(db, claims.Id)

				log.Error.Println("Error : 'session expired...login again!!!' Status : 440")
				c.Status(fiber.StatusGatewayTimeout)

				return c.JSON(fiber.Map{
					"status": 440,
					"Error":  "session expired...login again!!!",
				})
			} else if !token.Valid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"Error": "Invalid token",
				})
			}
		}

		// Check the user's role
		if claims.Subject == "1" {
			c.Locals("role", "admin")
		} else if claims.Subject == "2" {
			c.Locals("role", "user")
		}

		return c.Next()
	}
}

// Get a claims from the token
func GetTokenClaims(c *fiber.Ctx) jwt.StandardClaims {
	log := logs.Log()

	if err := helper.Config(".env"); err != nil {
		log.Error.Println("Error : 'Error at loading '.env' file'")
	}

	tokenString := c.Get("Authorization")

	for index, char := range tokenString {
		if char == ' ' {
			tokenString = tokenString[index+1:]
		}
	}

	claims := jwt.StandardClaims{}

	jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	return claims
}

// Admin authorization
func AdminAuth(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "admin" {
		return errors.New("unauthorized entry")
	}

	return nil
}

// User authorization
func UserAuth(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "user" {
		return errors.New("unauthorized entry")
	}
	
	return nil
}
