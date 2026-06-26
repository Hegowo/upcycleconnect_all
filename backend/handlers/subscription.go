package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	DB     *gorm.DB
	Stripe *services.StripeService
	Audit  *services.AuditService
}

func (h *SubscriptionHandler) Plans(c *gin.Context) {
	plans := []gin.H{}
	for _, plan := range models.PlansList() {
		plans = append(plans, planJSON(plan))
	}
	c.JSON(http.StatusOK, gin.H{"data": plans})
}

func (h *SubscriptionHandler) MySubscription(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	var sub models.Subscription
	if err := h.DB.Where("user_id = ?", user.ID).First(&sub).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"subscription": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"subscription": models.ToSubscriptionResponse(&sub)})
}

func (h *SubscriptionHandler) Checkout(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}

	var req struct {
		Plan string `json:"plan" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Plan requis"})
		return
	}

	var existing models.Subscription
	if h.DB.Where("user_id = ? AND status = ?", user.ID, "active").First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Tu as déjà un abonnement actif."})
		return
	}

	plan, ok := models.Plan(req.Plan)
	if !ok || !plan.IsActive || plan.IsDefault {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Plan invalide ou indisponible."})
		return
	}

	sess, err := h.Stripe.CreateSubscriptionCheckout(services.SubscriptionCheckoutParams{
		UserID:      user.ID,
		UserEmail:   user.Email,
		Plan:        req.Plan,
		Label:       plan.Label,
		AmountCents: plan.AmountCents,
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"checkout_url": sess.URL, "session_id": sess.ID})
}

func (h *SubscriptionHandler) Cancel(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	var sub models.Subscription
	if err := h.DB.Where("user_id = ? AND status = ?", user.ID, "active").First(&sub).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Aucun abonnement actif"})
		return
	}
	if err := h.Stripe.CancelSubscription(sub.StripeSubscriptionID); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Impossible d'annuler via Stripe : " + err.Error()})
		return
	}
	h.DB.Model(&sub).Updates(map[string]interface{}{"status": "cancelled"})
	h.Audit.Log(c, "subscription.cancelled", "Subscription", &sub.ID, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Abonnement annulé."})
}

func planJSON(p models.SubscriptionPlan) gin.H {
	return gin.H{
		"key":                    p.Key,
		"label":                  p.Label,
		"amount_cents":           p.AmountCents,
		"is_default":             p.IsDefault,
		"is_active":              p.IsActive,
		"sort_order":             p.SortOrder,
		"max_prestations":        p.MaxPrestations,
		"max_projects_per_month": p.MaxProjectsPerMonth,
		"max_campaigns":          p.MaxCampaigns,
		"max_events_per_month":   p.MaxEventsPerMonth,
		"feature_advanced_stats": p.FeatureAdvancedStats,
		"feature_premium_stats":  p.FeaturePremiumStats,
		"features":               p.Features(),
	}
}

type planRequest struct {
	Label                string    `json:"label"`
	AmountCents          *int64    `json:"amount_cents"`
	IsActive             *bool     `json:"is_active"`
	IsDefault            *bool     `json:"is_default"`
	SortOrder            *int      `json:"sort_order"`
	MaxPrestations       *int      `json:"max_prestations"`
	MaxProjectsPerMonth  *int      `json:"max_projects_per_month"`
	MaxCampaigns         *int      `json:"max_campaigns"`
	MaxEventsPerMonth    *int      `json:"max_events_per_month"`
	FeatureAdvancedStats *bool     `json:"feature_advanced_stats"`
	FeaturePremiumStats  *bool     `json:"feature_premium_stats"`
	Features             *[]string `json:"features"`
}

func applyPlan(p *models.SubscriptionPlan, req planRequest) {
	p.Label = strings.TrimSpace(req.Label)
	if req.AmountCents != nil {
		p.AmountCents = *req.AmountCents
	}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		p.SortOrder = *req.SortOrder
	}
	p.MaxPrestations = req.MaxPrestations
	p.MaxProjectsPerMonth = req.MaxProjectsPerMonth
	p.MaxCampaigns = req.MaxCampaigns
	p.MaxEventsPerMonth = req.MaxEventsPerMonth
	if req.FeatureAdvancedStats != nil {
		p.FeatureAdvancedStats = *req.FeatureAdvancedStats
	}
	if req.FeaturePremiumStats != nil {
		p.FeaturePremiumStats = *req.FeaturePremiumStats
	}
	if req.Features != nil {
		b, _ := json.Marshal(*req.Features)
		p.FeaturesJSON = string(b)
	}
}

func (h *SubscriptionHandler) uniquePlanKey(base string) string {
	if base == "" {
		base = "plan"
	}
	key := base
	for i := 2; ; i++ {
		var cnt int64
		h.DB.Model(&models.SubscriptionPlan{}).Where("`key` = ?", key).Count(&cnt)
		if cnt == 0 {
			return key
		}
		key = fmt.Sprintf("%s-%d", base, i)
	}
}

func (h *SubscriptionHandler) AdminPlans(c *gin.Context) {
	plans := []gin.H{}
	for _, plan := range models.AllPlans() {
		var subs int64
		h.DB.Model(&models.Subscription{}).Where("plan = ? AND status = ?", plan.Key, "active").Count(&subs)
		j := planJSON(plan)
		j["active_subscribers"] = subs
		plans = append(plans, j)
	}
	c.JSON(http.StatusOK, gin.H{"data": plans})
}

func (h *SubscriptionHandler) AdminCreatePlan(c *gin.Context) {
	var req planRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Label) == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Le nom du plan est requis."})
		return
	}
	plan := models.SubscriptionPlan{Key: h.uniquePlanKey(slugify(req.Label)), IsActive: true}
	applyPlan(&plan, req)
	if err := h.DB.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du plan."})
		return
	}
	if req.IsDefault != nil && *req.IsDefault {
		h.setDefaultPlan(plan.Key)
	}
	models.LoadSubscriptionPlans(h.DB)
	h.Audit.Log(c, "subscription_plan.created", "SubscriptionPlan", nil, nil, map[string]interface{}{"key": plan.Key})
	updated, _ := models.Plan(plan.Key)
	c.JSON(http.StatusCreated, planJSON(updated))
}

func (h *SubscriptionHandler) AdminUpdatePlan(c *gin.Context) {
	key := c.Param("key")
	var plan models.SubscriptionPlan
	if err := h.DB.Where("`key` = ?", key).First(&plan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Plan introuvable"})
		return
	}
	var req planRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Label) == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Le nom du plan est requis."})
		return
	}
	applyPlan(&plan, req)
	if err := h.DB.Save(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement."})
		return
	}
	if req.IsDefault != nil && *req.IsDefault {
		h.setDefaultPlan(plan.Key)
	}
	models.LoadSubscriptionPlans(h.DB)
	h.Audit.Log(c, "subscription_plan.updated", "SubscriptionPlan", nil, nil, map[string]interface{}{"key": key})
	updated, _ := models.Plan(key)
	c.JSON(http.StatusOK, planJSON(updated))
}

func (h *SubscriptionHandler) AdminDeletePlan(c *gin.Context) {
	key := c.Param("key")
	var plan models.SubscriptionPlan
	if err := h.DB.Where("`key` = ?", key).First(&plan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Plan introuvable"})
		return
	}
	if plan.IsDefault {
		c.JSON(http.StatusConflict, gin.H{"message": "Impossible de supprimer le plan par défaut. Définissez-en un autre par défaut d'abord."})
		return
	}
	var subs int64
	h.DB.Model(&models.Subscription{}).Where("plan = ? AND status = ?", key, "active").Count(&subs)
	if subs > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": fmt.Sprintf("%d abonné(s) actif(s) sur ce plan. Désactivez-le plutôt que de le supprimer.", subs)})
		return
	}
	h.DB.Where("`key` = ?", key).Delete(&models.SubscriptionPlan{})
	models.LoadSubscriptionPlans(h.DB)
	h.Audit.Log(c, "subscription_plan.deleted", "SubscriptionPlan", nil, map[string]interface{}{"key": key}, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Plan supprimé."})
}

func (h *SubscriptionHandler) setDefaultPlan(key string) {
	h.DB.Model(&models.SubscriptionPlan{}).Where("`key` <> ?", key).Update("is_default", false)
	h.DB.Model(&models.SubscriptionPlan{}).Where("`key` = ?", key).Update("is_default", true)
}

func (h *SubscriptionHandler) AdminIndex(c *gin.Context) {
	var subs []models.Subscription
	h.DB.Preload("User").Order("created_at DESC").Find(&subs)
	resp := make([]gin.H, 0, len(subs))
	for i := range subs {
		r := models.ToSubscriptionResponse(&subs[i])
		extra := gin.H{
			"id":                 r.ID,
			"plan":               r.Plan,
			"plan_label":         r.PlanLabel,
			"status":             r.Status,
			"current_period_end": r.CurrentPeriodEnd,
			"amount_cents":       r.AmountCents,
		}
		if subs[i].User != nil {
			extra["user_email"] = subs[i].User.Email
			extra["user_name"] = subs[i].User.FirstName + " " + subs[i].User.LastName
		}
		resp = append(resp, extra)
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
