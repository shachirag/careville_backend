package fitnessCenter

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity/subEntity"
	"careville_backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx = context.Background()

// @Summary Add fitnessCenter
// @Tags fitnessCenter
// @Description Add fitnessCenter
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider formData services.FitnessCenterRequestDto true "add fitnessCenter"
// @Param fitnessCenterImage formData file false "fitnessCenterImage"
// @Param certificate formData file false "certificate"
// @Param license formData file false "license"
// @Produce json
// @Success 200 {object} services.FitnessCenterResDto
// @Router /provider/services/add-fitness-center [post]
func AddFitnessCenter(c *fiber.Ctx) error {
	var (
		servicesColl  = database.GetCollection("service")
		data          services.FitnessCenterRequestDto
		fitnessCenter subEntity.UpdateServiceSubEntity
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.FitnessCenterResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Access the MultipartForm directly from the fiber.Ctx
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.FitnessCenterResDto{
			Status:  false,
			Message: "Failed to get multipart form: " + err.Error(),
		})
	}

	fitnessCenterImageFiles := form.File["fitnessCenterImage"]
	if len(fitnessCenterImageFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.FitnessCenterResDto{
			Status:  false,
			Message: "No fitnessCenterImage uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range fitnessCenterImageFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.FitnessCenterResDto{
				Status:  false,
				Message: "Failed to upload fitnessCenterImage to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("fitnessCenter/%v-image-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		fitnessCenterImage, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.FitnessCenterResDto{
				Status:  false,
				Message: "Failed to upload fitnessCenterImage to S3: " + err.Error(),
			})
		}

		if fitnessCenter.FitnessCenter != nil {
			fitnessCenter.FitnessCenter.Information.Image = fitnessCenterImage
		}

	}

	cerificateFiles := form.File["certificate"]
	licenseFiles := form.File["license"]
	if len(cerificateFiles) == 0 && len(licenseFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.FitnessCenterResDto{
			Status:  false,
			Message: "At least one document is mandatary",
		})
	}
	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range cerificateFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.FitnessCenterResDto{
				Status:  false,
				Message: "Failed to upload certificate to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("certificate/%v-doc-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		certificateURL, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.FitnessCenterResDto{
				Status:  false,
				Message: "Failed to upload certificate to S3: " + err.Error(),
			})
		}

		if fitnessCenter.FitnessCenter != nil {
			fitnessCenter.FitnessCenter.Documents.Certificate = certificateURL
		}

		// fitnessCenter.FitnessCenter.Documents.Certificate = certificateURL
	}

	// Combine the loops for license and certificate

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range licenseFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.FitnessCenterResDto{
				Status:  false,
				Message: "Failed to upload license to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("license/%v-doc-%s", id.Hex(), formFile.Filename)

		licenseURL, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.FitnessCenterResDto{
				Status:  false,
				Message: "Failed to upload license to S3: " + err.Error(),
			})
		}

		if fitnessCenter.FitnessCenter != nil {
			fitnessCenter.FitnessCenter.Documents.License = licenseURL
		}

	}

	longitude, err := strconv.ParseFloat(data.FitnessCenterReqDto.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.FitnessCenterResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.FitnessCenterReqDto.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.FitnessCenterResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	var additionalServices []subEntity.AdditionalServicesUpdateServiceSubEntity
	for _, inv := range data.FitnessCenterReqDto.AdditionalServices {
		convertedInv := subEntity.AdditionalServicesUpdateServiceSubEntity{
			Id:          primitive.NewObjectID(),
			Name:        inv.Name,
			Information: inv.Information,
		}
		additionalServices = append(additionalServices, convertedInv)
	}

	var trainers []subEntity.TrainersUpdateServiceSubEntity
	for _, inv := range data.FitnessCenterReqDto.Trainers {
		convertedInv := subEntity.TrainersUpdateServiceSubEntity{
			Id:          primitive.NewObjectID(),
			Category:    inv.Category,
			Name:        inv.Name,
			Information: inv.Information,
			Price:       inv.Price,
		}
		trainers = append(trainers, convertedInv)
	}

	var subscription []subEntity.SubscriptionUpdateServiceSubEntity
	for _, inv := range data.FitnessCenterReqDto.Subscription {
		convertedInv := subEntity.SubscriptionUpdateServiceSubEntity{
			Type:    inv.Type,
			Details: inv.Details,
			Price:   inv.Price,
		}
		subscription = append(subscription, convertedInv)
	}

	var fitnessCenterImage string
	var licenseDoc string
	var certificate string
	if fitnessCenter.FitnessCenter != nil {
		fitnessCenterImage = fitnessCenter.FitnessCenter.Information.Image
		licenseDoc = fitnessCenter.FitnessCenter.Documents.License
		certificate = fitnessCenter.FitnessCenter.Documents.Certificate
	}

	fitnessCenterData := subEntity.FitnessCenterUpdateServiceSubEntity{
		Information: subEntity.InformationUpdateServiceSubEntity{
			Name:           data.FitnessCenterReqDto.InformationName,
			AdditionalText: data.FitnessCenterReqDto.AdditionalText,
			Image:          fitnessCenterImage,
			Address: subEntity.AddressUpdateServiceSubEntity{
				Coordinates: []float64{longitude, latitude},
				Add:         data.FitnessCenterReqDto.Address,
				Type:        "Point",
			},
		},
		Documents: subEntity.DocumentsUpdateServiceSubEntity{
			Certificate: certificate,
			License:     licenseDoc,
		},
		AdditionalServices: additionalServices,
		Trainers:           trainers,
		Subscription:       subscription,
	}

	fitnessCenter = subEntity.UpdateServiceSubEntity{
		Role:                 "healthFacility",
		FacilityOrProfession: "fitnessCenter",
		ServiceStatus:        "pending",
		FitnessCenter:        &fitnessCenterData,
		UpdatedAt:            time.Now().UTC(),
	}

	fitnessCenterUpdate := bson.M{"$set": fitnessCenter}
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}
	_, err = servicesColl.UpdateOne(ctx, filter, fitnessCenterUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.FitnessCenterResDto{
			Status:  false,
			Message: "Failed to insert fitness center data into MongoDB: " + err.Error(),
		})
	}

	fitnessRes := services.FitnessCenterResDto{
		Status:  true,
		Message: "fitness center added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(fitnessRes)
}
