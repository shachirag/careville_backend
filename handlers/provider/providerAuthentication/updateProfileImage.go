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
// @Param provider body providerAuth.UpdateImageReqDto true "Update data of admin"
// @Param image formData file false "profile image"
// @Produce json
// @Success 200 {object} providerAuth.UpdateProviderResDto
// @Router /provider/profile/update-profile-image [put]
func UpdateImage(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
		data        providerAuth.UpdateImageReqDto
		providers   entity.ServiceEntity
	)

	ctx := context.TODO()

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}
	err = serviceColl.FindOne(ctx, filter).Decode(&providers)
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

	formFile, err := c.FormFile("image")
	var imageURL string
	if err != nil {
		imageURL = data.OldImage
	} else {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
				Status:  false,
				Message: "Failed to open image file: " + err.Error(),
			})
		}
		defer file.Close()

		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("provider/%v-profilepic.jpg", id.Hex())

		imageURL, err = utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
				Status:  false,
				Message: "Failed to upload image to S3: " + err.Error(),
			})
		}
	}

	update := bson.M{}

	if providers.Role == "healthFacility" && providers.FacilityOrProfession == "hospClinic" {
		update = bson.M{
			"$set": bson.M{
				"hospClinic": bson.M{
					"information": bson.M{
						"image": imageURL,
					},
				},
				"updatedAt": time.Now().UTC(),
			},
		}
	} else if providers.Role == "healthFacility" && providers.FacilityOrProfession == "laboratory" {
		update = bson.M{
			"$set": bson.M{
				"laboratory": bson.M{
					"information": bson.M{
						"image": imageURL,
					},
				},
				"updatedAt": time.Now().UTC(),
			},
		}
	} else if providers.Role == "healthFacility" && providers.FacilityOrProfession == "fitnessCenter" {
		update = bson.M{
			"$set": bson.M{
				"fitnessCenter": bson.M{
					"information": bson.M{
						"image": imageURL,
					},
				},
				"updatedAt": time.Now().UTC(),
			},
		}
	} else if providers.Role == "healthFacility" && providers.FacilityOrProfession == "pharmacy" {
		update = bson.M{
			"$set": bson.M{
				"pharmacy": bson.M{
					"information": bson.M{
						"image": imageURL,
					},
				},
				"updatedAt": time.Now().UTC(),
			},
		}
	} else if providers.Role == "healthProfessional" && providers.FacilityOrProfession == "medicalLabScientist" {
		update = bson.M{
			"$set": bson.M{
				"medicalLabScientist": bson.M{
					"information": bson.M{
						"image": imageURL,
					},
				},
				"updatedAt": time.Now().UTC(),
			},
		}
	} else if providers.Role == "healthProfessional" && providers.FacilityOrProfession == "nurse" {
		update = bson.M{
			"$set": bson.M{
				"doctor": bson.M{
					"information": bson.M{
						"image": imageURL,
					},
				},
				"updatedAt": time.Now().UTC(),
			},
		}
	} else if providers.Role == "healthProfessional" && providers.FacilityOrProfession == "doctor" {
		update = bson.M{
			"$set": bson.M{
				"doctor": bson.M{
					"information": bson.M{
						"image": imageURL,
					},
				},
				"updatedAt": time.Now().UTC(),
			},
		}
	} else if providers.Role == "healthProfessional" && providers.FacilityOrProfession == "physiotherapist" {
		update = bson.M{
			"$set": bson.M{
				"physiotherapist": bson.M{
					"information": bson.M{
						"image": imageURL,
					},
				},
				"updatedAt": time.Now().UTC(),
			},
		}
	}

	updateRes, err := serviceColl.UpdateOne(ctx, filter, update)
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

	response := providerAuth.UpdateProviderResDto{
		Status:  true,
		Message: "Successfully updated image",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
