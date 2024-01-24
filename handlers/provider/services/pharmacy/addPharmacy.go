package pharmacy

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

// @Summary Add pharmacy
// @Tags pharmacy
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
// @Router /provider/services/add-pharmacy [post]
func AddPharmacy(c *fiber.Ctx) error {
	var (
		servicesColl = database.GetCollection("service")
		data         services.PharmacyRequestDto
		pharmacy     subEntity.UpdateServiceSubEntity
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

		if pharmacy.Pharmacy != nil {
			pharmacy.Pharmacy.Information.Image = pharmacyImage
		}

	}

	cerificateFiles := form.File["certificate"]
	licenseFiles := form.File["license"]
	if len(cerificateFiles) == 0 && len(licenseFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PharmacyResDto{
			Status:  false,
			Message: "At least one document is mandatary",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range cerificateFiles {
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

		if pharmacy.Pharmacy != nil {
			pharmacy.Pharmacy.Documents.Certificate = certificateURL
		}

	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range licenseFiles {
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

		if pharmacy.Pharmacy != nil {
			pharmacy.Pharmacy.Documents.License = licenseURL
		}

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

	var additionalServices []subEntity.AdditionalServicesUpdateServiceSubEntity
	for _, inv := range data.PharmacyReqDto.AdditionalServices {
		convertedInv := subEntity.AdditionalServicesUpdateServiceSubEntity{
			Id:          primitive.NewObjectID(),
			Name:        inv.Name,
			Information: inv.Information,
		}
		additionalServices = append(additionalServices, convertedInv)
	}

	var pharmacyImage string
	var licenseDoc string
	var certificate string
	if pharmacy.Pharmacy != nil {
		pharmacyImage = pharmacy.Pharmacy.Information.Image
		licenseDoc = pharmacy.Pharmacy.Documents.License
		certificate = pharmacy.Pharmacy.Documents.Certificate
	}

	pharmacyData := subEntity.PharmacyUpdateServiceSubEntity{
		Information: subEntity.InformationUpdateServiceSubEntity{
			Name:           data.PharmacyReqDto.InformationName,
			AdditionalText: data.PharmacyReqDto.AdditionalText,
			Image:          pharmacyImage,
			Address: subEntity.AddressUpdateServiceSubEntity{
				Coordinates: []float64{longitude, latitude},
				Add:         data.PharmacyReqDto.Address,
				Type:        "Point",
			},
		},
		Documents: subEntity.DocumentsUpdateServiceSubEntity{
			Certificate: certificate,
			License:     licenseDoc,
		},
		AdditionalServices: additionalServices,
	}

	pharmacy = subEntity.UpdateServiceSubEntity{
		Role:                 "healthFacility",
		FacilityOrProfession: "pharmacy",
		ServiceStatus:        "pending",
		Pharmacy:             &pharmacyData,
		UpdatedAt:            time.Now().UTC(),
	}

	pharmacyUpdate := bson.M{"$set": pharmacy}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}

	_, err = servicesColl.UpdateOne(ctx, filter, pharmacyUpdate)
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
