package customerAuth

import (
	"careville_backend/database"
	"careville_backend/dto/customer/customerAuth"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get customer by ID
// @Tags customer authorization
// @Description Get customer by ID
// @Accept json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} customerAuth.GetCustomerResDto
// @Router /customer/get-customer-info [get]
func GetCustomer(c *fiber.Ctx) error {

	customerColl := database.GetCollection("customer")

	customerData := customerMiddleware.GetCustomerMiddlewareData(c)

	filter := bson.M{"_id": customerData.CustomerId}

	var customer entity.CustomerEntity

	projection := bson.M{
		"_id":                      1,
		"firstName":                1,
		"lastName":                 1,
		"image":                    1,
		"age":                      1,
		"sex":                      1,
		"phoneNumber.dialCode":     1,
		"phoneNumber.countryCode":  1,
		"phoneNumber.number":       1,
		"notification.deviceType":  1,
		"notification.deviceToken": 1,
		"notification.isEnabled":   1,
		"email":                    1,
		"createdAt":                1,
		"updatedAt":                1,
	}

	findOptions := options.FindOne().SetProjection(projection)
	err := customerColl.FindOne(ctx, filter, findOptions).Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(customerAuth.GetCustomerResDto{
			Status:  false,
			Message: "customer not found",
		})
	}

	customerRes := customerAuth.GetCustomerResDto{
		Status:  true,
		Message: "customer found",
		Data: customerAuth.CustomerResDto{
			Id:        customer.Id,
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Email:     customer.Email,
			Image:     customer.Image,
			PhoneNumber: customerAuth.PhoneNumber{
				DialCode:    customer.PhoneNumber.DialCode,
				CountryCode: customer.PhoneNumber.CountryCode,
				Number:      customer.PhoneNumber.Number,
			},
			Notification: customerAuth.Notification{
				DeviceToken: customer.Notification.DeviceToken,
				DeviceType:  customer.Notification.DeviceType,
				IsEnabled:   customer.Notification.IsEnabled,
			},
			Age:       customer.Age,
			Sex:       customer.Sex,
			CreatedAt: customer.CreatedAt,
			UpdatedAt: customer.UpdatedAt,
		},
	}

	return c.Status(fiber.StatusOK).JSON(customerRes)
}
