package hospClinic

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"careville_backend/utils"
	"encoding/json"
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
// @Param provider formData services.UpdateDoctorReqDto true "Update data of doctor"
// @Produce json
// @Success 200 {object} services.UpdateDoctorResDto
// @Router /provider/services/update-doctor-info/{doctorId} [put]
func UpdateDoctorInfo(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdateDoctorReqDto
		provider    entity.ServiceEntity
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.MedicalLabScientistResDto{
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

	// Access the MultipartForm directly from the fiber.Ctx
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorResDto{
			Status:  false,
			Message: "Failed to get multipart form: " + err.Error(),
		})
	}

	// Get the file header for the "images" field from the form
	formFiles := form.File["image"]

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.UpdateDoctorResDto{
				Status:  false,
				Message: "Failed to upload image to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("doctor/%v-image-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		imageURL, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorResDto{
				Status:  false,
				Message: "Failed to upload image to S3: " + err.Error(),
			})
		}

		// Append the image URL to the Images field
		provider.HospClinic.Information.Image = imageURL
	}

	update := bson.M{
		"$set": bson.M{
			"hospClinic.doctor.$.speciality": data.Speciality,
			"hospClinic.doctor.$.name":       data.Name,
			"hospClinic.doctor.$.schedule":   bson.A{},
			"hospClinic.doctor.$.image":      provider.HospClinic.Information.Image,
			"updatedAt":                      time.Now().UTC(),
		},
	}

	// Clearing existing schedule
	// update["$set"].(bson.M)["hospClinic.doctor.$.schedule"] = bson.A{}

	for _, schedule := range data.Schedule {
		scheduleUpdate := bson.M{
			"startTime": schedule.StartTime,
			"endTime":   schedule.EndTime,
			"days":      schedule.Days,
		}
		fmt.Print(schedule)
		update["$set"].(bson.M)["hospClinic.doctor.$.schedule"] = append(update["$set"].(bson.M)["hospClinic.doctor.$.schedule"].(bson.A), scheduleUpdate)
	}

	session, err := database.GetMongoClient().StartSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorResDto{
			Status:  false,
			Message: "Failed to start session",
		})
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		updateRes, err := serviceColl.UpdateOne(sessCtx, filter, update)
		if err != nil {
			return nil, err
		}

		if updateRes.MatchedCount == 0 {
			return nil, mongo.ErrNoDocuments
		}

		appointmentUpdate := bson.M{"$set": bson.M{
			"hospital.doctor.speciality": data.Speciality,
			"hospital.doctor.name":       data.Name,
		}}

		filter := bson.M{
			"serviceId":          providerData.ProviderId,
			"hospital.doctor.id": doctorObjID,
		}

		_, err = database.GetCollection("appointment").UpdateMany(sessCtx, filter, appointmentUpdate)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorResDto{
			Status:  false,
			Message: "Failed to update appointment data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdateDoctorResDto{
		Status:  true,
		Message: "Doctor data updated successfully",
	})
}
