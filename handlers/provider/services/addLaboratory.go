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

// @Summary Add laboratory
// @Tags services
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
		servicesColl = database.GetCollection("service")
		data         services.LaboratoryRequestDto
		laboratory   entity.ServiceEntity
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
			Message: "No laboratoryImage uploaded",
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
		if laboratory.Laboratory != nil {
			laboratory.Laboratory.Information.Image = laboratoryImage
		}
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
		if laboratory.Laboratory != nil {
			laboratory.Laboratory.Documents.Certificate = certificateURL
		}
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
		if laboratory.Laboratory != nil {
			laboratory.Laboratory.Documents.License = licenseURL
		}
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

	var investigations []entity.Investigations
	for _, inv := range data.LaboratoryReqDto.Investigations {
		convertedInv := entity.Investigations{
			Type:        inv.Type,
			Name:        inv.Name,
			Information: inv.Information,
			Price:       inv.Price,
		}
		investigations = append(investigations, convertedInv)
	}

	var laboratoryImage string
	var licenseDoc string
	var certificate string
	if laboratory.Laboratory != nil {
		laboratoryImage = laboratory.Laboratory.Information.Image
		licenseDoc = laboratory.Laboratory.Documents.License
		certificate = laboratory.Laboratory.Documents.Certificate
	}

	laboratoryData := entity.Laboratory{
		Information: entity.Information{
			Name:           data.LaboratoryReqDto.InformationName,
			AdditionalText: data.LaboratoryReqDto.AdditionalText,
			Image:          laboratoryImage,
			Address: entity.Address{
				Coordinates: []float64{longitude, latitude},
				Add:         data.LaboratoryReqDto.Address,
				Type:        "Point",
			},
		},
		Documents: entity.Documents{
			Certificate: certificate,
			License:     licenseDoc,
		},
		Investigations: investigations,
	}

	laboratory = entity.ServiceEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthFacility",
		FacilityOrProfession: "laboratory",
		ServiceStatus:        "pending",
		Laboratory:           &laboratoryData,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}

	_, err = servicesColl.InsertOne(ctx, laboratory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.LaboratoryResDto{
			Status:  false,
			Message: "Failed to insert laboratory data into MongoDB: " + err.Error(),
		})
	}

	laboratoryRes := services.LaboratoryResDto{
		Status:  true,
		Message: "laboratory added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(laboratoryRes)
}
