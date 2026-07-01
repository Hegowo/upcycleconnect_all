package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogsHandler struct {
	DB *gorm.DB
}

type auditLogRow struct {
	ID           uint      `json:"id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   *uint     `json:"resource_id"`
	AdminID      *uint     `json:"admin_id"`
	AdminName    string    `json:"admin_name"`
	AdminRole    string    `json:"admin_role"`
	IPAddress    *string   `json:"ip_address"`
	OldValues    *string   `json:"old_values"`
	NewValues    *string   `json:"new_values"`
	CreatedAt    time.Time `json:"created_at"`
}

func topRole(roles []string) string {
	for _, want := range []string{"super_admin", "admin", "employee"} {
		for _, r := range roles {
			if r == want {
				return want
			}
		}
	}
	if len(roles) > 0 {
		return roles[0]
	}
	return ""
}

func (h *LogsHandler) rolesByUser(ids []uint) map[uint]string {
	out := map[uint]string{}
	if len(ids) == 0 {
		return out
	}
	var rows []struct {
		UserID uint
		Name   string
	}
	h.DB.Table("user_roles ur").
		Select("ur.user_id, r.name").
		Joins("JOIN roles r ON r.id = ur.role_id").
		Where("ur.user_id IN ?", ids).
		Scan(&rows)
	grouped := map[uint][]string{}
	for _, r := range rows {
		grouped[r.UserID] = append(grouped[r.UserID], r.Name)
	}
	for uid, rs := range grouped {
		out[uid] = topRole(rs)
	}
	return out
}

func (h *LogsHandler) Index(c *gin.Context) {
	base := h.DB.Table("audit_logs al")

	if action := c.Query("action"); action != "" {
		base = base.Where("al.action = ?", action)
	}
	if rt := c.Query("resource_type"); rt != "" {
		base = base.Where("al.resource_type = ?", rt)
	}
	if actor := c.Query("actor"); actor != "" {
		if actor == "system" {
			base = base.Where("al.user_id IS NULL")
		} else {
			base = base.Where("al.user_id = ?", actor)
		}
	}
	if search := strings.TrimSpace(c.Query("q")); search != "" {
		like := "%" + search + "%"
		base = base.Where("al.action LIKE ? OR al.resource_type LIKE ?", like, like)
	}
	if from := c.Query("from"); from != "" {
		base = base.Where("al.created_at >= ?", from+" 00:00:00")
	}
	if to := c.Query("to"); to != "" {
		base = base.Where("al.created_at <= ?", to+" 23:59:59")
	}

	var total int64
	base.Count(&total)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	const limit = 30
	offset := (page - 1) * limit

	var rows []struct {
		models.AuditLog
		AdminName string
	}
	base.Select("al.*, CONCAT(COALESCE(u.first_name,''), ' ', COALESCE(u.last_name,'')) as admin_name").
		Joins("LEFT JOIN users u ON u.id = al.user_id").
		Order("al.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&rows)

	var actorIDs []uint
	for _, r := range rows {
		if r.UserID != nil {
			actorIDs = append(actorIDs, *r.UserID)
		}
	}
	roles := h.rolesByUser(actorIDs)

	result := make([]auditLogRow, len(rows))
	for i, r := range rows {
		role := ""
		if r.UserID != nil {
			role = roles[*r.UserID]
		}
		result[i] = auditLogRow{
			ID:           r.ID,
			Action:       r.Action,
			ResourceType: r.ResourceType,
			ResourceID:   r.ResourceID,
			AdminID:      r.UserID,
			AdminName:    r.AdminName,
			AdminRole:    role,
			IPAddress:    r.IPAddress,
			OldValues:    r.OldValues,
			NewValues:    r.NewValues,
			CreatedAt:    r.CreatedAt,
		}
	}

	var actions []string
	h.DB.Table("audit_logs").Distinct().Where("action <> ''").Order("action").Pluck("action", &actions)
	var resourceTypes []string
	h.DB.Table("audit_logs").Distinct().Where("resource_type <> ''").Order("resource_type").Pluck("resource_type", &resourceTypes)

	var actorRows []struct {
		ID   uint
		Name string
	}
	h.DB.Table("audit_logs al").
		Select("DISTINCT u.id, CONCAT(COALESCE(u.first_name,''), ' ', COALESCE(u.last_name,'')) as name").
		Joins("JOIN users u ON u.id = al.user_id").
		Scan(&actorRows)
	var actorIDList []uint
	for _, a := range actorRows {
		actorIDList = append(actorIDList, a.ID)
	}
	actorRoles := h.rolesByUser(actorIDList)
	type actorOption struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	}
	actors := make([]actorOption, 0, len(actorRows))
	for _, a := range actorRows {
		actors = append(actors, actorOption{ID: a.ID, Name: strings.TrimSpace(a.Name), Role: actorRoles[a.ID]})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
		"meta": gin.H{"total": total, "page": page, "limit": limit},
		"filters": gin.H{
			"actions":        actions,
			"resource_types": resourceTypes,
			"actors":         actors,
		},
	})
}

type activityEntry struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Company   string    `json:"company,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *LogsHandler) Activity(c *gin.Context) {
	type userRow struct {
		ID        uint
		FirstName string
		LastName  string
		Email     string
		CreatedAt time.Time
	}
	var users []userRow
	h.DB.Raw(`
		SELECT id, first_name, last_name, email, created_at
		FROM users
		WHERE deleted_at IS NULL
		  AND id NOT IN (SELECT user_id FROM user_roles)
		ORDER BY created_at DESC
		LIMIT 50
	`).Scan(&users)

	type providerRow struct {
		ID          uint
		FirstName   string
		LastName    string
		Email       string
		CompanyName string
		Status      string
		CreatedAt   time.Time
	}
	var providers []providerRow
	h.DB.Raw(`
		SELECT u.id, u.first_name, u.last_name, u.email, pp.company_name, pp.status, u.created_at
		FROM users u
		JOIN provider_profiles pp ON pp.user_id = u.id
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at DESC
		LIMIT 50
	`).Scan(&providers)

	var entries []activityEntry
	for _, u := range users {
		entries = append(entries, activityEntry{
			ID:        u.ID,
			Type:      "user",
			Name:      u.FirstName + " " + u.LastName,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		})
	}
	for _, p := range providers {
		entries = append(entries, activityEntry{
			ID:        p.ID,
			Type:      "provider",
			Name:      p.FirstName + " " + p.LastName,
			Email:     p.Email,
			Company:   p.CompanyName,
			Status:    p.Status,
			CreatedAt: p.CreatedAt,
		})
	}

	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].CreatedAt.After(entries[i].CreatedAt) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
	if len(entries) > 50 {
		entries = entries[:50]
	}

	c.JSON(http.StatusOK, gin.H{"data": entries})
}
