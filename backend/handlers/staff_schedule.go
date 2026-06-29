package handlers

import (
	"net/http"
	"strconv"

	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StaffScheduleHandler struct {
	DB            *gorm.DB
	Audit         *services.AuditService
	Notifications *services.NotificationService
}

func (h *StaffScheduleHandler) targetEmployeeID(c *gin.Context) (uint, bool) {
	user := middleware.GetAuthUser(c)
	if user == nil {
		return 0, false
	}
	if !user.IsAdmin() {
		return user.ID, true
	}
	if q := c.Query("employee_id"); q != "" {
		id, err := strconv.ParseUint(q, 10, 64)
		if err != nil {
			return 0, false
		}
		return uint(id), true
	}
	return 0, true
}

func (h *StaffScheduleHandler) ListShifts(c *gin.Context) {
	empID, ok := h.targetEmployeeID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Paramètre invalide"})
		return
	}
	q := h.DB.Preload("Employee").Order("weekday ASC, start_time ASC")
	if empID != 0 {
		q = q.Where("employee_id = ?", empID)
	}
	var shifts []models.StaffShift
	q.Find(&shifts)
	out := make([]models.StaffShiftResponse, 0, len(shifts))
	for i := range shifts {
		out = append(out, models.ToStaffShiftResponse(&shifts[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *StaffScheduleHandler) CreateShift(c *gin.Context) {
	var req struct {
		EmployeeID uint    `json:"employee_id" binding:"required"`
		Weekday    int     `json:"weekday"`
		StartTime  string  `json:"start_time" binding:"required"`
		EndTime    string  `json:"end_time" binding:"required"`
		Label      *string `json:"label"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides (employé, jour, horaires requis)."})
		return
	}
	if req.Weekday < 1 || req.Weekday > 7 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Jour invalide (1=lundi … 7=dimanche)."})
		return
	}
	if req.EndTime <= req.StartTime {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "L'heure de fin doit être après l'heure de début."})
		return
	}

	var emp models.User
	if err := h.DB.Preload("Roles").First(&emp, req.EmployeeID).Error; err != nil || !emp.IsEmployee() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Employé introuvable."})
		return
	}

	shift := models.StaffShift{
		EmployeeID: req.EmployeeID, Weekday: req.Weekday,
		StartTime: req.StartTime, EndTime: req.EndTime, Label: req.Label,
	}
	if err := h.DB.Create(&shift).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création."})
		return
	}
	h.Audit.Log(c, "staff_shift.created", "StaffShift", &shift.ID, nil, map[string]interface{}{"employee_id": req.EmployeeID})
	h.DB.Preload("Employee").First(&shift, shift.ID)
	c.JSON(http.StatusCreated, models.ToStaffShiftResponse(&shift))
}

func (h *StaffScheduleHandler) DeleteShift(c *gin.Context) {
	h.DB.Where("id = ?", c.Param("id")).Delete(&models.StaffShift{})
	c.Status(http.StatusNoContent)
}

func (h *StaffScheduleHandler) ListEvents(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	q := h.DB.Preload("Members").Order("start_at ASC")
	if from := c.Query("from"); from != "" {
		q = q.Where("end_at >= ?", from)
	}
	if to := c.Query("to"); to != "" {
		q = q.Where("start_at <= ?", to)
	}

	var events []models.StaffEvent
	q.Find(&events)

	out := make([]models.StaffEventResponse, 0, len(events))
	for i := range events {
		e := &events[i]
		if user != nil && !user.IsAdmin() {
			member := false
			for j := range e.Members {
				if e.Members[j].ID == user.ID {
					member = true
					break
				}
			}
			if !member {
				continue
			}
		} else if empID := c.Query("employee_id"); empID != "" && user != nil && user.IsAdmin() {
			id, _ := strconv.ParseUint(empID, 10, 64)
			member := false
			for j := range e.Members {
				if uint64(e.Members[j].ID) == id {
					member = true
					break
				}
			}
			if !member {
				continue
			}
		}
		out = append(out, models.ToStaffEventResponse(e))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *StaffScheduleHandler) parseEventReq(c *gin.Context) (*models.StaffEvent, []uint, bool) {
	var req struct {
		Title       string  `json:"title" binding:"required"`
		Description *string `json:"description"`
		Location    *string `json:"location"`
		StartAt     string  `json:"start_at" binding:"required"`
		EndAt       string  `json:"end_at" binding:"required"`
		EmployeeIDs []uint  `json:"employee_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Données invalides (titre, dates requis)."})
		return nil, nil, false
	}
	start, err := parseDateTime(req.StartAt)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Date de début invalide."})
		return nil, nil, false
	}
	end, err := parseDateTime(req.EndAt)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Date de fin invalide."})
		return nil, nil, false
	}
	if !end.After(start) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "La fin doit être après le début."})
		return nil, nil, false
	}
	return &models.StaffEvent{
		Title: req.Title, Description: req.Description, Location: req.Location,
		StartAt: start, EndAt: end,
	}, req.EmployeeIDs, true
}

func (h *StaffScheduleHandler) setMembers(event *models.StaffEvent, ids []uint) {
	var members []models.User
	if len(ids) > 0 {
		h.DB.Where("id IN ?", ids).Find(&members)
	}
	h.DB.Model(event).Association("Members").Replace(members)
}

func (h *StaffScheduleHandler) CreateEvent(c *gin.Context) {
	event, ids, ok := h.parseEventReq(c)
	if !ok {
		return
	}
	if u := middleware.GetAuthUser(c); u != nil {
		event.CreatedBy = &u.ID
	}
	if err := h.DB.Create(event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la création."})
		return
	}
	h.setMembers(event, ids)
	h.notifyMembers(event, ids)
	h.Audit.Log(c, "staff_event.created", "StaffEvent", &event.ID, nil, map[string]interface{}{"title": event.Title})
	h.DB.Preload("Members").First(event, event.ID)
	c.JSON(http.StatusCreated, models.ToStaffEventResponse(event))
}

func (h *StaffScheduleHandler) UpdateEvent(c *gin.Context) {
	var existing models.StaffEvent
	if err := h.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Événement introuvable."})
		return
	}
	event, ids, ok := h.parseEventReq(c)
	if !ok {
		return
	}
	h.DB.Model(&existing).Updates(map[string]interface{}{
		"title": event.Title, "description": event.Description, "location": event.Location,
		"start_at": event.StartAt, "end_at": event.EndAt,
	})
	h.setMembers(&existing, ids)
	h.Audit.Log(c, "staff_event.updated", "StaffEvent", &existing.ID, nil, nil)
	h.DB.Preload("Members").First(&existing, existing.ID)
	c.JSON(http.StatusOK, models.ToStaffEventResponse(&existing))
}

func (h *StaffScheduleHandler) DeleteEvent(c *gin.Context) {
	var event models.StaffEvent
	if err := h.DB.First(&event, c.Param("id")).Error; err != nil {
		c.Status(http.StatusNoContent)
		return
	}
	h.DB.Model(&event).Association("Members").Clear()
	h.DB.Delete(&event)
	c.Status(http.StatusNoContent)
}

func (h *StaffScheduleHandler) notifyMembers(event *models.StaffEvent, ids []uint) {
	if h.Notifications == nil {
		return
	}
	when := event.StartAt.Format("02/01/2006 15:04")
	for _, id := range ids {
		h.Notifications.MustNotify(id, "staff_event.assigned",
			"Nouvel événement à votre planning",
			"Vous êtes assigné à « "+event.Title+" » le "+when+".",
			"/admin/planning")
	}
}
