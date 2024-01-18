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

// @Summary Add MedicalLabScientist
// @Tags services
// @Description Add MedicalLabScientist
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider formData services.MedicalLabScientistRequestDto true "add MedicalLabScientist"
// @Param medicalLabScientistImage formData file false "physiotherapistImage"
// @Param professionalCertificate formData file false "professionalCertificate"
// @Param professionalLicense formData file false "professionalLicense"
// @Param personalLicense formData file false "personalLicense"
// @Param personalNimc formData file false "personalNimc"
// @Produce json
// @Success 200 {object} services.MedicalLabScientistResDto
// @Router /provider/services/add-medicalLab-scientist [post]
func AddMedicalLabScientist(c *fiber.Ctx) error {
	var (
		servicesColl        = database.GetCollection("service")
		data                services.MedicalLabScientistRequestDto
		medicalLabScientist entity.ServiceEntity
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.MedicalLabScientistResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Access the MultipartForm directly from the fiber.Ctx
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.MedicalLabScientistResDto{
			Status:  false,
			Message: "Failed to get multipart form: " + err.Error(),
		})
	}

	formFiles := form.File["medicalLabScientistImage"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
			Status:  false,
			Message: "No medicalLabScientistImage uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
				Status:  false,
				Message: "Failed to upload nurseImage to S3: " + err.Error(),
			})
		}

		// Generate a unique filename for each image
		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("physiotherapist/%v-image-%s", id.Hex(), formFile.Filename)

		// Upload the image to S3 and get the S3 URL
		physiotherapistImage, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.MedicalLabScientistResDto{
				Status:  false,
				Message: "Failed to upload doctorImage to S3: " + err.Error(),
			})
		}

		if medicalLabScientist.MedicalLabScientist != nil {
			medicalLabScientist.MedicalLabScientist.Information.Image = physiotherapistImage
		}

	}

	formFiles = form.File["professionalCertificate"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
			Status:  false,
			Message: "No professional certificate uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.MedicalLabScientistResDto{
				Status:  false,
				Message: "Failed to upload professional certificate to S3: " + err.Error(),
			})
		}

		if medicalLabScientist.MedicalLabScientist != nil {
			medicalLabScientist.MedicalLabScientist.ProfessionalDetailsDocs.Certificate = certificateURL
		}

	}

	formFiles = form.File["professionalLicense"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
			Status:  false,
			Message: "No professional license uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.MedicalLabScientistResDto{
				Status:  false,
				Message: "Failed to upload professional license to S3: " + err.Error(),
			})
		}

		if medicalLabScientist.MedicalLabScientist != nil {
			medicalLabScientist.MedicalLabScientist.ProfessionalDetailsDocs.License = licenseURL
		}

	}

	formFiles = form.File["personalNimc"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
			Status:  false,
			Message: "No personalNimc uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.MedicalLabScientistResDto{
				Status:  false,
				Message: "Failed to upload personalNimc to S3: " + err.Error(),
			})
		}

		if medicalLabScientist.MedicalLabScientist != nil {
			medicalLabScientist.MedicalLabScientist.PersonalIdentificationDocs.Nimc = nimcURL
		}

	}

	formFiles = form.File["personalLicense"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
			Status:  false,
			Message: "No personalLicense uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.MedicalLabScientistResDto{
				Status:  false,
				Message: "Failed to upload personalLicense to S3: " + err.Error(),
			})
		}

		if medicalLabScientist.MedicalLabScientist != nil {
			medicalLabScientist.MedicalLabScientist.PersonalIdentificationDocs.License = licenseURL
		}

	}

	longitude, err := strconv.ParseFloat(data.MedicalLabScientistReqDto.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.MedicalLabScientistReqDto.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.MedicalLabScientistResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	var schedule []entity.ServiceAndSchedule
	for _, scheduleItem := range data.MedicalLabScientistReqDto.Schedule {
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

	if medicalLabScientist.MedicalLabScientist != nil {
		medicalLabScientist.MedicalLabScientist.ServiceAndSchedule = schedule
	}

	var medicalLabScientistImage string
	var nimcDoc string
	var personalLicense string
	var professionalLicense string
	var professionalCertificate string
	if medicalLabScientist.MedicalLabScientist != nil {
		medicalLabScientistImage = medicalLabScientist.MedicalLabScientist.Information.Image
		nimcDoc = medicalLabScientist.MedicalLabScientist.PersonalIdentificationDocs.Nimc
		personalLicense = medicalLabScientist.MedicalLabScientist.PersonalIdentificationDocs.License
		professionalLicense = medicalLabScientist.MedicalLabScientist.ProfessionalDetailsDocs.License
		professionalCertificate = medicalLabScientist.MedicalLabScientist.ProfessionalDetailsDocs.Certificate
	}

	MedicalLabScientistData := entity.MedicalLabScientist{
		Information: entity.Information{
			FirstName:      data.MedicalLabScientistReqDto.InformationFirstName,
			LastName:       data.MedicalLabScientistReqDto.InformationLastName,
			AdditionalText: data.MedicalLabScientistReqDto.AdditionalText,
			Image:          medicalLabScientistImage,
			Address: entity.Address{
				Coordinates: []float64{longitude, latitude},
				Add:         data.MedicalLabScientistReqDto.Address,
				Type:        "Point",
			},
		},
		PersonalDetails: entity.PersonalDetails{
			Department: data.MedicalLabScientistReqDto.Department,
			Document:   data.MedicalLabScientistReqDto.Document,
		},
		PersonalIdentificationDocs: entity.PersonalIdentificationDocs{
			Nimc:    nimcDoc,
			License: personalLicense,
		},
		ProfessionalDetailsDocs: entity.ProfessionalDetailsDocs{
			Certificate: professionalCertificate,
			License:     professionalLicense,
		},
		ServiceAndSchedule: schedule,
	}

	medicalLabScientist = entity.ServiceEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthProfessional",
		FacilityOrProfession: "medicalLabScientist",
		Status:               "pending",
		MedicalLabScientist:  &MedicalLabScientistData,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}

	_, err = servicesColl.InsertOne(ctx, medicalLabScientist)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.MedicalLabScientistResDto{
			Status:  false,
			Message: "Failed to insert   medicalLabScientist data into MongoDB: " + err.Error(),
		})
	}

	fitnessRes := services.NurseResDto{
		Status:  true,
		Message: "medicalLabScientist added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(fitnessRes)
}
