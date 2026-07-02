package models

import "time"

type Ticket struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	Subject       string    `gorm:"size:200" json:"subject"`
	Status        string    `gorm:"size:20;default:open" json:"status"`
	AssignedTo    *uint     `gorm:"index" json:"assigned_to"`
	RefType       *string   `gorm:"size:20" json:"ref_type"`
	RefID         *uint     `json:"ref_id"`
	LastMessageAt time.Time `json:"last_message_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	User     *User           `gorm:"foreignKey:UserID" json:"-"`
	Messages []TicketMessage `gorm:"foreignKey:TicketID" json:"-"`
}

func (Ticket) TableName() string {
	return "tickets"
}

type TicketMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"index;not null" json:"ticket_id"`
	SenderID  uint      `gorm:"index;not null" json:"sender_id"`
	Content   string    `gorm:"type:text" json:"content"`
	Images    []string  `gorm:"serializer:json" json:"images"`
	IsStaff   bool      `gorm:"default:false" json:"is_staff"`
	CreatedAt time.Time `json:"created_at"`

	Sender *User `gorm:"foreignKey:SenderID" json:"-"`
}

func (TicketMessage) TableName() string {
	return "ticket_messages"
}

type TicketRefResponse struct {
	Type  string `json:"type"`
	ID    uint   `json:"id"`
	Label string `json:"label"`
	Path  string `json:"path"`
}

type TicketMessageResponse struct {
	ID         uint     `json:"id"`
	SenderID   uint     `json:"sender_id"`
	SenderName string   `json:"sender_name"`
	Avatar     *string  `json:"avatar"`
	Content    string   `json:"content"`
	Images     []string `json:"images"`
	IsStaff    bool     `json:"is_staff"`
	CreatedAt  string   `json:"created_at"`
}

type TicketResponse struct {
	ID            uint                    `json:"id"`
	UserID        uint                    `json:"user_id"`
	UserName      string                  `json:"user_name"`
	UserEmail     string                  `json:"user_email"`
	Subject       string                  `json:"subject"`
	Status        string                  `json:"status"`
	Ref           *TicketRefResponse      `json:"ref"`
	LastMessageAt string                  `json:"last_message_at"`
	CreatedAt     string                  `json:"created_at"`
	Messages      []TicketMessageResponse `json:"messages,omitempty"`
}

func ToTicketMessageResponse(m *TicketMessage) TicketMessageResponse {
	name := ""
	var avatar *string
	if m.Sender != nil {
		name = m.Sender.FirstName + " " + m.Sender.LastName
		avatar = m.Sender.AvatarURL
	}
	images := m.Images
	if images == nil {
		images = []string{}
	}
	return TicketMessageResponse{
		ID:         m.ID,
		SenderID:   m.SenderID,
		SenderName: name,
		Avatar:     avatar,
		Content:    m.Content,
		Images:     images,
		IsStaff:    m.IsStaff,
		CreatedAt:  m.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func ToTicketResponse(t *Ticket, ref *TicketRefResponse) TicketResponse {
	name := ""
	email := ""
	if t.User != nil {
		name = t.User.FirstName + " " + t.User.LastName
		email = t.User.Email
	}
	return TicketResponse{
		ID:            t.ID,
		UserID:        t.UserID,
		UserName:      name,
		UserEmail:     email,
		Subject:       t.Subject,
		Status:        t.Status,
		Ref:           ref,
		LastMessageAt: t.LastMessageAt.UTC().Format(time.RFC3339Nano),
		CreatedAt:     t.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}
