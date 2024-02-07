package physiotherapist

// import (
// 	"careville_backend/database"
// 	physiotherapist "careville_backend/dto/customer/physiotherapist"
// 	"careville_backend/entity"
// 	"context"
// 	"math"
// 	"strconv"

// 	"github.com/gofiber/fiber/v2"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// var ctx = context.Background()

// // @Summary Fetch physiotherapist With Filters
// // @Description Fetch physiotherapist With Filters
// // @Tags customer physiotherapist
// // @Accept application/json
// //
// //	@Param Authorization header	string true	"Authentication header"
// //
// // @Param page query int false "Page no. to fetch the products for 1"
// // @Param perPage query int false "Limit of products to fetch is 15"
// // @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// // @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// // @Param search query string false "Filter nurse by search"
// // @Produce json
// // @Success 200 {object} physiotherapist.GetPhysiotherapistPaginationRes
// // @Router /customer/healthProfessional/get-physiotherapists [get]
// func FetchPhysiotherapistWithPagination(c *fiber.Ctx) error {

// 	page, _ := strconv.Atoi(c.Query("page", "1"))
// 	limit, _ := strconv.Atoi(c.Query("limit", "15"))

// 	var lat, long float64
// 	latParam := c.Query("lat")
// 	longParam := c.Query("long")
// 	var err error

// 	if latParam != "" && longParam != "" {
// 		lat, err = strconv.ParseFloat(latParam, 64)
// 		if err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.GetPhysiotherapistPaginationRes{
// 				Status:  false,
// 				Message: "Invalid latitude format",
// 			})
// 		}

// 		long, err = strconv.ParseFloat(longParam, 64)
// 		if err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.GetPhysiotherapistPaginationRes{
// 				Status:  false,
// 				Message: "Invalid longitude format",
// 			})
// 		}
// 	}

// 	searchTitle := c.Query("search", "")

// 	serviceColl := database.GetCollection("service")

// 	filter := bson.M{
// 		"role":                 "healthProfessional",
// 		"facilityOrProfession": "physiotherapist",
// 	}

// 	if latParam != "" && longParam != "" {
// 		filter["physiotherapist.information.address"] = bson.M{
// 			"$nearSphere": bson.M{
// 				"$geometry": bson.M{
// 					"type":        "Point",
// 					"coordinates": []float64{long, lat},
// 				},
// 				"$maxDistance": 20000,
// 			},
// 		}
// 	}

// 	if searchTitle != "" {
// 		filter["physiotherapist.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
// 	}

// 	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

// 	skip := (page - 1) * limit

// 	projection := bson.M{
// 		"physiotherapist.information.name":  1,
// 		"physiotherapist.information.image": 1,
// 		"physiotherapist.information.id":    1,
// 		"avgRating":                         1,
// 	}

// 	findOptions := options.Find().SetProjection(projection).SetSkip(int64(skip)).SetLimit(int64(limit))

// 	cursor, err := serviceColl.Find(ctx, filter, findOptions, sortOptions)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return c.Status(fiber.StatusNotFound).JSON(physiotherapist.GetPhysiotherapistPaginationRes{
// 				Status:  false,
// 				Message: "physiotherapist not found",
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistPaginationRes{
// 			Status:  false,
// 			Message: "Failed to fetch physiotherapist from MongoDB: " + err.Error(),
// 		})
// 	}
// 	defer cursor.Close(ctx)

// 	response := physiotherapist.PhysiotherapistPaginationResponse{
// 		Total:              0,
// 		PerPage:            limit,
// 		CurrentPage:        page,
// 		TotalPages:         0,
// 		PhysiotherapistRes: []physiotherapist.GetPhysiotherapistRes{},
// 	}

// 	for cursor.Next(ctx) {
// 		var service entity.ServiceEntity
// 		err := cursor.Decode(&service)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistPaginationRes{
// 				Status:  false,
// 				Message: "Failed to decode physiotherapist data: " + err.Error(),
// 			})
// 		}

// 		// Check if hospClinic is not nil before accessing its properties
// 		if service.Physiotherapist != nil {
// 			nurseRes := physiotherapist.GetPhysiotherapistRes{
// 				Id:        service.Id,
// 				Image:     service.Physiotherapist.Information.Image,
// 				Name:      service.Physiotherapist.Information.Name,
// 				AvgRating: service.AvgRating,
// 			}

// 			response.PhysiotherapistRes = append(response.PhysiotherapistRes, nurseRes)
// 		}
// 	}

// 	totalCount, err := serviceColl.CountDocuments(ctx, bson.M{
// 		"role":                             "healthProfessional",
// 		"facilityOrProfession":             "physiotherapist",
// 		"physiotherapist.information.name": bson.M{"$regex": searchTitle, "$options": "i"},
// 	})
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistPaginationRes{
// 			Status:  false,
// 			Message: "Failed to count physiotherapist: " + err.Error(),
// 		})
// 	}

// 	response.Total = int(totalCount)
// 	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

// 	finalResponse := physiotherapist.GetPhysiotherapistPaginationRes{
// 		Status:  true,
// 		Message: "Sucessfully fetched data",
// 		Data:    response,
// 	}
// 	return c.Status(fiber.StatusOK).JSON(finalResponse)
// }
