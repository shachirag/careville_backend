package customerMiddleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerMiddlwareDto struct {
	CustomerId primitive.ObjectID
	Email      string
	Role       string
}

func SetCustomerMiddlewareData(c *fiber.Ctx) (*CustomerMiddlwareDto, error) {
	customer := c.Locals("user").(*jwt.Token)
	claims := customer.Claims.(jwt.MapClaims)

	cId := claims["Id"].(string)
	userId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		return nil, err
	}

	role := claims["role"].(string)
	email := claims["email"].(string)

	customerMiddlwareDto := CustomerMiddlwareDto{
		CustomerId: userId,
		Email:      email,
		Role:       role,
	}
	c.Locals("CustomerMiddlwareDto", customerMiddlwareDto)
	return &customerMiddlwareDto, nil
}

func GetCustomerMiddlewareData(c *fiber.Ctx) CustomerMiddlwareDto {
	return c.Locals("CustomerMiddlwareDto").(CustomerMiddlwareDto)
}
