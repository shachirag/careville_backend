package services

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strconv"
// 	"sync"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"go.mongodb.org/mongo-driver/bson/primitive"

// 	"careville_backend/database"
// 	"careville_backend/dto/provider/services"
// 	"careville_backend/entity"
// 	"careville_backend/utils"
// )

// // @Summary Add Physiotherapist
// // @Tags services
// // @Description Add Physiotherapist
// // @Accept multipart/form-data
// //
// //	@Param Authorization header	string true	"Authentication header"
// //
// // @Param provider formData services.PhysiotherapistRequestDto true "add Physiotherapist"
// // @Param physiotherapistImage formData file false "physiotherapistImage"
// // @Param professionalCertificate formData file false "professionalCertificate"
// // @Param professionalLicense formData file false "professionalLicense"
// // @Param personalLicense formData file false "personalLicense"
// // @Param personalNimc formData file false "personalNimc"
// // @Produce json
// // @Success 200 {object} services.PhysiotherapistResDto
// // @Router /provider/services/add-physiotherapist [post]
// func AddPhysiotherapists(c *fiber.Ctx) error {
// 	var (
// 		servicesColl    = database.GetCollection("service")
// 		data            services.PhysiotherapistRequestDto
// 		physiotherapist entity.ServiceEntity
// 	)

// 	dataStr := c.FormValue("data")
// 	dataBytes := []byte(dataStr)

// 	err := json.Unmarshal(dataBytes, &data)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
// 			Status:  false,
// 			Message: err.Error(),
// 		})
// 	}

// 	// Access the MultipartForm directly from the fiber.Ctx
// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
// 			Status:  false,
// 			Message: "Failed to get multipart form: " + err.Error(),
// 		})
// 	}

// 	formFiles := form.File["physiotherapistImage"]
// 	professionalCertificateFiles := form.File["professionalCertificate"]
// 	professionalLicenseFormFiles := form.File["professionalLicense"]
// 	personalNimcFiles := form.File["personalNimc"]
// 	personalLicenseFiles := form.File["personalLicense"]
// 	// if len(formFiles) == 0 {
// 	// 	return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
// 	// 		Status:  false,
// 	// 		Message: "No physiotherapistImage uploaded",
// 	// 	})
// 	// }

// 	// Upload images concurrently
// 	var physiotherapistImage1, certificateURL, licenseURL, nimcURL, personalLicenseUrl string
// 	var uploadErrors []error

// 	var wg sync.WaitGroup
// 	var mu sync.Mutex

// 	uploadImage := func(formFile *fiber.File, prefix string, url *string) {
// 		defer wg.Done()

// 		file, err := formFile.Open()
// 		if err != nil {
// 			mu.Lock()
// 			uploadErrors = append(uploadErrors, fmt.Errorf("Failed to upload %s to S3: %s", prefix, err.Error()))
// 			mu.Unlock()
// 			return
// 		}

// 		id := primitive.NewObjectID()
// 		fileName := fmt.Sprintf("%s/%v-doc-%s", prefix, id.Hex(), formFile.Filename)

// 		resultURL, err := utils.UploadToS3(fileName, file)
// 		if err != nil {
// 			mu.Lock()
// 			uploadErrors = append(uploadErrors, fmt.Errorf("Failed to upload %s to S3: %s", prefix, err.Error()))
// 			mu.Unlock()
// 			return
// 		}

// 		mu.Lock()
// 		*url = resultURL
// 		mu.Unlock()
// 	}

// 	wg.Add(5)
// 	go uploadImage(formFiles[0], "physiotherapistImage", &physiotherapistImage1)
// 	go uploadImage(professionalCertificateFiles[0], "professionalCertificate", &certificateURL)
// 	go uploadImage(professionalLicenseFormFiles[0], "professionalLicense", &licenseURL)
// 	go uploadImage(personalNimcFiles[0], "personalNimc", &nimcURL)
// 	go uploadImage(personalLicenseFiles[0], "personalLicense", &personalLicenseUrl)

// 	wg.Wait()

// 	if len(uploadErrors) > 0 {
// 		return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
// 			Status:  false,
// 			Message: "Failed to upload one or more images",
// 		})
// 	}

