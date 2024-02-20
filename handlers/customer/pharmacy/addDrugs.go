package pharmacy

import (
	"fmt"
	"time"

	"careville_backend/database"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/dto/customer/pharmacy"
	"careville_backend/entity"
	"careville_backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Add drugs
// @Tags customer pharmacy
// @Description Add drugs
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param serviceId query string true "service ID"
// @Param  customer formData pharmacy.PharmacyDrugsReqDto true "add pharmacy drugs"
// @Param prescriptionImages formData file true "prescription images"
// @Produce json
// @Success 200 {object} pharmacy.PharmacyDrugsResDto
// @Router /customer/healthFacility/add-pharmacy-drug [post]
func AddPharmacyDrugs(c *fiber.Ctx) error {

	var (
		appointmentColl = database.GetCollection("appointment")
		serviceColl     = database.GetCollection("service")
		customerColl    = database.GetCollection("customer")
		data            pharmacy.PharmacyDrugsReqDto
		appointment     entity.AppointmentEntity
		service         entity.ServiceEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(pharmacy.PharmacyDrugsResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	serviceId := c.Query("serviceId")

	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(pharmacy.PharmacyDrugsResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pharmacy.PharmacyDrugsResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	informationFilter := bson.M{
		"_id": serviceObjectID,
	}

	informationProjection := bson.M{
		"pharmacy.information.name":  1,
		"pharmacy.information.image": 1,
		"pharmacy.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
	}

	informationOpts := options.FindOne().SetProjection(informationProjection)

	err = serviceColl.FindOne(ctx, informationFilter, informationOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.PharmacyDrugsResDto{
			Status:  false,
			Message: "Failed to fetch pharmacy information data: " + err.Error(),
		})
	}

	if service.Pharmacy == nil {
		return c.Status(fiber.StatusNotFound).JSON(pharmacy.PharmacyDrugsResDto{
			Status:  false,
			Message: "Pharmacy data not found",
		})
	}

	customerMiddlewareData := customerMiddleware.GetCustomerMiddlewareData(c)

	var customer entity.CustomerEntity

	customerProjection := bson.M{
		"_id":       1,
		"firstName": 1,
		"lastName":  1,
		"image":     1,
		"email":     1,
		"phoneNumber": bson.M{
			"dialCode":    1,
			"number":      1,
			"countryCode": 1,
		},
	}

	customerOpts := options.FindOne().SetProjection(customerProjection)

	customerFilter := bson.M{
		"_id": customerMiddlewareData.CustomerId,
	}

	err = customerColl.FindOne(ctx, customerFilter, customerOpts).Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.PharmacyDrugsResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.PharmacyDrugsResDto{
			Status:  false,
			Message: "Failed to get multipart form: " + err.Error(),
		})
	}

	formFiles := form.File["prescriptionImages"]
	// if len(formFiles) == 0 {
	// 	return c.Status(fiber.StatusBadRequest).JSON(pharmacy.PharmacyDrugsResDto{
	// 		Status:  false,
	// 		Message: "No prescription uploaded",
	// 	})
	// }

	var informationName string
	var informationImage string
	var informationAddress entity.Address
	if service.Pharmacy != nil {
		informationName = service.Pharmacy.Information.Name
		informationImage = service.Pharmacy.Information.Image
		informationAddress = service.Pharmacy.Information.Address
	}

	// longitude, err := strconv.ParseFloat(data.Longitude, 64)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(pharmacy.PharmacyDrugsResDto{
	// 		Status:  false,
	// 		Message: "Invalid longitude format",
	// 	})
	// }

	// latitude, err := strconv.ParseFloat(data.Latitude, 64)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(pharmacy.PharmacyDrugsResDto{
	// 		Status:  false,
	// 		Message: "Invalid latitude format",
	// 	})
	// }

	drugData := entity.PharmacyAppointmentEntity{
		RequestedDrugs: entity.RequestedDrugsAppointmentEntity{
			ModeOfDelivery:  data.ModeOfDelivery,
			NameAndQuantity: data.NameAndQuantity,
			Address: entity.Address{
				Coordinates: customer.Address.Coordinates,
				Type:        customer.Address.Type,
				Add:         customer.Address.Add,
			},
			Prescription: make([]string, 0),
		},
		Information: entity.PharmacyInformationAppointmentEntity{
			Name:    informationName,
			Image:   informationImage,
			Address: informationAddress,
		},
	}

	appointment = entity.AppointmentEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthFacility",
		FacilityOrProfession: "pharmacy",
		ServiceID:            serviceObjectID,
		Customer: entity.CustomerAppointmentEntity{
			ID:          customerMiddlewareData.CustomerId,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Image:       customer.Image,
			Email:       customer.Email,
			PhoneNumber: customer.PhoneNumber,
		},
		Pharmacy:          &drugData,
		PaymentStatus:     "initiated",
		AppointmentStatus: "pending",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	for _, formFile := range formFiles {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(pharmacy.PharmacyDrugsResDto{
				Status:  false,
				Message: "Failed to upload prescription image to S3: " + err.Error(),
			})
		}

		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("prescription/%v-documents-%s", id.Hex(), formFile.Filename)

		imageURL, err := utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.PharmacyDrugsResDto{
				Status:  false,
				Message: "Failed to upload prescription to S3: " + err.Error(),
			})
		}

		if appointment.Pharmacy != nil {
			appointment.Pharmacy.RequestedDrugs.Prescription = append(appointment.Pharmacy.RequestedDrugs.Prescription, imageURL)
		}

	}

	_, err = appointmentColl.InsertOne(ctx, appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.PharmacyDrugsResDto{
			Status:  false,
			Message: "Failed to insert pharmacy appointment data into MongoDB: " + err.Error(),
		})
	}

	pharmacyRes := pharmacy.PharmacyDrugsResDto{
		Status:  true,
		Message: "Appointment added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(pharmacyRes)
}
