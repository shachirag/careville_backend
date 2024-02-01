package customerAuth

import (
	"careville_backend/database"
	"careville_backend/dto/customer/customerAuth"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// ChangeProviderPassword is the handler for changing provider passwords
// @Summary Change customer Password
// @Description Change customer Password
// @Tags customer authorization
// @Accept application/json
// @Param Authorization header string true "Authentication header"
// @Param customer body customerAuth.CustomerChangePasswordReqDto true "Change password of customer"
// @Produce json
// @Success 200 {object} customerAuth.CustomerChangePasswordResDto
// @Router /customer/change-password [put]
func ChangeCustomerPassword(c *fiber.Ctx) error {
	var (
		customerColl = database.GetCollection("customer")
		data         customerAuth.CustomerChangePasswordReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.CustomerChangePasswordResDto{
			Status:  false,
			Message: "Failed to parse request body: " + err.Error(),
		})
	}

	customerData := customerMiddleware.GetCustomerMiddlewareData(c)

	filter := bson.M{"_id": customerData.CustomerId}

	result := customerColl.FindOne(c.Context(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(customerAuth.CustomerChangePasswordResDto{
				Status:  false,
				Message: "Customer not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.CustomerChangePasswordResDto{
			Status:  false,
			Message: "Error finding customer: " + err.Error(),
		})
	}

	var customer entity.CustomerEntity
	err = result.Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.CustomerChangePasswordResDto{
			Status:  false,
			Message: "Failed to decode customer data: " + err.Error(),
		})
	}

	if data.CurrentPassword != "" {
		err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(data.CurrentPassword))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(customerAuth.CustomerChangePasswordResDto{
				Status:  false,
				Message: "Current password is incorrect",
			})
		}
	}

	// Hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 6)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.CustomerChangePasswordResDto{
			Status:  false,
			Message: "Failed to hash the new password",
		})
	}

	update := bson.M{
		"$set": bson.M{
			"password": string(hashedNewPassword),
		},
	}

	// Execute the update operation
	updateRes, err := customerColl.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.CustomerChangePasswordResDto{
			Status:  false,
			Message: "Failed to update customer password in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(customerAuth.CustomerChangePasswordResDto{
			Status:  false,
			Message: "Customer not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(customerAuth.CustomerChangePasswordResDto{
		Status:  true,
		Message: "Customer password updated successfully",
	})
}
