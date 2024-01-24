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
		servicesColl     = database.GetCollection("service")
		data             services.DoctorProfessionRequestDto
		doctorProfession subEntity.UpdateServiceSubEntity
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

		if doctorProfession.Doctor != nil {
			doctorProfession.Doctor.Information.Image = doctorImage
		}

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

		if doctorProfession.Doctor != nil {
			doctorProfession.Doctor.ProfessionalDetailsDocs.Certificate = certificateURL
		}

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

		if doctorProfession.Doctor != nil {
			doctorProfession.Doctor.ProfessionalDetailsDocs.License = licenseURL
		}

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

		if doctorProfession.Doctor != nil {
			doctorProfession.Doctor.PersonalIdentificationDocs.Nimc = nimcURL
		}

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

		if doctorProfession.Doctor != nil {
			doctorProfession.Doctor.PersonalIdentificationDocs.License = licenseURL
		}

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

	// Parse and add DoctorSchedule data
	// Parse and add DoctorSchedule data
	var schedule []subEntity.DoctorScheduleUpdateServiceSubEntity
	for _, scheduleItem := range data.DoctorProfessionReqDto.Schedule {
		var slots []subEntity.SlotsUpdateServiceSubEntity
		for _, slot := range scheduleItem.Slots {
			scheduleSlot := subEntity.SlotsUpdateServiceSubEntity{
				StartTime: slot.StartTime,
				EndTime:   slot.EndTime,
				Days:      slot.Days,
			}
			slots = append(slots, scheduleSlot)
		}
		scheduleData := subEntity.DoctorScheduleUpdateServiceSubEntity{
			ConsultationFees: scheduleItem.ConsultationFees,
			Slots:            slots,
		}
		schedule = append(schedule, scheduleData)
	}

	if doctorProfession.Doctor != nil {
		doctorProfession.Doctor.Schedule = schedule
	}

	var doctoryImage string
	var nimcDoc string
	var personalLicense string
	var professionalLicense string
	var professionalCertificate string
	if doctorProfession.Doctor != nil {
		doctoryImage = doctorProfession.Doctor.Information.Image
		nimcDoc = doctorProfession.Doctor.PersonalIdentificationDocs.Nimc
		personalLicense = doctorProfession.Doctor.PersonalIdentificationDocs.License
		professionalLicense = doctorProfession.Doctor.ProfessionalDetailsDocs.License
		professionalCertificate = doctorProfession.Doctor.ProfessionalDetailsDocs.Certificate
	}

	doctorData := subEntity.DoctorProfessionUpdateServiceSubEntity{
		Information: subEntity.InformationUpdateServiceSubEntity{
			Name:           data.DoctorProfessionReqDto.InformationName,
			AdditionalText: data.DoctorProfessionReqDto.AdditionalText,
			Image:          doctoryImage,
			Address: subEntity.AddressUpdateServiceSubEntity{
				Coordinates: []float64{longitude, latitude},
				Add:         data.DoctorProfessionReqDto.Address,
				Type:        "Point",
			},
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
		Schedule: schedule,
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
		Message: "doctor profession added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(fitnessRes)
}
