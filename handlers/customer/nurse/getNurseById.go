package nurse

import (
	"careville_backend/database"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/dto/customer/nurse"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get nurse by ID
// @Tags customer nurse
// @Description Get nurse by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "nurse ID"
// @Produce json
// @Success 200 {object} nurse.GetNurseResDto
// @Router /customer/healthProfessional/get-nurse/{id} [get]
func GetNurseByID(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
	)

	idParam := c.Params("id")
	nurseID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.GetNurseResDto{
			Status:  false,
			Message: "Invalid nurse ID",
		})
	}

	customerData := customerMiddleware.GetCustomerMiddlewareData(c)
	customerFilter := bson.M{
		"_id": customerData.CustomerId,
	}

	var customer entity.CustomerEntity
	err = database.GetCollection("customer").FindOne(ctx, customerFilter).Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNurseResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	isCustomerFamilyMember := false
	if len(customer.FamilyMembers) > 0 {
		isCustomerFamilyMember = true
	}

	filter := bson.M{"_id": nurseID}

	projection := bson.M{
		"nurse.information.name":           1,
		"nurse.information.image":          1,
		"_id":                              1,
		"nurse.information.additionalText": 1,
		"nurse.schedule.id":                1,
		"nurse.review.totalReviews":        1,
		"nurse.review.avgRating":           1,
		"nurse.schedule.serviceFees":       1,
		"nurse.schedule.slots": bson.M{
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	var nurseData entity.ServiceEntity
	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&nurseData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNurseResDto{
			Status:  false,
			Message: "Failed to fetch nurse data: " + err.Error(),
		})
	}

	if nurseData.Nurse == nil {
		return c.Status(fiber.StatusNotFound).JSON(nurse.GetNurseResDto{
			Status:  false,
			Message: "nurse data not found",
		})
	}

	scheduleData := make([]nurse.ServiceAndSchedule, 0)
	if nurseData.Nurse != nil && len(nurseData.Nurse.Schedule) > 0 {
		for _, service := range nurseData.Nurse.Schedule {
			var slotData []nurse.Slots
			for _, slots := range service.Slots {
				slotData = append(slotData, nurse.Slots{
					StartTime: slots.StartTime,
					EndTime:   slots.EndTime,
					Days:      slots.Days,
				})
			}
			scheduleData = append(scheduleData, nurse.ServiceAndSchedule{
				Id:          service.Id,
				ServiceFees: service.ServiceFees,
				Name:        service.Name,
				Slots:       slotData,
			})
		}

	}

	nurseRes := nurse.GetNurseResDto{
		Status:  true,
		Message: "Nurse data fetched successfully",
		Data: nurse.NurseResponse{
			Id:                     nurseData.Id,
			Image:                  nurseData.Nurse.Information.Image,
			Name:                   nurseData.Nurse.Information.Name,
			AboutMe:                nurseData.Nurse.Information.AdditionalText,
			ServiceAndSchedule:     scheduleData,
			TotalReviews:           nurseData.Nurse.Review.TotalReviews,
			AvgRating:              nurseData.Nurse.Review.AvgRating,
			IsCustomerFamilyMember: isCustomerFamilyMember,
		},
	}

	return c.Status(fiber.StatusOK).JSON(nurseRes)
}
