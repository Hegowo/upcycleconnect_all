package models

import "time"

type LegalDocument struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Slug      string    `gorm:"uniqueIndex;size:64" json:"slug"`
	Title     string    `gorm:"size:160" json:"title"`
	Content   string    `gorm:"type:longtext" json:"content"`
	UpdatedBy *uint     `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (LegalDocument) TableName() string {
	return "legal_documents"
}

type LegalDocumentResponse struct {
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updated_at"`
}

func ToLegalDocumentResponse(d *LegalDocument) LegalDocumentResponse {
	return LegalDocumentResponse{
		Slug:      d.Slug,
		Title:     d.Title,
		Content:   d.Content,
		UpdatedAt: d.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
