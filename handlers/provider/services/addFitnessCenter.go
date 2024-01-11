package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"careville_backend/database"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"careville_backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Add fitnessCenter
// @Tags services
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
		fitnessCenter entity.ServiceEntity
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

	formFiles := form.File["fitnessCenterImage"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.FitnessCenterResDto{
			Status:  false,
			Message: "No fitnessCenterImage uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
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

	formFiles = form.File["certificate"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.FitnessCenterResDto{
			Status:  false,
			Message: "No certificate uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
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

	}

	formFiles = form.File["license"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.FitnessCenterResDto{
			Status:  false,
			Message: "No license uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
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

		// Upload the image to S3 and get the S3 URL
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

	var additionalServices []entity.AdditionalServices
	for _, inv := range data.FitnessCenterReqDto.AdditionalServices {
		convertedInv := entity.AdditionalServices{
			Name:        inv.Name,
			Information: inv.Information,
		}
		additionalServices = append(additionalServices, convertedInv)
	}

	var trainers []entity.Trainers
	for _, inv := range data.FitnessCenterReqDto.Trainers {
		convertedInv := entity.Trainers{
			Category:    inv.Category,
			Name:        inv.Name,
			Information: inv.Information,
			Price:       inv.Price,
		}
		trainers = append(trainers, convertedInv)
	}

	var subscription []entity.Subscription
	for _, inv := range data.FitnessCenterReqDto.Subscription {
		convertedInv := entity.Subscription{
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

	fitnessCenterData := entity.FitnessCenter{
		Information: entity.Information{
			Name:           data.FitnessCenterReqDto.InformationName,
			AdditionalText: data.FitnessCenterReqDto.AdditionalText,
			Image:          fitnessCenterImage,
			Address: entity.Address{
				Coordinates: []float64{longitude, latitude},
				Add:         data.FitnessCenterReqDto.Address,
				Type:        "Point",
			},
		},
		Documents: entity.Documents{
			Certificate: certificate,
			License:     licenseDoc,
		},
		AdditionalServices: additionalServices,
		Trainers:           trainers,
		Subscription:       subscription,
	}

	fitnessCenter = entity.ServiceEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthFacility",
		FacilityOrProfession: "fitnessCenter",
		Status:               "pending",
		FitnessCenter:        &fitnessCenterData,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}

	_, err = servicesColl.InsertOne(ctx, fitnessCenter)
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
