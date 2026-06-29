package models

import "time"

type StaffShift struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID uint      `gorm:"index;not null" json:"employee_id"`
	Weekday    int       `gorm:"not null" json:"weekday"`
	StartTime  string    `gorm:"size:5;not null" json:"start_time"`
	EndTime    string    `gorm:"size:5;not null" json:"end_time"`
	Label      *string   `gorm:"size:120" json:"label"`
	CreatedAt  time.Time `json:"created_at"`

	Employee *User `gorm:"foreignKey:EmployeeID" json:"-"`
}

func (StaffShift) TableName() string { return "staff_shifts" }

type StaffEvent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Description *string   `gorm:"type:text" json:"description"`
	Location    *string   `gorm:"size:200" json:"location"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	CreatedBy   *uint     `gorm:"index" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Members []User `gorm:"many2many:staff_event_members;" json:"-"`
}

func (StaffEvent) TableName() string { return "staff_events" }

type StaffShiftResponse struct {
	ID           uint    `json:"id"`
	EmployeeID   uint    `json:"employee_id"`
	EmployeeName string  `json:"employee_name"`
	Weekday      int     `json:"weekday"`
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	Label        *string `json:"label"`
}

func ToStaffShiftResponse(s *StaffShift) StaffShiftResponse {
	name := ""
	if s.Employee != nil {
		name = s.Employee.FirstName + " " + s.Employee.LastName
	}
	return StaffShiftResponse{
		ID: s.ID, EmployeeID: s.EmployeeID, EmployeeName: name,
		Weekday: s.Weekday, StartTime: s.StartTime, EndTime: s.EndTime, Label: s.Label,
	}
}

type StaffEventMemberResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type StaffEventResponse struct {
	ID          uint                       `json:"id"`
	Title       string                     `json:"title"`
	Description *string                    `json:"description"`
	Location    *string                    `json:"location"`
	StartAt     string                     `json:"start_at"`
	EndAt       string                     `json:"end_at"`
	Members     []StaffEventMemberResponse `json:"members"`
}

func ToStaffEventResponse(e *StaffEvent) StaffEventResponse {
	members := make([]StaffEventMemberResponse, 0, len(e.Members))
	for i := range e.Members {
		m := &e.Members[i]
		members = append(members, StaffEventMemberResponse{ID: m.ID, Name: m.FirstName + " " + m.LastName})
	}
	return StaffEventResponse{
		ID: e.ID, Title: e.Title, Description: e.Description, Location: e.Location,
		StartAt: e.StartAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
		EndAt:   e.EndAt.UTC().Format("2006-01-02T15:04:05.000000Z"),
		Members: members,
	}
}
