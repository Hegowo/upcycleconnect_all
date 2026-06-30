package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TicketHandler struct {
	DB            *gorm.DB
	Audit         *services.AuditService
	Notifications *services.NotificationService
}

var ticketStatuses = map[string]bool{
	"open":        true,
	"in_progress": true,
	"resolved":    true,
	"closed":      true,
}

var ticketRefTypes = map[string]bool{
	"prestation": true,
	"event":      true,
	"project":    true,
}

func (h *TicketHandler) resolveRef(refType *string, refID *uint) *models.TicketRefResponse {
	if refType == nil || refID == nil || !ticketRefTypes[*refType] {
		return nil
	}
	out := &models.TicketRefResponse{Type: *refType, ID: *refID}
	switch *refType {
	case "prestation":
		var p models.Prestation
		if err := h.DB.Select("id", "title").First(&p, *refID).Error; err != nil {
			return nil
		}
		out.Label = p.Title
		out.Path = fmt.Sprintf("/prestations/%d", p.ID)
	case "event":
		var e models.Event
		if err := h.DB.Select("id", "title").First(&e, *refID).Error; err != nil {
			return nil
		}
		out.Label = e.Title
		out.Path = fmt.Sprintf("/evenements/%d", e.ID)
	case "project":
		var pr models.UpcyclingProject
		if err := h.DB.Select("id", "title").First(&pr, *refID).Error; err != nil {
			return nil
		}
		out.Label = pr.Title
		out.Path = fmt.Sprintf("/projets/%d", pr.ID)
	}
	return out
}

func (h *TicketHandler) saveImage(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Non authentifié"})
		return
	}
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Fichier image requis."})
		return
	}
	const maxSize = 5 << 20
	if file.Size > maxSize {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "L'image ne doit pas dépasser 5 Mo."})
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Format non supporté (jpg, jpeg, png, webp, gif)."})
		return
	}
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
		return
	}
	filename := fmt.Sprintf("%x%s", b, ext)
	dir := "/uploads/tickets"
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur serveur."})
		return
	}
	if err := c.SaveUploadedFile(file, filepath.Join(dir, filename)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'enregistrement."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": "/uploads/tickets/" + filename})
}

func (h *TicketHandler) UserUploadImage(c *gin.Context)  { h.saveImage(c) }
func (h *TicketHandler) StaffUploadImage(c *gin.Context) { h.saveImage(c) }

func sanitizeImages(in []string) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "/uploads/tickets/") {
			out = append(out, s)
		}
		if len(out) >= 6 {
			break
		}
	}
	return out
}

