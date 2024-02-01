package customerAuth

import (
	"careville_backend/database"
	customerAuth "careville_backend/dto/customer/customerAuth"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Update customer
// @Description Update customer
// @Tags customer authorization
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param customer formData customerAuth.UpdateCustomerReqDto true "Update data of customer"
// @Produce json
// @Success 200 {object} customerAuth.UpdateCustomerResDto
// @Router /customer/profile/update-customer-info [put]
func UpdateCustomer(c *fiber.Ctx) error {

	var (
		customerColl = database.GetCollection("customer")
		data         customerAuth.UpdateCustomerReqDto
		provider     entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(customerAuth.UpdateCustomerResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Get provider data from middleware
	customerData := customerMiddleware.GetCustomerMiddlewareData(c)

	filter := bson.M{"_id": customerData.CustomerId}
	err = customerColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(providerAuth.GetProviderResDto{
				Status:  false,
				Message: "customer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.GetProviderResDto{
			Status:  false,
			Message: "Failed to fetch customer from MongoDB: " + err.Error(),
		})
	}

	longitude, err := strconv.ParseFloat(data.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(customerAuth.UpdateCustomerResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(customerAuth.UpdateCustomerResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	update := bson.M{"$set": bson.M{
		"firstName": data.FirstName,
		"lastName":  data.LastName,
		"phoneNumber": bson.M{
			"dialCode":    data.DialCode,
			"number":      data.PhoneNumber,
			"countryCode": data.CountryCode,
		},
		"address": customerAuth.Address{
			Coordinates: []float64{longitude, latitude},
			Type:        "Point",
			Add:         data.Address,
		},
		"age":       data.Age,
		"sex":       data.Sex,
		"updatedAt": time.Now().UTC(),
	}}

	opts := options.Update().SetUpsert(true)
	// Execute the update operation
	updateRes, err := customerColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.UpdateCustomerResDto{
			Status:  false,
			Message: "Failed to update customer data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(customerAuth.UpdateCustomerResDto{
			Status:  false,
			Message: "customer not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(customerAuth.UpdateCustomerResDto{
		Status:  true,
		Message: "Customer data updated successfully",
	})
}
