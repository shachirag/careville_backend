package providerAuthenticate

import (
	"careville_backend/database"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Update provider
// @Description Update provider
// @Tags provider authorization
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "provider ID"
// @Param provider formData providerAuth.UpdateProviderReqDto true "Update data of provider"
// @Param newProviderImage formData file false "provider profile image"
// @Produce json
// @Success 200 {object} providerAuth.UpdateProviderResDto
// @Router /provider/update-provider-data/{id} [put]
func UpdateProvider(c *fiber.Ctx) error {

	var (
		providerColl = database.GetCollection("provider")
		data         providerAuth.UpdateProviderReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	providerID := c.Params("id")
	if providerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "provider ID is missing in the request",
		})
	}

	objID, err := primitive.ObjectIDFromHex(providerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "Invalid provider ID",
		})
	}

	filter := bson.M{"_id": objID}
	result := providerColl.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(providerAuth.UpdateProviderResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "internal server error " + err.Error(),
		})
	}

	formFile, err := c.FormFile("newProviderImage")
	var imageURL string
	if err != nil {
		imageURL = data.OldProfileImage
	} else {
		// New image uploaded
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
				Status:  false,
				Message: "Failed to open image file: " + err.Error(),
			})
		}
		defer file.Close()

		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("provider/%v-profilepic%v", id.Hex(), formFile.Filename)

		imageURL, err = utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
				Status:  false,
				Message: "Failed to upload image to S3: " + err.Error(),
			})
		}
	}

	// Update the admin document with new data
	update := bson.M{
		"name": data.Name,
		"phoneNumber": bson.M{
			"dialCode": data.DialCode,
			"number":   data.PhoneNumber,
		},
		"image":             imageURL,
		"additionalDetails": data.AdditionalDetails,
		"address":           data.Address,
		"updatedAt":         time.Now().UTC(),
	}

	updateFields := bson.M{"$set": update}
	// Execute the update operation
	updateRes, err := providerColl.UpdateOne(ctx, filter, updateFields)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "provider not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.UpdateProviderResDto{
		Status:  true,
		Message: "provider data updated successfully",
	})
}