func (h *TicketHandler) UserCreate(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	var req struct {
		Subject string   `json:"subject"`
		Content string   `json:"content"`
		Images  []string `json:"images"`
		RefType *string  `json:"ref_type"`
		RefID   *uint    `json:"ref_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return
	}
	req.Subject = strings.TrimSpace(req.Subject)
	req.Content = strings.TrimSpace(req.Content)
	images := sanitizeImages(req.Images)
	if req.Subject == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Le sujet est requis."})
		return
	}
	if req.Content == "" && len(images) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Un message ou une image est requis."})
		return
	}
	if len(req.Subject) > 200 {
		req.Subject = req.Subject[:200]
	}
	var refType *string
	var refID *uint
	if req.RefType != nil && req.RefID != nil && ticketRefTypes[*req.RefType] {
		if h.resolveRef(req.RefType, req.RefID) != nil {
			refType = req.RefType
			refID = req.RefID
		}
	}

	now := time.Now()
	ticket := models.Ticket{
		UserID:        user.ID,
		Subject:       req.Subject,
		Status:        "open",
		RefType:       refType,
		RefID:         refID,
		LastMessageAt: now,
	}
	if err := h.DB.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création du ticket."})
		return
	}
	msg := models.TicketMessage{
		TicketID: ticket.ID,
		SenderID: user.ID,
		Content:  req.Content,
		Images:   images,
		IsStaff:  false,
	}
	h.DB.Create(&msg)

	if h.Notifications != nil {
		h.Notifications.NotifyAdmins("ticket.created",
			"Nouveau ticket de contact",
			fmt.Sprintf("%s : %s", strings.TrimSpace(user.FirstName+" "+user.LastName), truncateMsg(req.Subject, 80)),
			fmt.Sprintf("/admin/tickets?id=%d", ticket.ID))
	}

	h.DB.Preload("User").First(&ticket, ticket.ID)
	resp := models.ToTicketResponse(&ticket, h.resolveRef(ticket.RefType, ticket.RefID))
	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (h *TicketHandler) UserIndex(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	var tickets []models.Ticket
	h.DB.Preload("User").Where("user_id = ?", user.ID).Order("last_message_at DESC").Find(&tickets)
	out := make([]models.TicketResponse, 0, len(tickets))
	for i := range tickets {
		out = append(out, models.ToTicketResponse(&tickets[i], h.resolveRef(tickets[i].RefType, tickets[i].RefID)))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *TicketHandler) loadWithMessages(id uint) (*models.Ticket, []models.TicketMessage, bool) {
	var ticket models.Ticket
	if err := h.DB.Preload("User").First(&ticket, id).Error; err != nil {
		return nil, nil, false
	}
	var msgs []models.TicketMessage
	h.DB.Preload("Sender").Where("ticket_id = ?", ticket.ID).Order("created_at ASC").Find(&msgs)
	return &ticket, msgs, true
}

func buildTicketDetail(ticket *models.Ticket, msgs []models.TicketMessage, ref *models.TicketRefResponse) models.TicketResponse {
	resp := models.ToTicketResponse(ticket, ref)
	resp.Messages = make([]models.TicketMessageResponse, 0, len(msgs))
	for i := range msgs {
		resp.Messages = append(resp.Messages, models.ToTicketMessageResponse(&msgs[i]))
	}
	return resp
}

func (h *TicketHandler) UserShow(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket introuvable."})
		return
	}
	ticket, msgs, ok := h.loadWithMessages(uint(id))
	if !ok || ticket.UserID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket introuvable."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": buildTicketDetail(ticket, msgs, h.resolveRef(ticket.RefType, ticket.RefID))})
}

func (h *TicketHandler) UserReply(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket introuvable."})
		return
	}
	var ticket models.Ticket
	if err := h.DB.First(&ticket, id).Error; err != nil || ticket.UserID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket introuvable."})
		return
	}
	if ticket.Status == "closed" {
		c.JSON(http.StatusConflict, gin.H{"message": "Ce ticket est clôturé."})
		return
	}
	msg, ok := h.storeMessage(c, &ticket, user, false)
	if !ok {
		return
	}
	if h.Notifications != nil {
		h.Notifications.NotifyAdmins("ticket.message",
			"Réponse à un ticket",
			fmt.Sprintf("%s : %s", strings.TrimSpace(user.FirstName+" "+user.LastName), truncateMsg(msg.Content, 80)),
			fmt.Sprintf("/admin/tickets?id=%d", ticket.ID))
	}
	c.JSON(http.StatusCreated, gin.H{"data": models.ToTicketMessageResponse(msg)})
}

func (h *TicketHandler) storeMessage(c *gin.Context, ticket *models.Ticket, sender *models.User, isStaff bool) (*models.TicketMessage, bool) {
	var req struct {
		Content string   `json:"content"`
		Images  []string `json:"images"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides."})
		return nil, false
	}
	req.Content = strings.TrimSpace(req.Content)
	images := sanitizeImages(req.Images)
	if req.Content == "" && len(images) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Un message ou une image est requis."})
		return nil, false
	}
	if len(req.Content) > 4000 {
		req.Content = req.Content[:4000]
	}
	msg := models.TicketMessage{
		TicketID: ticket.ID,
		SenderID: sender.ID,
		Content:  req.Content,
		Images:   images,
		IsStaff:  isStaff,
	}
	if err := h.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de l'envoi du message."})
		return nil, false
	}
	updates := map[string]interface{}{"last_message_at": time.Now()}
	if isStaff && ticket.Status == "open" {
		updates["status"] = "in_progress"
	}
	if !isStaff && ticket.Status == "resolved" {
		updates["status"] = "in_progress"
	}
	h.DB.Model(&models.Ticket{}).Where("id = ?", ticket.ID).Updates(updates)
	msg.Sender = sender
	return &msg, true
}

