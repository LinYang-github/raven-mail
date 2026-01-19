package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Mail represents an email entity
type Mail struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	SenderID  string         `gorm:"index;not null" json:"sender_id"`
	Subject   string         `gorm:"not null" json:"subject"`
	Content   string         `gorm:"type:text" json:"content"`
	ParentID  *string        `gorm:"index" json:"parent_id,omitempty"` // For threads/replies
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Attachments []Attachment    `gorm:"foreignKey:MailID" json:"attachments"`
	Recipients  []MailRecipient `gorm:"foreignKey:MailID" json:"recipients"`
}

// MailRecipient represents the relationship between a mail and a recipient
type MailRecipient struct {
	ID          string     `gorm:"primaryKey;type:uuid" json:"id"`
	MailID      string     `gorm:"index;not null" json:"mail_id"`
	RecipientID string     `gorm:"index;not null" json:"recipient_id"`              // UserID or GroupID
	Type        string     `gorm:"type:varchar(10);default:'to'" json:"type"`       // to, cc, bcc
	Status      string     `gorm:"type:varchar(20);default:'unread'" json:"status"` // unread, read, deleted
	ReadAt      *time.Time `json:"read_at,omitempty"`
}

// Attachment represents a file attachment
type Attachment struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	MailID    string    `gorm:"index;not null" json:"mail_id"`
	FileName  string    `gorm:"not null" json:"file_name"`
	FilePath  string    `gorm:"not null" json:"file_path"` // Path in object storage or simple FS
	FileSize  int64     `json:"file_size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate hooks to set UUIDs
func (m *Mail) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

func (mr *MailRecipient) BeforeCreate(tx *gorm.DB) (err error) {
	if mr.ID == "" {
		mr.ID = uuid.New().String()
	}
	return
}

func (a *Attachment) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return
}
