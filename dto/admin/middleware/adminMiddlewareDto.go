package customerMiddleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminMiddlwareDto struct {
	AdminId primitive.ObjectID
	Email   string
	Role    string
}

func SetAdminMiddlewareData(c *fiber.Ctx) (*AdminMiddlwareDto, error) {
	admin := c.Locals("user").(*jwt.Token)
	claims := admin.Claims.(jwt.MapClaims)

	aId := claims["Id"].(string)
	adminId, err := primitive.ObjectIDFromHex(aId)
	if err != nil {
		return nil, err
	}

	role := claims["role"].(string)
	email := claims["email"].(string)

	adminMiddlwareDto := AdminMiddlwareDto{
		AdminId: adminId,
		Email:   email,
		Role:    role,
	}
	c.Locals("AdminMiddlwareDto", adminMiddlwareDto)
	return &adminMiddlwareDto, nil
}

func GetAdminMiddlewareData(c *fiber.Ctx) AdminMiddlwareDto {
	return c.Locals("AdminMiddlwareDto").(AdminMiddlwareDto)
}
