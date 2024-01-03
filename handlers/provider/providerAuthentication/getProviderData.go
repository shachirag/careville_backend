package providerAuthenticate

import (
	"careville_backend/database"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Fetch provider By ID
// @Description Fetch provider By ID
// @Tags provider authorization
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "provider ID"
// @Produce json
// @Success 200 {object} providerAuth.GetProviderResDto
// @Router /provider/get-provider-info/{id} [get]
func FetchProviderById(c *fiber.Ctx) error {

	var provider entity.ProviderEntity

	// Get the user ID from the URL parameter
	userId := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return c.Status(400).JSON(providerAuth.GetProviderResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})

	}

	providerColl := database.GetCollection("provider")

	err = providerColl.FindOne(ctx, bson.M{"_id": objId}).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(providerAuth.GetProviderResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.GetProviderResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	providerRes := providerAuth.ProviderData{
		Id:                provider.Id,
		Name:              provider.Name,
		Email:             provider.Email,
		Image:             provider.Image,
		AdditionalDetails: provider.AdditionalDetails,
		Address:           provider.Address,
		PhoneNumber:       providerAuth.PhoneNumber(provider.PhoneNumber),
		CreatedAt:         provider.CreatedAt,
		UpdatedAt:         provider.UpdatedAt,
		Notification:      provider.Notification,
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.GetProviderResDto{
		Status:   true,
		Message:  "provider data retrieved successfully",
		Provider: providerRes,
	})
}
