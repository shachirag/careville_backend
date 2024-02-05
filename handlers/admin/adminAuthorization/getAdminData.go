package adminAuth

import (
	"careville_backend/database"
	adminAuth "careville_backend/dto/admin/adminAuth"
	adminMiddleware "careville_backend/dto/admin/middleware"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Fetch Admin By ID
// @Description Fetch Admin By ID
// @Tags admin authorization
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} adminAuth.GetAdminResDto
// @Router /admin/profile/get-admin-info [get]
func FetchAdminById(c *fiber.Ctx) error {

	var admin entity.AdminEntity

	adminData := adminMiddleware.GetAdminMiddlewareData(c)
	adminColl := database.GetCollection("admin")

	err := adminColl.FindOne(ctx, bson.M{"_id": adminData.AdminId}).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(adminAuth.GetAdminResDto{
				Status:  false,
				Message: "admin not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(adminAuth.GetAdminResDto{
			Status:  false,
			Message: "Failed to fetch admin from MongoDB: " + err.Error(),
		})
	}

	adminRes := adminAuth.GetAdminRes{
		Id:        admin.Id,
		FirstName: admin.FirstName,
		LastName:  admin.LastName,
		Email:     admin.Email,
		Image:     admin.Image,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(adminAuth.GetAdminResDto{
		Status:  true,
		Message: "admin data retrieved successfully",
		Data:    adminRes,
	})
}
