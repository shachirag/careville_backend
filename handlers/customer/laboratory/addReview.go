package laboratory

import (
	"careville_backend/database"
	reviews "careville_backend/dto/customer/reviews"
	"careville_backend/entity"

	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Add Review for service
// @Description Add Review for service
// @Tags customer reviews
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param review body reviews.ReviewsReqDto true "Review data"
// @Produce json
// @Success 200 {object} reviews.ReviewsResDto
// @Router /customer/healthFacility/add-laboratory-review [post]
func AddLaboratoryReview(c *fiber.Ctx) error {
	var data reviews.ReviewsReqDto

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(reviews.ReviewsResDto{
			Status:  false,
			Message: "Invalid request format",
		})
	}

	customerColl := database.GetCollection("customer")

	var customer entity.CustomerEntity
	customerFilter := bson.M{"_id": data.CustomerId}
	err := customerColl.FindOne(ctx, customerFilter).Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(reviews.ReviewsResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	review := entity.ReviewEntity{
		Id: primitive.NewObjectID(),
		Customer: entity.Customer{
			Id:        data.CustomerId,
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Image:     customer.Image,
		},
		Description: data.Description,
		ServiceId:   data.ServiceId,
		Rating:      data.Rating,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	session, err := database.GetMongoClient().StartSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(reviews.ReviewsResDto{
			Status:  false,
			Message: "Failed to start session",
		})
	}
	defer session.EndSession(ctx)

	var newAverageRating float64
	var newTotalReviews int32

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		reviewColl := database.GetCollection("review")
		var existingReviews entity.ReviewEntity
		err1 := reviewColl.FindOne(sessCtx, bson.M{"serviceId": review.ServiceId,
			"customer.id": review.Customer.Id}).Decode(&existingReviews)
		if err1 != nil && err1 != mongo.ErrNoDocuments {
			return nil, err1
		}

		if err1 == mongo.ErrNoDocuments {
			if _, err := reviewColl.InsertOne(sessCtx, review); err != nil {
				return nil, err
			}
			var newAverageRating float64
			var newTotalReviews int32

			serviceColl := database.GetCollection("service")
			var service entity.ServiceEntity
			filter := bson.M{"_id": review.ServiceId}
			err = serviceColl.FindOne(sessCtx, filter).Decode(&service)
			if err != nil {
				return nil, err
			}

			if service.Laboratory != nil {
				newTotalReviews = service.Laboratory.Review.TotalReviews + 1
				newAverageRating = ((service.Laboratory.Review.AvgRating * float64(service.Laboratory.Review.TotalReviews)) +
					review.Rating) / float64(newTotalReviews)
			}

			update := bson.M{
				"$set": bson.M{
					"laboratory.review.totalReviews": newTotalReviews,
					"laboratory.review.avgRating":    newAverageRating,
				},
			}

			_, err := serviceColl.UpdateOne(
				sessCtx,
				bson.M{"_id": review.ServiceId},
				update,
			)
			if err != nil {
				return nil, err
			}
		} else {
			update := bson.M{
				"$set": bson.M{
					"description": review.Description,
					"rating":      review.Rating,
					"updatedAt":   time.Now().UTC(),
				},
			}
			_, err := reviewColl.UpdateOne(
				sessCtx,
				bson.M{"serviceId": review.ServiceId, "customer.id": review.Customer.Id},
				update,
			)
			if err != nil {
				return nil, err
			}

			serviceColl := database.GetCollection("service")
			var service entity.ServiceEntity
			err = serviceColl.FindOne(sessCtx, bson.M{"_id": review.ServiceId}).Decode(&service)
			if err != nil {
				return nil, err
			}

			if service.Laboratory != nil {
				newTotalReviews = service.Laboratory.Review.TotalReviews
				newAverageRating = ((service.Laboratory.Review.AvgRating * float64(service.Laboratory.Review.TotalReviews)) - existingReviews.Rating + review.Rating) / float64(newTotalReviews)
			}

			update = bson.M{
				"$set": bson.M{
					"laboratory.review.totalReviews": newTotalReviews,
					"laboratory.review.avgRating":    newAverageRating,
				},
			}
			_, err = serviceColl.UpdateOne(
				sessCtx,
				bson.M{"_id": review.ServiceId},
				update,
			)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(reviews.ReviewsResDto{
			Status:  false,
			Message: "Failed to insert review data or update totalReviews" + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(reviews.ReviewsResDto{
		Status:  true,
		Message: "Review added successfully",
	})
}
