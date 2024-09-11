package customerAuth

import (
	"careville_backend/database"
	customerAuth "careville_backend/dto/customer/customerAuth"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"
	"careville_backend/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// @Param newCustomerImage formData file false "customer profile image"
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

	formFile, err := c.FormFile("newCustomerImage")
	var imageURL string
	if err != nil {
		imageURL = data.OldImage
	} else {
		file, err := formFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.UpdateCustomerResDto{
				Status:  false,
				Message: "Failed to open image file: " + err.Error(),
			})
		}
		defer file.Close()

		id := primitive.NewObjectID()
		fileName := fmt.Sprintf("customer/%v-profilepic%s", id.Hex(), formFile.Filename)

		imageURL, err = utils.UploadToS3(fileName, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.UpdateCustomerResDto{
				Status:  false,
				Message: "Failed to upload image to S3: " + err.Error(),
			})
		}
	}

	update := bson.M{"$set": bson.M{
		"firstName": data.FirstName,
		"lastName":  data.LastName,
		"phoneNumber": bson.M{
			"dialCode":    data.DialCode,
			"number":      data.PhoneNumber,
			"countryCode": data.CountryCode,
		},
		"image": imageURL,
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

	session, err := database.GetMongoClient().StartSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.UpdateCustomerResDto{
			Status:  false,
			Message: "Failed to start session",
		})
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		updateRes, err := customerColl.UpdateOne(sessCtx, filter, update, opts)
		if err != nil {
			return nil, err
		}

		if updateRes.MatchedCount == 0 {
			return nil, mongo.ErrNoDocuments
		}

		appointmentUpdate := bson.M{"$set": bson.M{
			"customer.firstName": data.FirstName,
			"customer.lastName":  data.LastName,
			"customer.image":     imageURL,
			"customer.phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"age": data.Age,
		}}

		_, err = database.GetCollection("appointment").UpdateMany(sessCtx, bson.M{"customer.id": customerData.CustomerId}, appointmentUpdate)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.UpdateCustomerResDto{
			Status:  false,
			Message: "Failed to update appointment data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(customerAuth.UpdateCustomerResDto{
		Status:  true,
		Message: "Customer data updated successfully",
	})
}