// 	// Update the physiotherapist entity with the uploaded image URLs
// 	if physiotherapist.Physiotherapist != nil {
// 		physiotherapist.Physiotherapist.Information.Image = physiotherapistImage1
// 		physiotherapist.Physiotherapist.PersonalIdentificationDocs.Nimc = nimcURL
// 		physiotherapist.Physiotherapist.ProfessionalDetailsDocs.Certificate = certificateURL
// 		physiotherapist.Physiotherapist.ProfessionalDetailsDocs.License = licenseURL
// 	}

// 	longitude, err := strconv.ParseFloat(data.PhysiotherapistReqDto.Longitude, 64)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
// 			Status:  false,
// 			Message: "Invalid longitude format",
// 		})
// 	}

// 	latitude, err := strconv.ParseFloat(data.PhysiotherapistReqDto.Latitude, 64)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(services.PhysiotherapistResDto{
// 			Status:  false,
// 			Message: "Invalid latitude format",
// 		})
// 	}

// 	// Parse and add NurseSchedule data
// 	var schedule []entity.ServiceAndSchedule
// 	for _, scheduleItem := range data.PhysiotherapistReqDto.Schedule {
// 		var slots []entity.Slots
// 		for _, slot := range scheduleItem.Slots {
// 			scheduleSlot := entity.Slots{
// 				StartTime: slot.StartTime,
// 				EndTime:   slot.EndTime,
// 				Days:      slot.Days,
// 			}
// 			slots = append(slots, scheduleSlot)
// 		}
// 		scheduleData := entity.ServiceAndSchedule{
// 			Name:        scheduleItem.Name,
// 			ServiceFees: scheduleItem.ServiceFees,
// 			Slots:       slots,
// 		}

// 		schedule = append(schedule, scheduleData)
// 	}

// 	if physiotherapist.Physiotherapist != nil {
// 		physiotherapist.Physiotherapist.ServiceAndSchedule = schedule
// 	}

// 	var physiotherapistImage string
// 	var nimcDoc string
// 	var personalLicense string
// 	var professionalLicense string
// 	var professionalCertificate string
// 	if physiotherapist.Physiotherapist != nil {
// 		physiotherapistImage = physiotherapist.Physiotherapist.Information.Image
// 		nimcDoc = physiotherapist.Physiotherapist.PersonalIdentificationDocs.Nimc
// 		personalLicense = physiotherapist.Physiotherapist.PersonalIdentificationDocs.License
// 		professionalLicense = physiotherapist.Physiotherapist.ProfessionalDetailsDocs.License
// 		professionalCertificate = physiotherapist.Physiotherapist.ProfessionalDetailsDocs.Certificate
// 	}

// 	physiotherapistData := entity.Physiotherapist{
// 		Information: entity.Information{
// 			Name:           data.PhysiotherapistReqDto.InformationName,
// 			AdditionalText: data.PhysiotherapistReqDto.AdditionalText,
// 			Image:          physiotherapistImage,
// 			Address: entity.Address{
// 				Coordinates: []float64{longitude, latitude},
// 				Add:         data.PhysiotherapistReqDto.Address,
// 				Type:        "Point",
// 			},
// 		},
// 		ProfessionalDetails: entity.ProfessionalDetails{
// 			Qualifications: data.PhysiotherapistReqDto.Qualifications,
// 		},
// 		PersonalIdentificationDocs: entity.PersonalIdentificationDocs{
// 			Nimc:    nimcDoc,
// 			License: personalLicense,
// 		},
// 		ProfessionalDetailsDocs: entity.ProfessionalDetailsDocs{
// 			Certificate: professionalCertificate,
// 			License:     professionalLicense,
// 		},
// 		ServiceAndSchedule: schedule,
// 	}

// 	physiotherapist = entity.ServiceEntity{
// 		Id:                   primitive.NewObjectID(),
// 		Role:                 "healthProfessional",
// 		FacilityOrProfession: "physiotherpist",
// 		ServiceStatus:        "pending",
// 		Physiotherapist:      &physiotherapistData,
// 		CreatedAt:            time.Now().UTC(),
// 		UpdatedAt:            time.Now().UTC(),
// 	}

// 	_, err = servicesColl.InsertOne(ctx, physiotherapist)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistResDto{
// 			Status:  false,
// 			Message: "Failed to insert Physiotherapist data into MongoDB: " + err.Error(),
// 		})
// 	}

// 	fitnessRes := services.NurseResDto{
// 		Status:  true,
// 		Message: "Physiotherapist added successfully",
// 	}
// 	return c.Status(fiber.StatusOK).JSON(fitnessRes)
// }
