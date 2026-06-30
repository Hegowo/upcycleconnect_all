package database

import (
	"upcycleconnect/backend/models"

	"gorm.io/gorm"
)

func SeedLegalDocuments(db *gorm.DB) {
	defaults := []models.LegalDocument{
		{
			Slug:    "mentions-legales",
			Title:   "Mentions légales",
			Content: "<h2>Mentions légales</h2><p>Éditeur du site : UpcycleConnect.</p><p>Ce contenu est à compléter depuis l'espace d'administration.</p>",
		},
		{
			Slug:    "cgu-cgv",
			Title:   "Conditions générales d'utilisation et de vente",
			Content: "<h2>Conditions générales</h2><p>Les présentes conditions régissent l'utilisation de la plateforme UpcycleConnect.</p><p>Ce contenu est à compléter depuis l'espace d'administration.</p>",
		},
		{
			Slug:    "confidentialite",
			Title:   "Politique de confidentialité",
			Content: "<h2>Politique de confidentialité</h2><p>UpcycleConnect s'engage à protéger les données personnelles de ses utilisateurs.</p><p>Ce contenu est à compléter depuis l'espace d'administration.</p>",
		},
	}
	for _, d := range defaults {
		var count int64
		db.Model(&models.LegalDocument{}).Where("slug = ?", d.Slug).Count(&count)
		if count == 0 {
			db.Create(&d)
		}
	}
}