func (h *TicketHandler) AdminIndex(c *gin.Context) {
	q := h.DB.Preload("User").Model(&models.Ticket{})
	if status := c.Query("status"); status != "" && ticketStatuses[status] {
		q = q.Where("status = ?", status)
	}
	if search := strings.TrimSpace(c.Query("q")); search != "" {
		q = q.Where("subject LIKE ?", "%"+search+"%")
	}
	var tickets []models.Ticket
	q.Order("last_message_at DESC").Find(&tickets)
	out := make([]models.TicketResponse, 0, len(tickets))
	for i := range tickets {
		out = append(out, models.ToTicketResponse(&tickets[i], h.resolveRef(tickets[i].RefType, tickets[i].RefID)))
	}
	var counts struct {
		Open       int64
		InProgress int64
	}
	h.DB.Model(&models.Ticket{}).Where("status = ?", "open").Count(&counts.Open)
	h.DB.Model(&models.Ticket{}).Where("status = ?", "in_progress").Count(&counts.InProgress)
	c.JSON(http.StatusOK, gin.H{"data": out, "meta": gin.H{"open": counts.Open, "in_progress": counts.InProgress}})
}

func (h *TicketHandler) AdminShow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket introuvable."})
		return
	}
	ticket, msgs, ok := h.loadWithMessages(uint(id))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket introuvable."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": buildTicketDetail(ticket, msgs, h.resolveRef(ticket.RefType, ticket.RefID))})
}

func (h *TicketHandler) AdminReply(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket introuvable."})
		return
	}
	var ticket models.Ticket
	if err := h.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket introuvable."})
		return
	}
	msg, ok := h.storeMessage(c, &ticket, user, true)
	if !ok {
		return
	}
	if h.Notifications != nil {
		h.Notifications.MustNotify(ticket.UserID, "ticket.reply",
			"Réponse à votre demande",
			fmt.Sprintf("Notre équipe a répondu à : %s", truncateMsg(ticket.Subject, 80)),
			fmt.Sprintf("/contact?id=%d", ticket.ID))
	}
	if h.Audit != nil {
		h.Audit.Log(c, "ticket.replied", "Ticket", &ticket.ID, nil, nil)
	}
	c.JSON(http.StatusCreated, gin.H{"data": models.ToTicketMessageResponse(msg)})
}

func (h *TicketHandler) AdminUpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket introuvable."})
		return
	}
	var ticket models.Ticket
	if err := h.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket introuvable."})
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || !ticketStatuses[req.Status] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Statut invalide."})
		return
	}
	h.DB.Model(&ticket).Update("status", req.Status)
	if h.Audit != nil {
		h.Audit.Log(c, "ticket.status_updated", "Ticket", &ticket.ID, nil, map[string]interface{}{"status": req.Status})
	}
	if h.Notifications != nil && (req.Status == "resolved" || req.Status == "closed") {
		label := "résolue"
		if req.Status == "closed" {
			label = "clôturée"
		}
		h.Notifications.MustNotify(ticket.UserID, "ticket.status",
			"Mise à jour de votre demande",
			fmt.Sprintf("Votre demande \"%s\" a été %s.", truncateMsg(ticket.Subject, 60), label),
			fmt.Sprintf("/contact?id=%d", ticket.ID))
	}
	c.JSON(http.StatusOK, gin.H{"message": "Statut mis à jour."})
}
