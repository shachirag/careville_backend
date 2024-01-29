package hospClinic

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"careville_backend/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Update doctor info
// @Description Update doctor info
// @Tags hospClinic
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param doctorId path string true "doctor ID"
// @Param image formData file false "profile image"
// @Param provider body services.UpdateDoctorReqDto true "Update data of doctor"
// @Produce json
// @Success 200 {object} services.UpdateDoctorResDto
// @Router /provider/services/update-doctor-info/{doctorId} [put]
func UpdateDoctorInfo(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdateDoctorReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateDoctorResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Get provider data from middleware
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

	err = serviceColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorResDto{
				Status:  false,
				Message: "doctor not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorResDto{
			Status:  false,
			Message: "Failed to fetch doctor from MongoDB: " + err.Error(),
		})
	}

	formFile, err := c.FormFile("image")
	var imageURL string
	if err != nil {
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
			"hospClinic.doctor.$.speciality": data.Speciality,
			"hospClinic.doctor.$.name":       data.Name,
			"hospClinic.doctor.$.schedule":   bson.A{},
			"hospClinic.information.$.image": imageURL,
			"updatedAt":                      time.Now().UTC(),
		},
	}

	for _, schedule := range data.Schedule {
		scheduleUpdate := bson.M{
			"startTime": schedule.StartTime,
			"endTime":   schedule.EndTime,
			"days":      schedule.Days,
		}
		update["$set"].(bson.M)["hospClinic.doctor.$.schedule"] = append(update["$set"].(bson.M)["hospClinic.doctor.$.schedule"].(bson.A), scheduleUpdate)
	}

	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorResDto{
			Status:  false,
			Message: "Failed to update doctor data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorResDto{
			Status:  false,
			Message: "doctor not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdateDoctorResDto{
		Status:  true,
		Message: "doctor data updated successfully",
	})
}
