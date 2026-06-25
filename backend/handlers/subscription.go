package handlers

import (
	"net/http"

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
		plans = append(plans, gin.H{
			"key":          plan.Key,
			"label":        plan.Label,
			"amount_cents": plan.AmountCents,
			"features":     plan.Features,
		})
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
		Plan string `json:"plan" binding:"required,oneof=basic premium"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Plan invalide (basic ou premium)"})
		return
	}

	var existing models.Subscription
	if h.DB.Where("user_id = ? AND status = ?", user.ID, "active").First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Tu as déjà un abonnement actif."})
		return
	}

	plan, ok := models.Plan(req.Plan)
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Plan inconnu"})
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

func (h *SubscriptionHandler) AdminPlans(c *gin.Context) {
	plans := []gin.H{}
	for _, plan := range models.PlansList() {
		plans = append(plans, gin.H{
			"key":          plan.Key,
			"label":        plan.Label,
			"amount_cents": plan.AmountCents,
			"features":     plan.Features,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": plans})
}

func (h *SubscriptionHandler) AdminUpdatePlan(c *gin.Context) {
	key := c.Param("key")
	if _, ok := models.Plan(key); !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "Plan inconnu"})
		return
	}
	var req struct {
		AmountCents *int64 `json:"amount_cents"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.AmountCents == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Montant requis (amount_cents)."})
		return
	}
	if *req.AmountCents < 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Le montant doit être positif."})
		return
	}

	var cfg models.SubscriptionPlanConfig
	if err := h.DB.Where(models.SubscriptionPlanConfig{Key: key}).
		Assign(models.SubscriptionPlanConfig{AmountCents: *req.AmountCents}).
		FirstOrCreate(&cfg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement."})
		return
	}
	models.LoadSubscriptionPlans(h.DB)

	h.Audit.Log(c, "subscription_plan.updated", "SubscriptionPlan", nil, nil, map[string]interface{}{
		"key":          key,
		"amount_cents": *req.AmountCents,
	})

	plan, _ := models.Plan(key)
	c.JSON(http.StatusOK, gin.H{
		"key":          plan.Key,
		"label":        plan.Label,
		"amount_cents": plan.AmountCents,
		"features":     plan.Features,
	})
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
