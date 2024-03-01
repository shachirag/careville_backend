package laboratory

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

// @Summary Add laboratory
// @Tags laboratory
// @Description Add laboratory
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider formData services.LaboratoryRequestDto true "add laboratory"
// @Param laboratoryImage formData file false "laboratoryImage"
// @Param certificate formData file false "certificate"
// @Param license formData file false "license"
// @Produce json
// @Success 200 {object} services.LaboratoryResDto
// @Router /provider/services/add-laboratory [post]
func AddLaboratory(c *fiber.Ctx) error {
	var (
		servicesColl             = database.GetCollection("service")
		data                     services.LaboratoryRequestDto
		laboratory               subEntity.UpdateServiceSubEntity
		laboratoryImageUrl       string
		laboratoryLicenceUrl     string
		laboratoryCertificateUrl string
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.LaboratoryResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Access the MultipartForm directly from the fiber.Ctx
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.LaboratoryResDto{
			Status:  false,
			Message: "Failed to get multipart form: " + err.Error(),
		})
	}

	formFiles := form.File["laboratoryImage"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.LaboratoryResDto{
			Status:  false,
			Message: "Laboratory image is required",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.LaboratoryResDto{
				Status:  false,
				Message: "Failed to upload laboratoryImage to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("laboratory/%v-image-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		laboratoryImage, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.LaboratoryResDto{
				Status:  false,
				Message: "Failed to upload laboratoryImage to S3: " + err.Error(),
			})
		}
		laboratoryImageUrl = laboratoryImage
	}

	cerificateFiles := form.File["certificate"]
	licenseFiles := form.File["license"]
	if len(cerificateFiles) == 0 && len(licenseFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.LaboratoryResDto{
			Status:  false,
			Message: "At least one document is mandatary",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range cerificateFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.LaboratoryResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.LaboratoryResDto{
				Status:  false,
				Message: "Failed to upload certificate to S3: " + err.Error(),
			})
		}
		laboratoryLicenceUrl = certificateURL
		// Append the image URL to the Images field

	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range licenseFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.LaboratoryResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.LaboratoryResDto{
				Status:  false,
				Message: "Failed to upload license to S3: " + err.Error(),
			})
		}
		laboratoryLicenceUrl = licenseURL
		// Append the image URL to the Images field

	}

	longitude, err := strconv.ParseFloat(data.LaboratoryReqDto.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.LaboratoryResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.LaboratoryReqDto.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.LaboratoryResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	var investigations []subEntity.InvestigationsUpdateServiceSubEntity
	for _, inv := range data.LaboratoryReqDto.Investigations {
		convertedInv := subEntity.InvestigationsUpdateServiceSubEntity{
			Id:          primitive.NewObjectID(),
			Type:        inv.Type,
			Name:        inv.Name,
			Information: inv.Information,
			Price:       inv.Price,
		}
		investigations = append(investigations, convertedInv)
	}

	laboratoryData := &subEntity.LaboratoryUpdateServiceSubEntity{
		Review: subEntity.Review{
			TotalReviews: 0,
			AvgRating:    0,
		},
		Information: subEntity.InformationUpdateServiceSubEntity{
			Name:           data.LaboratoryReqDto.InformationName,
			AdditionalText: data.LaboratoryReqDto.AdditionalText,
			Image:          laboratoryImageUrl,
			Address: subEntity.AddressUpdateServiceSubEntity{
				Coordinates: []float64{longitude, latitude},
				Add:         data.LaboratoryReqDto.Address,
				Type:        "Point",
			},
			IsEmergencyAvailable: false,
		},
		Documents: subEntity.DocumentsUpdateServiceSubEntity{
			Certificate: laboratoryCertificateUrl,
			License:     laboratoryLicenceUrl,
		},
		Investigations: investigations,
	}

	currentCount, err := servicesColl.CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "Failed to count documents in service collection: " + err.Error(),
		})
	}

	profileID := fmt.Sprintf("%06d", currentCount+1)

	laboratory = subEntity.UpdateServiceSubEntity{
		Role:                 "healthFacility",
		FacilityOrProfession: "laboratory",
		ServiceStatus:        "pending",
		ProfileId:            profileID,
		Laboratory:           laboratoryData,
		UpdatedAt:            time.Now().UTC(),
	}

	laboratoryUpdate := bson.M{"$set": laboratory}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}

	_, err = servicesColl.UpdateOne(ctx, filter, laboratoryUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.LaboratoryResDto{
			Status:  false,
			Message: "Failed to insert laboratoryUpdate data into MongoDB: " + err.Error(),
		})
	}

	laboratoryRes := services.LaboratoryResDto{
		Status:  true,
		Message: "Laboratory added successfully",
		Role: services.Role{
			Role:                 "healthFacility",
			FacilityOrProfession: "laboratory",
			ServiceStatus:        "pending",
			Image:                laboratoryImageUrl,
			Name:                 data.LaboratoryReqDto.InformationName,
			IsEmergencyAvailable: false,
		},
	}
	return c.Status(fiber.StatusOK).JSON(laboratoryRes)
}
