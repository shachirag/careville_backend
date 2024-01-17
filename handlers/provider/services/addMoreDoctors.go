package services

// import (
// 	"careville_backend/database"
// 	"careville_backend/dto/provider/services"
// 	"careville_backend/entity"

// 	"github.com/gofiber/fiber/v2"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// // @Summary Add HospitalClinic
// // @Tags services
// // @Description Add HospitalClinic
// // @Accept multipart/form-data
// //
// //	@Param Authorization header	string true	"Authentication header"
// //
// // @Param  provider formData services.HospitalClinicRequestDto true "add HospitalClinic"
// // @Param hospitalImage formData file false "hospitalImage"
// // @Param certificate formData file false "certificate"
// // @Param license formData file false "license"
// // @Produce json
// // @Success 200 {object} services.HospitalClinicResDto
// // @Router /provider/services/add-hospitalClinic [post]
// func AddMoreDoctors(c *fiber.Ctx) error {
// 	var (
// 		servicesColl = database.GetCollection("service")
// 		data         services.MoreDoctorReqDto
// 		hospClinic   entity.ServiceEntity
// 	)

// 	// Parsing the request body
// 	err := c.BodyParser(&data)
// 	if err != nil {
// 		return c.Status(500).JSON(services.UpdateDoctorResDto{
// 			Status:  false,
// 			Message: err.Error(),
// 		})
// 	}

// 	hospClinicData := entity.HospClinic{
// 		Doctor: doctors,
// 	}

// 	hospClinic = entity.ServiceEntity{
// 		HospClinic: &hospClinicData,
// 	}

// 	_, err = servicesColl.InsertOne(ctx, hospClinic)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
// 			Status:  false,
// 			Message: "Failed to insert doctor data into MongoDB: " + err.Error(),
// 		})
// 	}

// 	hospClinicRes := services.HospitalClinicResDto{
// 		Status:  true,
// 		Message: "doctor added successfully",
// 	}
// 	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
// }
