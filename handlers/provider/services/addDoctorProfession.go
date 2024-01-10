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

// @Summary Add doctorProfession
// @Tags services
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
// @Router /provider/add-doctor-profession [post]
func AddDoctorProfession(c *fiber.Ctx) error {
	var (
		servicesColl     = database.GetCollection("service")
		data             services.DoctorProfessionRequestDto
		doctorProfession entity.ServiceEntity
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

	formFiles := form.File["doctorImage"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "No doctorImage uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
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

	formFiles = form.File["professionalCertificate"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "No certificate uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload certificate to S3: " + err.Error(),
			})
		}
		if doctorProfession.Doctor != nil {
			doctorProfession.Doctor.ProfessionalDetailsDocs.Certificate = certificateURL
		}
		// Append the image URL to the Images field

	}

	formFiles = form.File["professionalLicense"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "No license uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionResDto{
				Status:  false,
				Message: "Failed to upload license to S3: " + err.Error(),
			})
		}
		if doctorProfession.Doctor != nil {
			doctorProfession.Doctor.ProfessionalDetailsDocs.License = licenseURL
		}
		// Append the image URL to the Images field

	}

	formFiles = form.File["personalNimc"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "No personalNimc uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
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
		// Append the image URL to the Images field

	}

	formFiles = form.File["personalLicense"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.DoctorProfessionResDto{
			Status:  false,
			Message: "No personalLicense uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
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
		// Append the image URL to the Images field

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
	var schedule []entity.DoctorSchedule
	for _, scheduleItem := range data.DoctorProfessionReqDto.Schedule {
		scheduleData := entity.DoctorSchedule{
			ConsultationFees: scheduleItem.ConsultationFees,
			Slots: entity.Slots{
				StartTime: scheduleItem.Slots.StartTime,
				EndTime:   scheduleItem.Slots.EndTime,
				Days:      scheduleItem.Slots.Days,
			},
		}

		schedule = append(schedule, scheduleData)
	}

	if doctorProfession.Doctor != nil {
		doctorProfession.Doctor.Schedule = schedule
	}

	var doctoryImage string
	if doctorProfession.Doctor != nil {
		doctoryImage = doctorProfession.Doctor.Information.Image
	}

	var nimcDoc string
	if doctorProfession.Doctor != nil {
		nimcDoc = doctorProfession.Doctor.PersonalIdentificationDocs.Nimc
	}

	var personalLicense string
	if doctorProfession.Doctor != nil {
		personalLicense = doctorProfession.Doctor.PersonalIdentificationDocs.License
	}

	var professionalLicense string
	if doctorProfession.Doctor != nil {
		professionalLicense = doctorProfession.Doctor.ProfessionalDetailsDocs.License
	}

	var professionalCertificate string
	if doctorProfession.Doctor != nil {
		professionalCertificate = doctorProfession.Doctor.ProfessionalDetailsDocs.Certificate
	}

	doctorData := entity.DoctorEntityDto{
		Information: entity.Information{
			Name:           data.DoctorProfessionReqDto.InformationName,
			AdditionalText: data.DoctorProfessionReqDto.AdditionalText,
			Image:          doctoryImage,
			Address: entity.Address{
				Coordinates: []float64{longitude, latitude},
				Add:         data.DoctorProfessionReqDto.Address,
				Type:        "Point",
			},
		},
		AdditionalServices: entity.AdditionalService{
			Qualifications: data.DoctorProfessionReqDto.Qualifications,
			Speciality:     data.DoctorProfessionReqDto.Speciality,
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

	doctorProfession = entity.ServiceEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthProfessional",
		FacilityOrProfession: "doctor",
		Status:               "pending",
		Doctor:               &doctorData,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}

	_, err = servicesColl.InsertOne(ctx, doctorProfession)
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
