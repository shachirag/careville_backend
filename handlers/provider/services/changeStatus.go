package services

// import (
// 	"careville_backend/database"
// 	"careville_backend/dto/provider/services"
// 	"github.com/gofiber/fiber/v2"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // @Summary change status
// // @Description change status
// // @Tags services
// // @Accept multipart/form-data
// //
// //	@Param Authorization header	string true	"Authentication header"
// //
// // @Param status query string false "change status approved or rejected"
// // @Produce json
// // @Success 200 {object} services.StatusResDto
// // @Router /provider/change-status/{id} [put]
// func ChangeStatus(c *fiber.Ctx) error {

// 	var (
// 		serviceColl = database.GetCollection("service")
// 	)

// 	status := c.Query("status", "")

// 	serviceID := c.Params("id")
// 	if serviceID == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(services.StatusResDto{
// 			Status:  false,
// 			Message: "service ID is missing in the request",
// 		})
// 	}

// 	objID, err := primitive.ObjectIDFromHex(serviceID)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(services.StatusResDto{
// 			Status:  false,
// 			Message: "Invalid service ID",
// 		})
// 	}

// 	filter := bson.M{
// 		"_id": objID,
// 	}

// 	result := serviceColl.FindOne(ctx, filter)
// 	if result.Err() != nil {
// 		if result.Err() == mongo.ErrNoDocuments {
// 			return c.Status(fiber.StatusNotFound).JSON(services.StatusResDto{
// 				Status:  false,
// 				Message: "service not found",
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(services.StatusResDto{
// 			Status:  false,
// 			Message: "internal server error " + err.Error(),
// 		})
// 	}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"status": status,
// 		},
// 	}

// 	statusRes, err := serviceColl.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(services.StatusResDto{
// 			Status:  false,
// 			Message: "Failed to change status in MongoDB: " + err.Error(),
// 		})
// 	}

// 	if statusRes.MatchedCount == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(services.StatusResDto{
// 			Status:  false,
// 			Message: "service not found",
// 		})
// 	}

// 	var actionMessage string
// 	if status == "approved" {
// 		actionMessage = "approved"
// 	} else if status == "rejected" {
// 		actionMessage = "rejected"
// 	}

// 	return c.Status(fiber.StatusOK).JSON(services.StatusResDto{
// 		Status:  true,
// 		Message: actionMessage,
// 	})
// }
