package nurse

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

// @Summary Add nurse
// @Tags nurse
// @Description Add nurse
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider formData services.NurseRequestDto true "add nurse"
// @Param nurseImage formData file false "nurseImage"
// @Param professionalCertificate formData file false "professionalCertificate"
// @Param professionalLicense formData file false "professionalLicense"
// @Param personalLicense formData file false "personalLicense"
// @Param personalNimc formData file false "personalNimc"
// @Produce json
// @Success 200 {object} services.NurseResDto
// @Router /provider/services/add-nurse [post]
func AddNurse(c *fiber.Ctx) error {
	var (
		servicesColl            = database.GetCollection("service")
		data                    services.NurseRequestDto
		nurse                   subEntity.UpdateServiceSubEntity
		nurseImage              string
		nimcDoc                 string
		personalLicense         string
		professionalLicense     string
		professionalCertificate string
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NurseResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Access the MultipartForm directly from the fiber.Ctx
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NurseResDto{
			Status:  false,
			Message: "Failed to get multipart form: " + err.Error(),
		})
	}

	formFiles := form.File["nurseImage"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
			Status:  false,
			Message: "No nurseImage uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
				Status:  false,
				Message: "Failed to upload nurseImage to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("nurse/%v-image-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		nurseImage, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.NurseResDto{
				Status:  false,
				Message: "Failed to upload doctorImage to S3: " + err.Error(),
			})
		}

		nurseImage = nurseImage

	}

	professionalCertificateFiles := form.File["professionalCertificate"]
	professionalLicenseFormFiles := form.File["professionalLicense"]
	if len(professionalCertificateFiles) == 0 && len(professionalLicenseFormFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
			Status:  false,
			Message: "At least one document is mandatary",
		})
	}
	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range professionalCertificateFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.NurseResDto{
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
			return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.NurseResDto{
				Status:  false,
				Message: "Failed to upload professional license to S3: " + err.Error(),
			})
		}

		professionalLicense = licenseURL

	}

	personalNimcFiles := form.File["personalNimc"]
	personalLicenseFiles := form.File["personalLicense"]
	if len(personalNimcFiles) == 0 && len(personalLicenseFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
			Status:  false,
			Message: "At least one document is mandatary",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range personalNimcFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.NurseResDto{
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
			return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.NurseResDto{
				Status:  false,
				Message: "Failed to upload personalLicense to S3: " + err.Error(),
			})
		}

		personalLicense = licenseURL
	}

	longitude, err := strconv.ParseFloat(data.NurseReqDto.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.NurseReqDto.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	// Parse and add DoctorSchedule data
	// Parse and add NurseSchedule data
	var schedule []subEntity.ServiceAndScheduleUpdateServiceSubEntity
	for _, scheduleItem := range data.NurseReqDto.Schedule {
		var slots []subEntity.SlotsUpdateServiceSubEntity
		for _, slot := range scheduleItem.Slots {
			scheduleSlot := subEntity.SlotsUpdateServiceSubEntity{
				StartTime: slot.StartTime,
				EndTime:   slot.EndTime,
				Days:      slot.Days,
			}
			slots = append(slots, scheduleSlot)
		}
		scheduleData := subEntity.ServiceAndScheduleUpdateServiceSubEntity{
			Name:        scheduleItem.Name,
			ServiceFees: scheduleItem.ServiceFees,
			Slots:       slots,
		}

		schedule = append(schedule, scheduleData)
	}

	if nurse.Nurse != nil {
		nurse.Nurse.Schedule = schedule
	}

	NurseData := subEntity.NurseUpdateServiceSubEntity{
		Information: subEntity.InformationUpdateServiceSubEntity{
			Name:           data.NurseReqDto.InformationName,
			AdditionalText: data.NurseReqDto.AdditionalText,
			Image:          nurseImage,
			Address: subEntity.AddressUpdateServiceSubEntity{
				Coordinates: []float64{longitude, latitude},
				Add:         data.NurseReqDto.Address,
				Type:        "Point",
			},
			IsEmergencyAvailable: false,
		},
		ProfessionalDetails: subEntity.ProfessionalDetailsUpdateServiceSubEntity{
			Qualifications: data.NurseReqDto.Qualifications,
		},
		PersonalIdentificationDocs: subEntity.PersonalIdentificationDocsUpdateServiceSubEntity{
			Nimc:    nimcDoc,
			License: personalLicense,
		},
		ProfessionalDetailsDocs: subEntity.ProfessionalDetailsDocsUpdateServiceSubEntity{
			Certificate: professionalCertificate,
			License:     professionalLicense,
		},
		Schedule: schedule,
	}

	nurse = subEntity.UpdateServiceSubEntity{
		Role:                 "healthProfessional",
		FacilityOrProfession: "nurse",
		ServiceStatus:        "pending",
		Nurse:                &NurseData,
		UpdatedAt:            time.Now().UTC(),
	}

	nurseUpdate := bson.M{"$set": nurse}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}

	_, err = servicesColl.UpdateOne(ctx, filter, nurseUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NurseResDto{
			Status:  false,
			Message: "Failed to insert nurse data into MongoDB: " + err.Error(),
		})
	}

	fitnessRes := services.NurseResDto{
		Status:  true,
		Message: "doctor profession added successfully",
		Role: services.Role{
			Role:                 "healthProfessional",
			FacilityOrProfession: "nurse",
			ServiceStatus:        "pending",
			Image:                nurseImage,
			Name:                 data.NurseReqDto.InformationName,
			IsEmergencyAvailable: false,
		},
	}
	return c.Status(fiber.StatusOK).JSON(fitnessRes)
}
