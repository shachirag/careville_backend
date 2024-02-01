package adminAuth

import (
	"careville_backend/database"
	adminMiddleware "careville_backend/dto/admin/middleware"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Change adminPassword is the handler for changing provider passwords
// @Summary Change provider Password
// @Description Change provider Password
// @Tags provider authorization
// @Accept application/json
// @Param Authorization header string true "Authentication header"
// @Param provider body providerAuth.ProviderChangePasswordReqDto true "Change password of provider"
// @Produce json
// @Success 200 {object} providerAuth.ProviderChangePasswordResDto
// @Router /admin/profile/change-password [put]
func ChangeAdminPassword(c *fiber.Ctx) error {
	var (
		adminColl = database.GetCollection("admin")
		data      providerAuth.ProviderChangePasswordReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to parse request body: " + err.Error(),
		})
	}

	adminData := adminMiddleware.GetAdminMiddlewareData(c)

	filter := bson.M{"_id": adminData.AdminId}

	result := adminColl.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(providerAuth.ProviderChangePasswordResDto{
				Status:  false,
				Message: "admin not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Error finding provider: " + err.Error(),
		})
	}

	// Decode the provider data
	var admin entity.AdminEntity
	err = result.Decode(&admin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to decode admin data: " + err.Error(),
		})
	}

	if data.CurrentPassword != "" {
		// Validate the current password
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(data.CurrentPassword))
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
	updateRes, err := adminColl.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to update admin password in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "admin not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.ProviderChangePasswordResDto{
		Status:  true,
		Message: "admin password updated successfully",
	})
}
