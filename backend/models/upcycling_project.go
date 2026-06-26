package models

import (
	"time"

	"gorm.io/gorm"
)

type UpcyclingProject struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	ProviderID  uint    `gorm:"index;not null" json:"provider_id"`
	Title       string  `gorm:"size:200;not null" json:"title"`
	Description string  `gorm:"type:text" json:"description"`
	Category    string  `gorm:"size:80" json:"category"`
	CoverImage  *string `gorm:"size:500" json:"cover_image"`

	ImpactLabel *string        `gorm:"size:200" json:"impact_label"`
	Status      string         `gorm:"size:20;default:in_progress;index" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Provider *User          `gorm:"foreignKey:ProviderID" json:"-"`
	Steps    []ProjectStep  `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"steps,omitempty"`
	Images   []ProjectImage `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}

func (UpcyclingProject) TableName() string { return "upcycling_projects" }

type ProjectImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"index;not null" json:"project_id"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	Position  int       `gorm:"not null;default:0" json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

func (ProjectImage) TableName() string { return "project_images" }

type ProjectImageResponse struct {
	ID  uint   `json:"id"`
	URL string `json:"url"`
}

type ProjectStep struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ProjectID   uint       `gorm:"index;not null" json:"project_id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	ImageURL    *string    `gorm:"size:500" json:"image_url"`
	StepOrder   int        `gorm:"not null;default:0" json:"step_order"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (ProjectStep) TableName() string { return "project_steps" }

type ProjectResponse struct {
	ID           uint          `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Category     string        `json:"category"`
	CoverImage   *string                `json:"cover_image"`
	ImpactLabel  *string                `json:"impact_label"`
	Status       string                 `json:"status"`
	ProviderID   uint                   `json:"provider_id"`
	ProviderName string                 `json:"provider_name"`
	Steps        []ProjectStep          `json:"steps"`
	StepCount    int                    `json:"step_count"`
	Images       []ProjectImageResponse `json:"images"`
	CreatedAt    string                 `json:"created_at"`
}

func ToProjectResponse(p *UpcyclingProject, withSteps bool) ProjectResponse {
	name := ""
	if p.Provider != nil {
		name = p.Provider.FirstName + " " + p.Provider.LastName
	}
	steps := p.Steps
	if !withSteps {
		steps = nil
	}
	images := make([]ProjectImageResponse, 0, len(p.Images))
	for _, img := range p.Images {
		images = append(images, ProjectImageResponse{ID: img.ID, URL: img.URL})
	}
	return ProjectResponse{
		ID:           p.ID,
		Title:        p.Title,
		Description:  p.Description,
		Category:     p.Category,
		CoverImage:   p.CoverImage,
		ImpactLabel:  p.ImpactLabel,
		Status:       p.Status,
		ProviderID:   p.ProviderID,
		ProviderName: name,
		Steps:        steps,
		StepCount:    len(p.Steps),
		Images:       images,
		CreatedAt:    p.CreatedAt.UTC().Format(time.RFC3339),
	}
}
