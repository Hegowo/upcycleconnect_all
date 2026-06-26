package database

import (
	"errors"

	"upcycleconnect/backend/models"

	"gorm.io/gorm"
)

func SeedSubscriptionPlans(db *gorm.DB) {
	for _, seed := range models.DefaultPlanSeeds {
		var existing models.SubscriptionPlan
		err := db.Where("`key` = ?", seed.Key).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			db.Create(&seed)
			continue
		}
		if err != nil {
			continue
		}
		if existing.Label == "" {
			db.Model(&models.SubscriptionPlan{}).Where("`key` = ?", seed.Key).Updates(map[string]interface{}{
				"label":                  seed.Label,
				"is_default":             seed.IsDefault,
				"is_active":              seed.IsActive,
				"sort_order":             seed.SortOrder,
				"max_prestations":        seed.MaxPrestations,
				"max_projects_per_month": seed.MaxProjectsPerMonth,
				"max_campaigns":          seed.MaxCampaigns,
				"max_events_per_month":   seed.MaxEventsPerMonth,
				"feature_advanced_stats": seed.FeatureAdvancedStats,
				"feature_premium_stats":  seed.FeaturePremiumStats,
				"features_json":          seed.FeaturesJSON,
			})
		}
	}
}
