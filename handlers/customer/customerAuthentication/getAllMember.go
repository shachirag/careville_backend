package customerAuth

import (
	"careville_backend/database"
	"careville_backend/dto/customer/customerAuth"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all investigations
// @Description Get all investigations
// @Tags customer authorization
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} customerAuth.MembeResDto
// @Router /customer/profile/get-members [get]
func GetMembers(c *fiber.Ctx) error {

	var customer entity.CustomerEntity

	customerData := customerMiddleware.GetCustomerMiddlewareData(c)

	customerColl := database.GetCollection("customer")

	filter := bson.M{
		"_id": customerData.CustomerId,
	}

	projection := bson.M{
		"familyMembers.id":           1,
		"familyMembers.name":         1,
		"familyMembers.age":          1,
		"familyMembers.sex":          1,
		"familyMembers.relationShip": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := customerColl.FindOne(ctx, filter, findOptions).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(customerAuth.MembeResDto{
				Status:  false,
				Message: "member not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.MembeResDto{
			Status:  false,
			Message: "Failed to fetch member from MongoDB: " + err.Error(),
		})
	}

	membersData := make([]customerAuth.MembeRes, 0)
	if customer.FamilyMembers != nil && len(customer.FamilyMembers) > 0 {
		for _, member := range customer.FamilyMembers {
			membersData = append(membersData, customerAuth.MembeRes{
				Id:           member.Id,
				Sex:          member.Sex,
				Name:         member.Name,
				RelationShip: member.RelationShip,
				Age:          member.Age,
			})
		}
	}

	if len(membersData) == 0 {
		return c.Status(fiber.StatusOK).JSON(customerAuth.MembeResDto{
			Status:  false,
			Message: "No member data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(customerAuth.MembeResDto{
		Status:  true,
		Message: "Members retrieved successfully",
		Data:    membersData,
	})
}
