package providerMiddleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProviderMiddlwareDto struct {
	ProviderId           primitive.ObjectID
	Email                string
	Role                 string
	ServiceRole          string
	FacilityOrProfession string
}

func SetProviderMiddlewareData(c *fiber.Ctx) (*ProviderMiddlwareDto, error) {
	customer := c.Locals("user").(*jwt.Token)
	claims := customer.Claims.(jwt.MapClaims)

	pId := claims["Id"].(string)
	userId, err := primitive.ObjectIDFromHex(pId)
	if err != nil {
		return nil, err
	}

	role := claims["role"].(string)
	email := claims["email"].(string)
	serviceRole := claims["serviceRole"].(string)
	facilityOrProfesion := claims["facilityOrProfession"].(string)

	providerMiddlwareDto := ProviderMiddlwareDto{
		ProviderId:           userId,
		Email:                email,
		Role:                 role,
		ServiceRole:          serviceRole,
		FacilityOrProfession: facilityOrProfesion,
	}
	c.Locals("ProviderMiddlwareDto", providerMiddlwareDto)
	return &providerMiddlwareDto, nil
}

func GetProviderMiddlewareData(c *fiber.Ctx) ProviderMiddlwareDto {
	return c.Locals("ProviderMiddlwareDto").(ProviderMiddlwareDto)
}


