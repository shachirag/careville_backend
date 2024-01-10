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

// @Summary Add pharmacy
// @Tags services
// @Description Add pharmacy
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider formData services.PharmacyRequestDto true "add pharmacy"
// @Param pharmacyImage formData file false "pharmacyImage"
// @Param certificate formData file false "certificate"
// @Param license formData file false "license"
// @Produce json
// @Success 200 {object} services.PharmacyResDto
// @Router /provider/add-pharmacy [post]
func AddPharmacy(c *fiber.Ctx) error {
	var (
		servicesColl = database.GetCollection("service")
		data         services.PharmacyRequestDto
		pharmacy     entity.ServiceEntity
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.PharmacyResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Access the MultipartForm directly from the fiber.Ctx
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.PharmacyResDto{
			Status:  false,
			Message: "Failed to get multipart form: " + err.Error(),
		})
	}

	formFiles := form.File["pharmacyImage"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PharmacyResDto{
			Status:  false,
			Message: "No pharmacyImage uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.PharmacyResDto{
				Status:  false,
				Message: "Failed to upload pharmacyImage to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("pharmacy/%v-image-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		pharmacyImage, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.PharmacyResDto{
				Status:  false,
				Message: "Failed to upload pharmacyImage to S3: " + err.Error(),
			})
		}

		// Append the image URL to the Images field
		pharmacy.Pharmacy.Information.Image = pharmacyImage
	}

	formFiles = form.File["certificate"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PharmacyResDto{
			Status:  false,
			Message: "No certificate uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.PharmacyResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.PharmacyResDto{
				Status:  false,
				Message: "Failed to upload certificate to S3: " + err.Error(),
			})
		}

		// Append the image URL to the Images field
		pharmacy.Pharmacy.Documents.Certificate = certificateURL
	}

	formFiles = form.File["license"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PharmacyResDto{
			Status:  false,
			Message: "No license uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.PharmacyResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.PharmacyResDto{
				Status:  false,
				Message: "Failed to upload license to S3: " + err.Error(),
			})
		}

		// Append the image URL to the Images field
		pharmacy.Pharmacy.Documents.License = licenseURL
	}

	longitude, err := strconv.ParseFloat(data.PharmacyReqDto.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.PharmacyResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.PharmacyReqDto.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.PharmacyResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	var additionalServices []entity.AdditionalServices
	for _, inv := range data.PharmacyReqDto.AdditionalServices {
		convertedInv := entity.AdditionalServices{
			Name:        inv.Name,
			Information: inv.Information,
		}
		additionalServices = append(additionalServices, convertedInv)
	}

	pharmacyData := entity.Pharmacy{
		Information: entity.Information{
			Name:           data.PharmacyReqDto.InformationName,
			AdditionalText: data.PharmacyReqDto.AdditionalText,
			Image:          pharmacy.Pharmacy.Information.Image,
			Address: entity.Address{
				Coordinates: []float64{longitude, latitude},
				Add:         data.PharmacyReqDto.Address,
				Type:        "Point",
			},
		},
		Documents: entity.Documents{
			Certificate: pharmacy.Pharmacy.Documents.Certificate,
			License:     pharmacy.Pharmacy.Documents.License,
		},
		AdditionalServices: additionalServices,
	}

	pharmacy = entity.ServiceEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthFacility",
		FacilityOrProfession: "pharmacy",
		Status:               "pending",
		Pharmacy:             &pharmacyData,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}

	_, err = servicesColl.InsertOne(ctx, pharmacy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.PharmacyResDto{
			Status:  false,
			Message: "Failed to insert pharmacy data into MongoDB: " + err.Error(),
		})
	}

	pharmacyRes := services.PharmacyResDto{
		Status:  true,
		Message: "pharmacy added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(pharmacyRes)
}
