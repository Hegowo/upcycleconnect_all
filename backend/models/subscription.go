package models

import (
	"encoding/json"
	"sort"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	UserID               uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	StripeCustomerID     string     `gorm:"size:100" json:"stripe_customer_id"`
	StripeSubscriptionID string     `gorm:"size:100;uniqueIndex" json:"stripe_subscription_id"`
	Plan                 string     `gorm:"size:40;default:basic" json:"plan"`
	Status               string     `gorm:"size:30;default:inactive" json:"status"`
	StartedAt            *time.Time `json:"started_at"`
	CurrentPeriodEnd     *time.Time `json:"current_period_end"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (Subscription) TableName() string { return "subscriptions" }

type SubscriptionPlan struct {
	Key         string `gorm:"size:40;primaryKey" json:"key"`
	Label       string `gorm:"size:80;default:''" json:"label"`
	AmountCents int64  `gorm:"not null;default:0" json:"amount_cents"`
	IsDefault   bool   `gorm:"not null;default:false" json:"is_default"`
	IsActive    bool   `gorm:"not null;default:true" json:"is_active"`
	SortOrder   int    `gorm:"not null;default:0" json:"sort_order"`

	MaxPrestations      *int `json:"max_prestations"`
	MaxProjectsPerMonth *int `json:"max_projects_per_month"`
	MaxCampaigns        *int `json:"max_campaigns"`
	MaxEventsPerMonth   *int `json:"max_events_per_month"`

	FeatureAdvancedStats bool `gorm:"not null;default:false" json:"feature_advanced_stats"`
	FeaturePremiumStats  bool `gorm:"not null;default:false" json:"feature_premium_stats"`

	FeaturesJSON string    `gorm:"type:text" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (SubscriptionPlan) TableName() string { return "subscription_plans" }

func (p SubscriptionPlan) Features() []string {
	if p.FeaturesJSON == "" {
		return []string{}
	}
	var f []string
	if json.Unmarshal([]byte(p.FeaturesJSON), &f) != nil {
		return []string{}
	}
	return f
}

func intPtr(v int) *int { return &v }

var DefaultPlanSeeds = []SubscriptionPlan{
	{
		Key: "free", Label: "Gratuit", AmountCents: 0, IsDefault: true, IsActive: true, SortOrder: 0,
		MaxPrestations: intPtr(1), MaxProjectsPerMonth: intPtr(1), MaxCampaigns: intPtr(0), MaxEventsPerMonth: intPtr(1),
		FeatureAdvancedStats: false, FeaturePremiumStats: false,
		FeaturesJSON: mustJSON([]string{"1 prestation publiée", "1 projet par mois", "1 événement par mois"}),
	},
	{
		Key: "basic", Label: "Basic", AmountCents: 1500, IsDefault: false, IsActive: true, SortOrder: 1,
		MaxPrestations: intPtr(10), MaxProjectsPerMonth: intPtr(5), MaxCampaigns: intPtr(0), MaxEventsPerMonth: intPtr(5),
		FeatureAdvancedStats: true, FeaturePremiumStats: false,
		FeaturesJSON: mustJSON([]string{
			"Jusqu'à 10 prestations publiées",
			"5 projets et 5 événements par mois",
			"Tableaux de bord avancés (chiffre d'affaires, réservations)",
			"Statistiques sur les matériaux disponibles",
		}),
	},
	{
		Key: "premium", Label: "Premium", AmountCents: 3000, IsDefault: false, IsActive: true, SortOrder: 2,
		MaxPrestations: nil, MaxProjectsPerMonth: nil, MaxCampaigns: intPtr(5), MaxEventsPerMonth: nil,
		FeatureAdvancedStats: true, FeaturePremiumStats: true,
		FeaturesJSON: mustJSON([]string{
			"Prestations, projets et événements illimités",
			"Jusqu'à 5 campagnes publicitaires",
			"Analyse d'impact écologique détaillée",
			"Alertes priorisées pour la collecte",
		}),
	},
}

func mustJSON(v []string) string {
	b, _ := json.Marshal(v)
	return string(b)
}

var (
	planMu    sync.RWMutex
	planCache = map[string]SubscriptionPlan{}
)

func LoadSubscriptionPlans(db *gorm.DB) {
	var rows []SubscriptionPlan
	if err := db.Find(&rows).Error; err != nil {
		return
	}
	m := make(map[string]SubscriptionPlan, len(rows))
	for _, r := range rows {
		m[r.Key] = r
	}
	planMu.Lock()
	planCache = m
	planMu.Unlock()
}

func Plan(key string) (SubscriptionPlan, bool) {
	planMu.RLock()
	defer planMu.RUnlock()
	p, ok := planCache[key]
	return p, ok
}

func sortedPlans(activeOnly bool) []SubscriptionPlan {
	planMu.RLock()
	out := make([]SubscriptionPlan, 0, len(planCache))
	for _, p := range planCache {
		if activeOnly && !p.IsActive {
			continue
		}
		out = append(out, p)
	}
	planMu.RUnlock()
	sort.Slice(out, func(i, j int) bool {
		if out[i].SortOrder != out[j].SortOrder {
			return out[i].SortOrder < out[j].SortOrder
		}
		return out[i].AmountCents < out[j].AmountCents
	})
	return out
}

func PlansList() []SubscriptionPlan {
	out := make([]SubscriptionPlan, 0)
	for _, p := range sortedPlans(true) {
		if p.IsDefault {
			continue
		}
		out = append(out, p)
	}
	return out
}

func AllPlans() []SubscriptionPlan { return sortedPlans(false) }

func DefaultPlan() SubscriptionPlan {
	planMu.RLock()
	defer planMu.RUnlock()
	for _, p := range planCache {
		if p.IsDefault {
			return p
		}
	}
	return SubscriptionPlan{Key: "free", Label: "Gratuit"}
}

type SubscriptionResponse struct {
	ID               uint    `json:"id"`
	Plan             string  `json:"plan"`
	PlanLabel        string  `json:"plan_label"`
	Status           string  `json:"status"`
	StartedAt        *string `json:"started_at"`
	CurrentPeriodEnd *string `json:"current_period_end"`
	AmountCents      int64   `json:"amount_cents"`
}

func ToSubscriptionResponse(s *Subscription) SubscriptionResponse {
	fmtTime := func(t *time.Time) *string {
		if t == nil {
			return nil
		}
		v := t.UTC().Format(time.RFC3339)
		return &v
	}
	started := s.StartedAt
	if started == nil {
		started = &s.CreatedAt
	}
	plan, _ := Plan(s.Plan)
	return SubscriptionResponse{
		ID:               s.ID,
		Plan:             s.Plan,
		PlanLabel:        plan.Label,
		Status:           s.Status,
		StartedAt:        fmtTime(started),
		CurrentPeriodEnd: fmtTime(s.CurrentPeriodEnd),
		AmountCents:      plan.AmountCents,
	}
}
