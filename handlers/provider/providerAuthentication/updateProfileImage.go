package providerAuthenticate

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"
	"careville_backend/utils"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Update Profile Image
// @Description Update Profile Image
// @Tags provider authorization
// @Accept multipart/form-data
// @Param id path string true "provider ID"
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param image formData file false "profile image"
// @Produce json
// @Success 200 {object} providerAuth.UpdateProviderImageResDto
// @Router /provider/profile/update-profile-image [put]
func UpdateImage(c *fiber.Ctx) error {
	var (
		serviceColl  = database.GetCollection("service")
		providers    entity.ServiceEntity
		profileImage string
	)

	ctx := context.TODO()

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}
	err := serviceColl.FindOne(ctx, filter).Decode(&providers)
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

	// Access the MultipartForm directly from the fiber.Ctx
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.GetProviderResDto{
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
			return c.Status(fiber.StatusBadRequest).JSON(providerAuth.GetProviderResDto{
				Status:  false,
				Message: "Failed to upload image to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("profile/%v-image-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		imageURL, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.GetProviderResDto{
				Status:  false,
				Message: "Failed to upload image to S3: " + err.Error(),
			})
		}

		profileImage = imageURL
	}

	var subfield string
	switch {
	case providers.Role == "healthFacility":
		switch providers.FacilityOrProfession {
		case "hospClinic":
			subfield = "hospClinic.information.image"
		case "laboratory":
			subfield = "laboratory.information.image"
		case "fitnessCenter":
			subfield = "fitnessCenter.information.image"
		case "pharmacy":
			subfield = "pharmacy.information.image"
		}
	case providers.Role == "healthProfessional":
		switch providers.FacilityOrProfession {
		case "medicalLabScientist":
			subfield = "medicalLabScientist.information.image"
		case "nurse":
			subfield = "nurse.information.image"
		case "doctor":
			subfield = "doctor.information.image"
		case "physiotherapist":
			subfield = "physiotherapist.information.image"
		}
	}

	// Update the subfield in the document
	update := bson.M{
		"$set": bson.M{
			subfield:    profileImage,
			"updatedAt": time.Now().UTC(),
		},
	}

	session, err := database.GetMongoClient().StartSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
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

		filter := bson.M{"serviceId": providerData.ProviderId}

		var appointment entity.AppointmentEntity
		err = database.GetCollection("appointment").FindOne(ctx, filter).Decode(&appointment)
		if err != nil {
			return nil, err
		}

		var appointmentUpdate bson.M
		switch appointment.Role {
		case "healthFacility":
			switch appointment.FacilityOrProfession {
			case "hospital":
				appointmentUpdate = bson.M{"$set": bson.M{
					"hospital.information.image": profileImage,
				}}
			case "laboratory":
				appointmentUpdate = bson.M{"$set": bson.M{
					"laboratory.information.image": profileImage,
				}}
			case "fitnessCenter":
				appointmentUpdate = bson.M{"$set": bson.M{
					"fitnessCenter.information.image": profileImage,
				}}
			case "pharmacy":
				appointmentUpdate = bson.M{"$set": bson.M{
					"pharmacy.information.image": profileImage,
				}}
			}
		case "healthProfessional":
			switch appointment.FacilityOrProfession {
			case "doctor":
				appointmentUpdate = bson.M{"$set": bson.M{
					"hospital.information.image": profileImage,
				}}
			case "physiotherapist":
				appointmentUpdate = bson.M{"$set": bson.M{
					"physiotherapist.information.image": profileImage,
				}}
			case "medicalLabScientist":
				appointmentUpdate = bson.M{"$set": bson.M{
					"medicalLabScientist.information.image": profileImage,
				}}
			case "nurse":
				appointmentUpdate = bson.M{"$set": bson.M{
					"nurse.information.image": profileImage,
				}}
			}
		}

		_, err = database.GetCollection("appointment").UpdateMany(sessCtx, filter, appointmentUpdate)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "Failed to update appointment data: " + err.Error(),
		})
	}

	response := providerAuth.UpdateProviderImageResDto{
		Status:  true,
		Message: "Successfully updated image",
		Image:   profileImage,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
