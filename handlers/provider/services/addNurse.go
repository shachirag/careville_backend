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

// @Summary Add nurse
// @Tags services
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
		servicesColl = database.GetCollection("service")
		data         services.NurseRequestDto
		nurse        entity.ServiceEntity
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

		if nurse.Nurse != nil {
			nurse.Nurse.Information.Image = nurseImage
		}

	}

	formFiles = form.File["professionalCertificate"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
			Status:  false,
			Message: "No professional certificate uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
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

		if nurse.Nurse != nil {
			nurse.Nurse.ProfessionalDetailsDocs.Certificate = certificateURL
		}

	}

	formFiles = form.File["professionalLicense"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
			Status:  false,
			Message: "No professional license uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
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

		if nurse.Nurse != nil {
			nurse.Nurse.ProfessionalDetailsDocs.License = licenseURL
		}

	}

	formFiles = form.File["personalNimc"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
			Status:  false,
			Message: "No personalNimc uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
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

		if nurse.Nurse != nil {
			nurse.Nurse.PersonalIdentificationDocs.Nimc = nimcURL
		}

	}

	formFiles = form.File["personalLicense"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.NurseResDto{
			Status:  false,
			Message: "No personalLicense uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
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

		if nurse.Nurse != nil {
			nurse.Nurse.PersonalIdentificationDocs.License = licenseURL
		}

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
	var schedule []entity.ServiceAndSchedule
	for _, scheduleItem := range data.NurseReqDto.Schedule {
		scheduleData := entity.ServiceAndSchedule{
			Name:        scheduleItem.Name,
			ServiceFees: scheduleItem.ServiceFees,
			Slots: entity.Slots{
				StartTime: scheduleItem.Slots.StartTime,
				EndTime:   scheduleItem.Slots.EndTime,
				Days:      scheduleItem.Slots.Days,
			},
		}

		schedule = append(schedule, scheduleData)
	}

	if nurse.Nurse != nil {
		nurse.Nurse.Schedule = schedule
	}

	var nurseImage string
	var nimcDoc string
	var personalLicense string
	var professionalLicense string
	var professionalCertificate string
	if nurse.Nurse != nil {
		nurseImage = nurse.Nurse.Information.Image
		nimcDoc = nurse.Nurse.PersonalIdentificationDocs.Nimc
		personalLicense = nurse.Nurse.PersonalIdentificationDocs.License
		professionalLicense = nurse.Nurse.ProfessionalDetailsDocs.License
		professionalCertificate = nurse.Nurse.ProfessionalDetailsDocs.Certificate
	}

	NurseData := entity.Nurse{
		Information: entity.Information{
			FirstName:      data.NurseReqDto.InformationFirstName,
			LastName:       data.NurseReqDto.InformationLastName,
			AdditionalText: data.NurseReqDto.AdditionalText,
			Image:          nurseImage,
			Address: entity.Address{
				Coordinates: []float64{longitude, latitude},
				Add:         data.NurseReqDto.Address,
				Type:        "Point",
			},
		},
		ProfessionalDetails: entity.ProfessionalDetails{
			Qualifications: data.NurseReqDto.Qualifications,
		},
		PersonalIdentificationDocs: entity.PersonalIdentificationDocs{
			Nimc:    nimcDoc,
			License: personalLicense,
		},
		ProfessionalDetailsDocs: entity.ProfessionalDetailsDocs{
			Certificate: professionalCertificate,
			License:     professionalLicense,
		},
		Schedule: schedule,
	}

	nurse = entity.ServiceEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthProfessional",
		FacilityOrProfession: "nurse",
		Status:               "pending",
		Nurse:                &NurseData,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}

	_, err = servicesColl.InsertOne(ctx, nurse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NurseResDto{
			Status:  false,
			Message: "Failed to insert doctor profession data into MongoDB: " + err.Error(),
		})
	}

	fitnessRes := services.NurseResDto{
		Status:  true,
		Message: "doctor profession added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(fitnessRes)
}
