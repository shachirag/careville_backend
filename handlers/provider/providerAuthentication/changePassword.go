package providerAuthenticate

import (
	"careville_backend/database"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Change provider Password
// @Description change provider Password
// @Tags provider authorization
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "provider ID"
// @Param provider body  providerAuth.ProviderChangePasswordReqDto true "Change password of provider"
// @Produce json
// @Success 200 {object}  providerAuth.ProviderChangePasswordResDto
// @Router /provider/change-password/{id} [put]
func ChangeProviderPassword(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        providerAuth.ProviderChangePasswordReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to parse request body: " + err.Error(),
		})
	}

	// Get the customer ID from the request parameters
	providerID := c.Params("id")

	// Convert the admin ID string to primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(providerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Invalid provider ID",
		})
	}

	filter := bson.M{"_id": objID}

	result := serviceColl.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(providerAuth.ProviderChangePasswordResDto{
				Status:  false,
				Message: "provider not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Error by finding provider " + err.Error(),
		})
	}

	// Decode the customer data
	var provider entity.ServiceEntity
	err = result.Decode(&provider)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to decode provider data: " + err.Error(),
		})
	}

	if data.CurrentPassword != "" {
		// Validate the current password
		err = bcrypt.CompareHashAndPassword([]byte(provider.Password), []byte(data.CurrentPassword))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(providerAuth.ProviderChangePasswordResDto{
				Status:  false,
				Message: "Current password is incorrect",
			})
		}
	}

	// Hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 6)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to hash the new password",
		})
	}

	update := bson.M{
		"$set": bson.M{
			"password": string(hashedNewPassword),
		},
	}

	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to update provider password in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "provider not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.ProviderChangePasswordResDto{
		Status:  true,
		Message: "provider password updated successfully",
	})
}
