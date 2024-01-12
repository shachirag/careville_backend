package services

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"careville_backend/utils"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Update Doctor Profile Image
// @Description Update Doctor Profile Image
// @Tags services
// @Accept multipart/form-data
// @Param serviceId path string true "service ID"
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider body services.UpdateDoctorImageReqDto true "Update data of doctor image"
// @Param image formData file false "profile image"
// @Produce json
// @Success 200 {object} services.UpdateDoctorImageResDto
// @Router /provider/services/update-profile-image/{doctorId} [put]
func UpdateDoctorImage(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdateDoctorImageReqDto
		providers   entity.ServiceEntity
	)

	ctx := context.TODO()

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	doctorId := c.Params("doctorId")
	doctorObjID, err := primitive.ObjectIDFromHex(doctorId)

	if err != nil {
		return c.Status(400).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	filter := bson.M{
		"_id": providerData.ProviderId,
		"hospClinic.doctor": bson.M{
			"$elemMatch": bson.M{
				"id": doctorObjID,
			},
		},
	}

	err = serviceColl.FindOne(ctx, filter).Decode(&providers)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorImageResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	formFile, err := c.FormFile("image")
	var imageURL string
	if err != nil {
		imageURL = data.OldImage
	} else {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorImageResDto{
				Status:  false,
				Message: "Failed to open image file: " + err.Error(),
			})
		}
		defer file.Close()

		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("doctor/%v-profilepic.jpg", id.Hex())

		imageURL, err = utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorImageResDto{
				Status:  false,
				Message: "Failed to upload image to S3: " + err.Error(),
			})
		}
	}

	update := bson.M{
		"$set": bson.M{
			"hospClinic.doctor.$.image": imageURL,
		},
	}

	updateRes, err := serviceColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "Failed to update doctor image in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "doctor not found",
		})
	}

	response := services.UpdateDoctorImageResDto{
		Status:  true,
		Message: "Successfully updated doctor image",
		Image:   imageURL,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
