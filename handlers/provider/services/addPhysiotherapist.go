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

// @Summary Add Physiotherapist
// @Tags services
// @Description Add Physiotherapist
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider formData services.PhysiotherapistRequestDto true "add Physiotherapist"
// @Param physiotherapistImage formData file false "physiotherapistImage"
// @Param professionalCertificate formData file false "professionalCertificate"
// @Param professionalLicense formData file false "professionalLicense"
// @Param personalLicense formData file false "personalLicense"
// @Param personalNimc formData file false "personalNimc"
// @Produce json
// @Success 200 {object} services.PhysiotherapistResDto
// @Router /provider/add-physiotherapist [post]
func AddPhysiotherapist(c *fiber.Ctx) error {
	var (
		servicesColl    = database.GetCollection("service")
		data            services.PhysiotherapistRequestDto
		physiotherapist entity.ServiceEntity
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Access the MultipartForm directly from the fiber.Ctx
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "Failed to get multipart form: " + err.Error(),
		})
	}

	formFiles := form.File["physiotherapistImage"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "No physiotherapistImage uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
				Status:  false,
				Message: "Failed to upload doctorImage to S3: " + err.Error(),
			})
		}

		// Append the image URL to the Images field
		physiotherapist.Physiotherapist.Information.Image = physiotherapistImage
	}

	formFiles = form.File["professionalCertificate"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "No certificate uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
				Status:  false,
				Message: "Failed to upload certificate to S3: " + err.Error(),
			})
		}

		// Append the image URL to the Images field
		physiotherapist.Physiotherapist.ProfessionalDetailsDocs.Certificate = certificateURL
	}

	formFiles = form.File["professionalLicense"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "No license uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
				Status:  false,
				Message: "Failed to upload license to S3: " + err.Error(),
			})
		}

		// Append the image URL to the Images field
		physiotherapist.Physiotherapist.ProfessionalDetailsDocs.License = licenseURL
	}

	formFiles = form.File["personalNimc"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "No personalNimc uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
				Status:  false,
				Message: "Failed to upload personalNimc to S3: " + err.Error(),
			})
		}

		// Append the image URL to the Images field
		physiotherapist.Physiotherapist.PersonalIdentificationDocs.Nimc = nimcURL
	}

	formFiles = form.File["personalLicense"]
	if len(formFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "No personalLicense uploaded",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
				Status:  false,
				Message: "Failed to upload personalLicense to S3: " + err.Error(),
			})
		}

		// Append the image URL to the Images field
		physiotherapist.Physiotherapist.PersonalIdentificationDocs.License = licenseURL
	}

	longitude, err := strconv.ParseFloat(data.PhysiotherapistReqDto.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.PhysiotherapistReqDto.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	// Parse and add DoctorSchedule data
	// Parse and add NurseSchedule data
	var schedule []entity.ServiceAndSchedule
	for _, scheduleItem := range data.PhysiotherapistReqDto.Schedule {
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

	physiotherapist.Physiotherapist.ServiceAndSchedule = schedule

	// Assign schedule to doctorProfession
	// physiotherapist.Physiotherapist.ServiceAndSchedule = schedule

	physiotherapist = entity.ServiceEntity{
		Id:                   primitive.NewObjectID(),
		ProviderId:           data.PhysiotherapistReqDto.ProviderId,
		Role:                 data.PhysiotherapistReqDto.Role,
		FacilityOrProfession: data.PhysiotherapistReqDto.FacilityOrProfession,
		IsApproved:           false,
		Physiotherapist: entity.Physiotherapist{
			Information: entity.Information{
				Name:           data.PhysiotherapistReqDto.InformationName,
				AdditionalText: data.PhysiotherapistReqDto.AdditionalText,
				Image:          physiotherapist.Physiotherapist.Information.Image,
				Address: entity.Address{
					Coordinates: []float64{longitude, latitude},
					Add:         data.PhysiotherapistReqDto.Address,
					Type:        "Point",
				},
			},
			ProfessionalDetails: entity.ProfessionalDetails{
				Qualifications: data.PhysiotherapistReqDto.Qualifications,
			},
			PersonalIdentificationDocs: entity.PersonalIdentificationDocs{
				Nimc:    physiotherapist.Physiotherapist.PersonalIdentificationDocs.Nimc,
				License: physiotherapist.Physiotherapist.PersonalIdentificationDocs.License,
			},
			ProfessionalDetailsDocs: entity.ProfessionalDetailsDocs{
				Certificate: physiotherapist.Physiotherapist.ProfessionalDetailsDocs.Certificate,
				License:     physiotherapist.Physiotherapist.ProfessionalDetailsDocs.License,
			},
			ServiceAndSchedule: schedule,
		},

		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = servicesColl.InsertOne(ctx, physiotherapist)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "Failed to insert Physiotherapist data into MongoDB: " + err.Error(),
		})
	}

	fitnessRes := services.NurseResDto{
		Status:  true,
		Message: "Physiotherapist added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(fitnessRes)
}
