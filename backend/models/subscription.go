package models

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	UserID               uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	StripeCustomerID     string     `gorm:"size:100" json:"stripe_customer_id"`
	StripeSubscriptionID string     `gorm:"size:100;uniqueIndex" json:"stripe_subscription_id"`
	Plan                 string     `gorm:"size:20;default:basic" json:"plan"`
	Status               string     `gorm:"size:30;default:inactive" json:"status"`
	CurrentPeriodEnd     *time.Time `json:"current_period_end"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (Subscription) TableName() string { return "subscriptions" }

type SubscriptionPlan struct {
	Key         string
	Label       string
	AmountCents int64
	Features    []string
}

type SubscriptionPlanConfig struct {
	Key         string    `gorm:"size:20;primaryKey" json:"key"`
	AmountCents int64     `gorm:"not null" json:"amount_cents"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SubscriptionPlanConfig) TableName() string { return "subscription_plans" }

var PlanOrder = []string{"basic", "premium"}

var defaultSubscriptionPlans = map[string]SubscriptionPlan{
	"basic": {
		Key:         "basic",
		Label:       "Basic",
		AmountCents: 1500,
		Features: []string{
			"Tableaux de bord avancés (chiffre d'affaires, réservations)",
			"Statistiques sur les matériaux disponibles",
		},
	},
	"premium": {
		Key:         "premium",
		Label:       "Premium",
		AmountCents: 3000,
		Features: []string{
			"Tout ce qui est inclus dans Basic",
			"Analyse d'impact écologique détaillée",
			"Alertes priorisées pour la collecte",
			"Gestion de campagnes publicitaires (mise en avant)",
		},
	},
}

var (
	planMu    sync.RWMutex
	planCache = clonePlans()
)

func clonePlans() map[string]SubscriptionPlan {
	out := make(map[string]SubscriptionPlan, len(defaultSubscriptionPlans))
	for k, v := range defaultSubscriptionPlans {
		out[k] = v
	}
	return out
}

func LoadSubscriptionPlans(db *gorm.DB) {
	plans := clonePlans()
	var rows []SubscriptionPlanConfig
	if err := db.Find(&rows).Error; err == nil {
		for _, r := range rows {
			if p, ok := plans[r.Key]; ok {
				p.AmountCents = r.AmountCents
				plans[r.Key] = p
			}
		}
	}
	planMu.Lock()
	planCache = plans
	planMu.Unlock()
}

func Plan(key string) (SubscriptionPlan, bool) {
	planMu.RLock()
	defer planMu.RUnlock()
	p, ok := planCache[key]
	return p, ok
}

func PlansList() []SubscriptionPlan {
	planMu.RLock()
	defer planMu.RUnlock()
	out := make([]SubscriptionPlan, 0, len(planCache))
	for _, key := range PlanOrder {
		if p, ok := planCache[key]; ok {
			out = append(out, p)
		}
	}
	return out
}

type SubscriptionResponse struct {
	ID               uint    `json:"id"`
	Plan             string  `json:"plan"`
	PlanLabel        string  `json:"plan_label"`
	Status           string  `json:"status"`
	CurrentPeriodEnd *string `json:"current_period_end"`
	AmountCents      int64   `json:"amount_cents"`
}

func ToSubscriptionResponse(s *Subscription) SubscriptionResponse {
	var end *string
	if s.CurrentPeriodEnd != nil {
		v := s.CurrentPeriodEnd.UTC().Format(time.RFC3339)
		end = &v
	}
	plan, _ := Plan(s.Plan)
	return SubscriptionResponse{
		ID:               s.ID,
		Plan:             s.Plan,
		PlanLabel:        plan.Label,
		Status:           s.Status,
		CurrentPeriodEnd: end,
		AmountCents:      plan.AmountCents,
	}
}
