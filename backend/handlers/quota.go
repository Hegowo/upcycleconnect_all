package handlers

import (
	"fmt"
	"net/http"
	"time"

	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func effectivePlan(db *gorm.DB, userID uint) models.SubscriptionPlan {
	var sub models.Subscription
	if err := db.Where("user_id = ? AND status = ?", userID, "active").First(&sub).Error; err == nil {
		if p, ok := models.Plan(sub.Plan); ok {
			return p
		}
	}
	return models.DefaultPlan()
}

func startOfMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}

func quotaExceeded(c *gin.Context, limit *int, used int64, resource string) bool {
	if limit == nil {
		return false
	}
	if used >= int64(*limit) {
		c.JSON(http.StatusForbidden, gin.H{
			"message":          fmt.Sprintf("Limite de votre abonnement atteinte pour %s (%d max). Passez à un abonnement supérieur pour en ajouter davantage.", resource, *limit),
			"upgrade_required": true,
			"limit":            *limit,
			"resource":         resource,
		})
		return true
	}
	return false
}

func enforcePrestationQuota(c *gin.Context, db *gorm.DB, userID uint) bool {
	plan := effectivePlan(db, userID)
	var used int64
	db.Model(&models.Prestation{}).
		Where("provider_id = ? AND deleted_at IS NULL AND status IN ?", userID, []string{"pending", "published", "suspended"}).
		Count(&used)
	return quotaExceeded(c, plan.MaxPrestations, used, "les prestations publiées")
}

func enforceProjectQuota(c *gin.Context, db *gorm.DB, userID uint) bool {
	plan := effectivePlan(db, userID)
	var used int64
	db.Model(&models.UpcyclingProject{}).
		Where("provider_id = ? AND created_at >= ?", userID, startOfMonth()).
		Count(&used)
	return quotaExceeded(c, plan.MaxProjectsPerMonth, used, "les projets ce mois-ci")
}

func enforceCampaignQuota(c *gin.Context, db *gorm.DB, userID uint) bool {
	plan := effectivePlan(db, userID)
	var used int64
	db.Model(&models.Campaign{}).
		Where("provider_id = ? AND deleted_at IS NULL AND status IN ?", userID, []string{"pending", "active"}).
		Count(&used)
	return quotaExceeded(c, plan.MaxCampaigns, used, "les campagnes")
}

func enforceEventQuota(c *gin.Context, db *gorm.DB, userID uint) bool {
	plan := effectivePlan(db, userID)
	var used int64
	db.Model(&models.Event{}).
		Where("created_by = ? AND deleted_at IS NULL AND created_at >= ?", userID, startOfMonth()).
		Count(&used)
	return quotaExceeded(c, plan.MaxEventsPerMonth, used, "les événements ce mois-ci")
}
