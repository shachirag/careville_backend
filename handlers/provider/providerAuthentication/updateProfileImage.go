package providerAuthenticate

import (
    "careville_backend/database"
    providerMiddleware "careville_backend/dto/provider/middleware"
    providerAuth "careville_backend/dto/provider/providerAuth"
    "careville_backend/entity"
    "careville_backend/utils"
    "context"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "time"
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
            subfield:   profileImage,
            "updatedAt": time.Now().UTC(),
        },
    }

    // Update the document in the database
    opts := options.Update().SetUpsert(true)
    updateRes, err := serviceColl.UpdateOne(ctx, filter, update, opts)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderImageResDto{
            Status:  false,
            Message: "Failed to update provider data in MongoDB: " + err.Error(),
        })
    }

    if updateRes.MatchedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(providerAuth.UpdateProviderImageResDto{
            Status:  false,
            Message: "provider not found",
        })
    }

    response := providerAuth.UpdateProviderImageResDto{
        Status:  true,
        Message: "Successfully updated image",
        Image:   profileImage,
    }

    return c.Status(fiber.StatusOK).JSON(response)
}
