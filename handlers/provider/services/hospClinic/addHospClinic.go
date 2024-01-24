package hospClinic

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

// @Summary Add HospitalClinic
// @Tags hospClinic
// @Description Add HospitalClinic
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param  provider formData services.HospitalClinicRequestDto true "add HospitalClinic"
// @Param hospitalImage formData file false "hospitalImage"
// @Param certificate formData file false "certificate"
// @Param license formData file false "license"
// @Produce json
// @Success 200 {object} services.HospitalClinicResDto
// @Router /provider/services/add-hospitalClinic [post]
func AddHospClinic(c *fiber.Ctx) error {
	var (
		servicesColl = database.GetCollection("service")
		data         services.HospitalClinicRequestDto
		hospClinic   subEntity.UpdateServiceSubEntity
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Access the MultipartForm directly from the fiber.Ctx
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "Failed to get multipart form: " + err.Error(),
		})
	}

	formFiles := form.File["hospitalImage"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "No hospitalImage uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.HospitalClinicResDto{
				Status:  false,
				Message: "Failed to upload hospitalImage to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("hospital/%v-image-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		hospitalImage, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
				Status:  false,
				Message: "Failed to upload hospitalImage to S3: " + err.Error(),
			})
		}

		if hospClinic.HospClinic != nil {
			hospClinic.HospClinic.Information.Image = hospitalImage
		}

	}

	cerificateFiles := form.File["certificate"]
	licenseFiles := form.File["license"]
	if len(cerificateFiles) == 0 && len(licenseFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "At least one document is mandatary",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range cerificateFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.HospitalClinicResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
				Status:  false,
				Message: "Failed to upload certificate to S3: " + err.Error(),
			})
		}

		if hospClinic.HospClinic != nil {
			hospClinic.HospClinic.Documents.Certificate = certificateURL
		}

	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range licenseFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.HospitalClinicResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
				Status:  false,
				Message: "Failed to upload license to S3: " + err.Error(),
			})
		}

		if hospClinic.HospClinic != nil {
			hospClinic.HospClinic.Documents.License = licenseURL
		}

	}

	longitude, err := strconv.ParseFloat(data.HospitalClinicReqDto.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.HospitalClinicReqDto.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	// Convert Doctor data from request into the required structure
	doctors := make([]subEntity.DoctorUpdateServiceSubEntity, len(data.HospitalClinicReqDto.Doctor))
	for i, doc := range data.HospitalClinicReqDto.Doctor {
		schedule := make([]subEntity.ScheduleUpdateServiceSubEntity, len(doc.Schedule))
		for j, sch := range doc.Schedule {
			schedule[j] = subEntity.ScheduleUpdateServiceSubEntity{
				StartTime: sch.StartTime,
				EndTime:   sch.EndTime,
				Days:      sch.Days,
			}
		}

		doctors[i] = subEntity.DoctorUpdateServiceSubEntity{
			Id:         primitive.NewObjectID(),
			Name:       doc.Name,
			Speciality: doc.Speciality,
			Schedule:   schedule,
		}
	}

	var hospitalImage string
	var licenseDoc string
	var certificate string
	if hospClinic.HospClinic != nil {
		hospitalImage = hospClinic.HospClinic.Information.Image
		licenseDoc = hospClinic.HospClinic.Documents.License
		certificate = hospClinic.HospClinic.Documents.Certificate
	}

	hospClinicData := subEntity.HospClinicUpdateServiceSubEntity{
		Information: subEntity.InformationUpdateServiceSubEntity{
			Name:           data.HospitalClinicReqDto.InformationName,
			AdditionalText: data.HospitalClinicReqDto.AdditionalText,
			Image:          hospitalImage,
			Address: subEntity.AddressUpdateServiceSubEntity{
				Coordinates: []float64{longitude, latitude},
				Add:         data.HospitalClinicReqDto.Address,
				Type:        "Point",
			},
		},
		Documents: subEntity.DocumentsUpdateServiceSubEntity{
			Certificate: certificate,
			License:     licenseDoc,
		},
		OtherServices: data.HospitalClinicReqDto.OtherServices,
		Insurances:    data.HospitalClinicReqDto.Insurances,
		Doctor:        doctors,
	}

	hospClinic = subEntity.UpdateServiceSubEntity{
		Role:                 "healthFacility",
		FacilityOrProfession: "hospClinic",
		ServiceStatus:        "pending",
		HospClinic:           &hospClinicData,
		UpdatedAt:            time.Now().UTC(),
	}
	healthUpdate := bson.M{"$set": hospClinic}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}

	_, err = servicesColl.UpdateOne(ctx, filter, healthUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "Failed to insert hospital/clinic data into MongoDB: " + err.Error(),
		})
	}

	hospClinicRes := services.HospitalClinicResDto{
		Status:  true,
		Message: "hospital/clinic added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
