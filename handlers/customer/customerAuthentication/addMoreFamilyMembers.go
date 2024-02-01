package customerAuth

import (
	"careville_backend/database"
	"careville_backend/dto/customer/customerAuth"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Add more members
// @Tags customer authorization
// @Description Add more members
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param  customer body customerAuth.AddMemberReqDto true "add members"
// @Produce json
// @Success 200 {object} customerAuth.AddMemberResDto
// @Router /customer/profile/add-more-family-member [post]
func AddMoreMembers(c *fiber.Ctx) error {
	var (
		customerColl = database.GetCollection("customer")
		data         customerAuth.AddMemberReqDto
		customer     entity.CustomerEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(customerAuth.AddMemberResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	customerData := customerMiddleware.GetCustomerMiddlewareData(c)

	filter := bson.M{
		"_id": customerData.CustomerId,
	}

	err = customerColl.FindOne(ctx, filter).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorImageResDto{
				Status:  false,
				Message: "customer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "Failed to fetch customer from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{
		"$addToSet": bson.M{
			"familyMembers": bson.M{
				"$each": []entity.FamilyMembers{
					{
						Id:           primitive.NewObjectID(),
						Name:         data.Name,
						Age:          data.Age,
						RelationShip: data.RelationShip,
						Sex:          data.Sex,
					},
				},
			},
		},
	}

	updateRes, err := customerColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.AddMemberResDto{
			Status:  false,
			Message: "Failed to update customer data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(customerAuth.AddMemberResDto{
			Status:  false,
			Message: "customer not found",
		})
	}

	hospClinicRes := customerAuth.AddMemberResDto{
		Status:  true,
		Message: "Member added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
