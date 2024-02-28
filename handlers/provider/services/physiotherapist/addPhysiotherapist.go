package physiotherapist

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

// @Summary Add Physiotherapist
// @Tags physiotherapist
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
// @Router /provider/services/add-physiotherapist [post]
func AddPhysiotherapist(c *fiber.Ctx) error {
	var (
		servicesColl            = database.GetCollection("service")
		data                    services.PhysiotherapistRequestDto
		physiotherapist         subEntity.UpdateServiceSubEntity
		physiotherapistImageUrl string
		nimcDoc                 string
		personalLicense         string
		professionalLicense     string
		professionalCertificate string
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
			Message: "Physiotherapist image is required",
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

		physiotherapistImageUrl = physiotherapistImage

	}

	professionalCertificateFiles := form.File["professionalCertificate"]
	professionalLicenseFormFiles := form.File["professionalLicense"]
	if len(professionalCertificateFiles) == 0 && len(professionalLicenseFormFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "At least one document is mandatary",
		})
	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range professionalCertificateFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
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
			return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
				Status:  false,
				Message: "Failed to upload professional license to S3: " + err.Error(),
			})
		}

		professionalLicense = licenseURL

	}

	personalNimcFiles := form.File["personalNimc"]
	personalLicenseFiles := form.File["personalLicense"]
	if len(personalNimcFiles) == 0 && len(personalLicenseFiles) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "At least one document is mandatary",
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

		nimcDoc = nimcURL

	}

	// Upload each image to S3 and get the S3 URLs
	for _, formFile := range personalNimcFiles {
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

		personalLicense = licenseURL

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

	var schedule []subEntity.ServiceAndScheduleUpdateServiceSubEntity
	for _, scheduleItem := range data.PhysiotherapistReqDto.Schedule {
		var slots []subEntity.SlotsUpdateServiceSubEntity
		for _, slot := range scheduleItem.Slots {
			breakingSlots := generateBreakingSlots(slot.StartTime, slot.EndTime)

			scheduleSlot := subEntity.SlotsUpdateServiceSubEntity{
				StartTime:     slot.StartTime,
				EndTime:       slot.EndTime,
				Days:          slot.Days,
				BreakingSlots: breakingSlots,
			}
			slots = append(slots, scheduleSlot)
		}
		scheduleData := subEntity.ServiceAndScheduleUpdateServiceSubEntity{
			Id:          primitive.NewObjectID(),
			Name:        scheduleItem.Name,
			ServiceFees: scheduleItem.ServiceFees,
			Slots:       slots,
		}

		schedule = append(schedule, scheduleData)
	}

	if physiotherapist.Physiotherapist != nil {
		physiotherapist.Physiotherapist.ServiceAndSchedule = schedule
	}

	physiotherapistData := subEntity.PhysiotherapistUpdateServiceSubEntity{
		Information: subEntity.InformationUpdateServiceSubEntity{
			Name:           data.PhysiotherapistReqDto.InformationName,
			AdditionalText: data.PhysiotherapistReqDto.AdditionalText,
			Image:          physiotherapistImageUrl,
			Address: subEntity.AddressUpdateServiceSubEntity{
				Coordinates: []float64{longitude, latitude},
				Add:         data.PhysiotherapistReqDto.Address,
				Type:        "Point",
			},
		},
		ProfessionalDetails: subEntity.ProfessionalDetailsUpdateServiceSubEntity{
			Qualifications: data.PhysiotherapistReqDto.Qualifications,
		},
		PersonalIdentificationDocs: subEntity.PersonalIdentificationDocsUpdateServiceSubEntity{
			Nimc:    nimcDoc,
			License: personalLicense,
		},
		ProfessionalDetailsDocs: subEntity.ProfessionalDetailsDocsUpdateServiceSubEntity{
			Certificate: professionalCertificate,
			License:     professionalLicense,
		},
		ServiceAndSchedule: schedule,
	}

	physiotherapist = subEntity.UpdateServiceSubEntity{
		Role:                 "healthProfessional",
		FacilityOrProfession: "physiotherapist",
		ServiceStatus:        "pending",
		Physiotherapist:      &physiotherapistData,
		UpdatedAt:            time.Now().UTC(),
	}

	physiotherapistUpdate := bson.M{"$set": physiotherapist}

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}

	_, err = servicesColl.UpdateOne(ctx, filter, physiotherapistUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
			Status:  false,
			Message: "Failed to insert physiotherapist data into MongoDB: " + err.Error(),
		})
	}

	fitnessRes := services.PhysiotherapistResDto{
		Status:  true,
		Message: "Physiotherapist added successfully",
		Role: services.Role{
			Role:                 "healthProfessional",
			FacilityOrProfession: "physiotherapist",
			ServiceStatus:        "pending",
			Image:                physiotherapistImageUrl,
			Name:                 data.PhysiotherapistReqDto.InformationName,
		},
	}
	return c.Status(fiber.StatusOK).JSON(fitnessRes)
}

func generateBreakingSlots(startTime, endTime string) []subEntity.BreakingSlots {
	layout := "15:04"
	start, _ := time.Parse(layout, startTime)
	end, _ := time.Parse(layout, endTime)

	if start.After(end) {
		return []subEntity.BreakingSlots{}
	}

	var breakingSlots []subEntity.BreakingSlots

	for start.Before(end) {
		next := start.Add(20 * time.Minute)
		if next.After(end) {
			break
		}
		breakingSlots = append(breakingSlots, subEntity.BreakingSlots{
			StartTime: start.Format(layout),
			EndTime:   next.Format(layout),
		})
		start = next
	}

	return breakingSlots
}
