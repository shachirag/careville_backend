package adminAuth

import (
	"careville_backend/database"
	"careville_backend/dto/admin/adminAuth"
	"careville_backend/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Update Admin
// @Description Update Admin
// @Tags admin authorization
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "Admin ID"
// @Param admin formData adminAuth.UpdateAdminReqDto true "Update data of admin"
// @Param newAdminImage formData file false "admin profile image"
// @Produce json
// @Success 200 {object} adminAuth.UpdateAdminResDto
// @Router /admin/update-admin-info/{adminId} [put]
func UpdateAdmin(c *fiber.Ctx) error {
	var (
		adminColl = database.GetCollection("admin")
		data      adminAuth.UpdateAdminReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(adminAuth.UpdateAdminResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Check if the admin ID is provided in the request
	adminID := c.Params("adminId")
	if adminID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(adminAuth.UpdateAdminResDto{
			Status:  false,
			Message: "Admin ID is missing in the request",
		})
	}

	// Find the admin document in MongoDB based on the provided admin ID
	objID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(adminAuth.UpdateAdminResDto{
			Status:  false,
			Message: "Invalid Admin ID",
		})
	}

	// Find the admin document in MongoDB
	filter := bson.M{"_id": objID}
	result := adminColl.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(adminAuth.UpdateAdminResDto{
				Status:  false,
				Message: "Admin not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(adminAuth.UpdateAdminResDto{
			Status:  false,
			Message: "internal server error " + err.Error(),
		})
	}

	formFile, err := c.FormFile("newAdminImage")
	var imageURL string
	if err != nil {
		imageURL = data.OldImage
	} else {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(adminAuth.UpdateAdminResDto{
				Status:  false,
				Message: "Failed to open image file: " + err.Error(),
			})
		}
		defer file.Close()

		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("admin/%v-profilepic%s", id.Hex(), formFile.Filename)

		imageURL, err = utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(adminAuth.UpdateAdminResDto{
				Status:  false,
				Message: "Failed to upload image to S3: " + err.Error(),
			})
		}
	}

	// Update the admin document with new data
	update := bson.M{
		"$set": bson.M{
			"firstName": data.FirstName,
			"lastName":  data.LastName,
			"image":     imageURL,
			"updatedAt": time.Now().UTC(),
		},
	}

	// Execute the update operation
	updateRes, err := adminColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(adminAuth.UpdateAdminResDto{
			Status:  false,
			Message: "Failed to update admin data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(adminAuth.UpdateAdminResDto{
			Status:  false,
			Message: "Admin not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(adminAuth.UpdateAdminResDto{
		Status:  true,
		Message: "Admin data updated successfully",
	})
}
