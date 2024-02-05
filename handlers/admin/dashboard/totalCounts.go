package dashboardHandlers

import (
	"careville_backend/database"
	"careville_backend/dto/admin/dashboard"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

var ctx = context.Background()

// @Summary Count All
// @Description Count the number of documents in various collections
// @Tags dashboard in admin pannel
// @Accept application/json
//
// @Param Authorization header	string true	"Authentication header"
// @Produce json
// @Success 200 {object} dashboard.GetAllCounts
// @Router /admin/total-counts [get]
func CountAll(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
	)

	hospitalFilter := bson.M{"role": "healthFacility", "facilityOrProfession": "hospClinic", "serviceStatus": "approved"}
	pharmacyFilter := bson.M{"role": "healthFacility", "facilityOrProfession": "pharmacy", "serviceStatus": "approved"}
	fitnessCenterFilter := bson.M{"role": "healthFacility", "facilityOrProfession": "fitnessCenter", "serviceStatus": "approved"}
	laboratoryFilter := bson.M{"role": "healthFacility", "facilityOrProfession": "laboratory", "serviceStatus": "approved"}
	doctorFilter := bson.M{"role": "healthProfessional", "facilityOrProfession": "doctor", "serviceStatus": "approved"}
	nurseFilter := bson.M{"role": "healthProfessional", "facilityOrProfession": "nurse", "serviceStatus": "approved"}
	medicalLabScientistFilter := bson.M{"role": "healthProfessional", "facilityOrProfession": "medicalLabScientist", "serviceStatus": "approved"}
	physiotherapistFilter := bson.M{"role": "healthProfessional", "facilityOrProfession": "physiotherapist", "serviceStatus": "approved"}

	hospitalCount, err1 := serviceColl.CountDocuments(ctx, hospitalFilter)
	pharmacyCount, err2 := serviceColl.CountDocuments(ctx, pharmacyFilter)
	fitnessCenterCount, err3 := serviceColl.CountDocuments(ctx, fitnessCenterFilter)
	laboratoryCount, err4 := serviceColl.CountDocuments(ctx, laboratoryFilter)
	doctorCount, err5 := serviceColl.CountDocuments(ctx, doctorFilter)
	nurseCount, err6 := serviceColl.CountDocuments(ctx, nurseFilter)
	medicalLabScientistCount, err7 := serviceColl.CountDocuments(ctx, medicalLabScientistFilter)
	physiotherapistCount, err8 := serviceColl.CountDocuments(ctx, physiotherapistFilter)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil || err8 != nil {
		return c.Status(500).JSON(dashboard.GetAllCounts{
			Status:  false,
			Message: "Error while counting documents",
		})
	}

	return c.Status(200).JSON(dashboard.GetAllCounts{
		Status:  true,
		Message: "Successfully counted documents in various collections.",
		Data: dashboard.DashboardCount{
			HealthProfessionals: dashboard.HealthProfessionals{
				DoctorCount:              doctorCount,
				MedicalLabScientistCount: medicalLabScientistCount,
				PhysiotherapistCount:     physiotherapistCount,
				NurseCount:               nurseCount,
			},
			HealthFacilitities: dashboard.HealthFacility{
				LaboratoryCount:    laboratoryCount,
				HospitalCount:      hospitalCount,
				PharmacyCount:      pharmacyCount,
				FitnessCenterCount: fitnessCenterCount,
			},
		},
	})
}
