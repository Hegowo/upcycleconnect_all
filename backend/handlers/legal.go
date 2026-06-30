package handlers

import (
	"net/http"
	"strings"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LegalHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

var legalSlugs = map[string]bool{
	"mentions-legales": true,
	"cgu-cgv":          true,
	"confidentialite":  true,
}

func (h *LegalHandler) Show(c *gin.Context) {
	slug := c.Param("slug")
	if !legalSlugs[slug] {
		c.JSON(http.StatusNotFound, gin.H{"message": "Document introuvable."})
		return
	}
	var doc models.LegalDocument
	if err := h.DB.Where("slug = ?", slug).First(&doc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Document introuvable."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": models.ToLegalDocumentResponse(&doc)})
}

func (h *LegalHandler) AdminIndex(c *gin.Context) {
	var docs []models.LegalDocument
	h.DB.Order("slug ASC").Find(&docs)
	out := make([]models.LegalDocumentResponse, 0, len(docs))
	for i := range docs {
		out = append(out, models.ToLegalDocumentResponse(&docs[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *LegalHandler) AdminUpdate(c *gin.Context) {
	slug := c.Param("slug")
	if !legalSlugs[slug] {
		c.JSON(http.StatusNotFound, gin.H{"message": "Document introuvable."})
		return
	}
	var doc models.LegalDocument
	if err := h.DB.Where("slug = ?", slug).First(&doc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Document introuvable."})
		return
	}
	var req struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}
	if req.Title != nil && strings.TrimSpace(*req.Title) != "" {
		doc.Title = strings.TrimSpace(*req.Title)
	}
	if req.Content != nil {
		doc.Content = *req.Content
	}
	if user := middleware.GetAuthUser(c); user != nil {
		doc.UpdatedBy = &user.ID
	}
	if err := h.DB.Save(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement."})
		return
	}
	if h.Audit != nil {
		h.Audit.Log(c, "legal.updated", "LegalDocument", &doc.ID, nil, map[string]interface{}{"slug": doc.Slug})
	}
	c.JSON(http.StatusOK, gin.H{"data": models.ToLegalDocumentResponse(&doc)})
}
