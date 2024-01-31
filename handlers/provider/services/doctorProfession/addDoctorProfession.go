package doctorProfession

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

// @Summary Add doctorProfession
// @Tags doctorProfession
// @Description Add doctorProfession
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider formData services.DoctorProfessionRequestDto true "add doctorProfession"
// @Param doctorImage formData file false "doctorImage"
// @Param professionalCertificate formData file false "professionalCertificate"
// @Param professionalLicense formData file false "professionalLicense"
// @Param personalLicense formData file false "personalLicense"
// @Param personalNimc formData file false "personalNimc"
// @Produce json
// @Success 200 {object} services.DoctorProfessionResDto
// @Router /provider/services/add-doctor-profession [post]
func AddDoctorProfession(c *fiber.Ctx) error {
	var (
		servicesColl            = database.GetCollection("service")
		data                    services.DoctorProfessionRequestDto
		doctorProfession        subEntity.UpdateServiceSubEntity
		doctoryImageUrl         string
		nimcDoc                 string
		personalLicense         string
		professionalLicense     string
		professionalCertificate string
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Access the MultipartForm directly from the fiber.Ctx
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "Failed to get multipart form: " + err.Error(),
		})
	}

	doctorImageFiles := form.File["doctorImage"]
	if len(doctorImageFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "No doctorImage uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range doctorImageFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload doctorImage to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("doctor/%v-image-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		doctorImage, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload doctorImage to S3: " + err.Error(),
			})
		}

		doctoryImageUrl = doctorImage

	}

	professionalCertificateFiles := form.File["professionalCertificate"]
	professionalLicenseFormFiles := form.File["professionalLicense"]
	if len(professionalCertificateFiles) == 0 && len(professionalLicenseFormFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "At least one document is mandatary",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range professionalCertificateFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload professional certificate to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("certificate/%v-doc-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		certificateURL, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload professional certificate to S3: " + err.Error(),
			})
		}

		professionalCertificate = certificateURL

	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range professionalLicenseFormFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload professional license to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("license/%v-doc-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		licenseURL, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload professional license to S3: " + err.Error(),
			})
		}

		professionalLicense = licenseURL

	}

	personalNimcFiles := form.File["personalNimc"]
	personalLicenseFiles := form.File["personalLicense"]
	if len(personalNimcFiles) == 0 && len(personalLicenseFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "At least one document is mandatary",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range personalNimcFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload personalNimc to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("nimc/%v-doc-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		nimcURL, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload personalNimc to S3: " + err.Error(),
			})
		}

		nimcDoc = nimcURL

	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range personalLicenseFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload personalLicense to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("license/%v-doc-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		licenseURL, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload personalLicense to S3: " + err.Error(),
			})
		}

		personalLicense = licenseURL

	}

	longitude, err := strconv.ParseFloat(data.DoctorProfessionReqDto.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.DoctorProfessionReqDto.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	var slots []subEntity.DoctorSlotsUpdateServiceSubEntity
	for _, slot := range data.DoctorProfessionReqDto.Slots {
		scheduleSlot := subEntity.DoctorSlotsUpdateServiceSubEntity{
			Id:        primitive.NewObjectID(),
			StartTime: slot.StartTime,
			EndTime:   slot.EndTime,
			Days:      slot.Days,
		}
		slots = append(slots, scheduleSlot)
	}

	if doctorProfession.Doctor != nil && len(doctorProfession.Doctor.Schedule.Slots) > 0 {
		doctorProfession.Doctor.Schedule.Slots = slots
	}

	doctorData := subEntity.DoctorProfessionUpdateServiceSubEntity{
		Information: subEntity.InformationUpdateServiceSubEntity{
			Name:           data.DoctorProfessionReqDto.InformationName,
			AdditionalText: data.DoctorProfessionReqDto.AdditionalText,
			Image:          doctoryImageUrl,
			Address: subEntity.AddressUpdateServiceSubEntity{
				Coordinates: []float64{longitude, latitude},
				Add:         data.DoctorProfessionReqDto.Address,
				Type:        "Point",
			},
			IsEmergencyAvailable: false,
		},
		AdditionalServices: subEntity.AdditionalServiceUpdateServiceSubEntity{
			Qualifications: data.DoctorProfessionReqDto.Qualifications,
			Speciality:     data.DoctorProfessionReqDto.Speciality,
		},
		PersonalIdentificationDocs: subEntity.PersonalIdentificationDocsUpdateServiceSubEntity{
			Nimc:    nimcDoc,
			License: personalLicense,
		},
		ProfessionalDetailsDocs: subEntity.ProfessionalDetailsDocsUpdateServiceSubEntity{
			Certificate: professionalCertificate,
			License:     professionalLicense,
		},
		Schedule: subEntity.DoctorScheduleUpdateServiceSubEntity{
			ConsultationFees: data.DoctorProfessionReqDto.ConsultationFees,
			Slots:            slots,
		},
	}

	doctorProfession = subEntity.UpdateServiceSubEntity{
		Role:                 "healthProfessional",
		FacilityOrProfession: "doctor",
		ServiceStatus:        "pending",
		Doctor:               &doctorData,
		UpdatedAt:            time.Now().UTC(),
	}

	doctorProfessionUpdate := bson.M{"$set": doctorProfession}
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}
	_, err = servicesColl.UpdateOne(ctx, filter, doctorProfessionUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "Failed to insert doctor profession data into MongoDB: " + err.Error(),
		})
	}

	fitnessRes := services.DoctorProfessionResDto{
		Status:  true,
		Message: "Doctor Profession added successfully",
		Role: services.Role{
			Role:                 "healthProfessional",
			FacilityOrProfession: "doctor",
			ServiceStatus:        "pending",
			Image:                doctoryImageUrl,
			Name:                 data.DoctorProfessionReqDto.InformationName,
			IsEmergencyAvailable: false,
		},
	}
	return c.Status(fiber.StatusOK).JSON(fitnessRes)
}
