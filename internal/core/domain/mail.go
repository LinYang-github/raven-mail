package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Mail 代表文电实体
type Mail struct {
	ID           string         `gorm:"primaryKey;type:uuid" json:"id"`
	SessionID    string         `gorm:"index;not null;default:'default'" json:"session_id"`
	SenderID     string         `gorm:"index;not null" json:"sender_id"`
	SenderStatus string         `gorm:"type:varchar(20);default:'normal'" json:"-"` // normal, deleted
	Subject      string         `gorm:"not null" json:"subject"`
	Content      string         `gorm:"type:text" json:"content"`
	ContentType  string         `gorm:"type:varchar(32);default:'text'" json:"content_type"`
	ParentID     *string        `gorm:"index" json:"parent_id,omitempty"` // 用于会话/回复
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Attachments []Attachment    `gorm:"foreignKey:MailID" json:"attachments"`
	Recipients  []MailRecipient `gorm:"foreignKey:MailID" json:"recipients"`
}

// MailRecipient 代表文电与接收者之间的关系
type MailRecipient struct {
	ID          string     `gorm:"primaryKey;type:uuid" json:"id"`
	MailID      string     `gorm:"index;not null" json:"mail_id"`
	SessionID   string     `gorm:"index;not null;default:'default'" json:"session_id"` // 为了方便查询收件箱而冗余的字段
	RecipientID string     `gorm:"index;not null" json:"recipient_id"`                 // 用户 ID 或组 ID
	Type        string     `gorm:"type:varchar(10);default:'to'" json:"type"`          // to (收件人), cc (抄送), bcc (密送)
	Status      string     `gorm:"type:varchar(20);default:'unread'" json:"status"`    // unread, read, deleted
	ReadAt      *time.Time `json:"read_at,omitempty"`
}

// Attachment 代表附件文件
type Attachment struct {
	ID            string    `gorm:"primaryKey;type:uuid" json:"id"`
	MailID        *string   `gorm:"index" json:"mail_id,omitempty"`
	ChatMessageID *string   `gorm:"index" json:"chat_message_id,omitempty"`
	SessionID     string    `gorm:"index;not null;default:'default'" json:"session_id"`
	FileName      string    `gorm:"not null" json:"file_name"`
	FilePath      string    `gorm:"not null" json:"file_path"` // 对象存储路径或本地路径
	FileSize      int64     `json:"file_size"`
	MimeType      string    `json:"mime_type"`
	CreatedAt     time.Time `json:"created_at"`
}

type ChatMessage struct {
	ID          string       `gorm:"primaryKey" json:"id"`
	SessionID   string       `gorm:"index" json:"session_id"`
	SenderID    string       `gorm:"index" json:"sender_id"`
	ReceiverID  string       `gorm:"index" json:"receiver_id"` // 接收用户 ID
	Content     string       `json:"content"`
	IsRead      bool         `gorm:"default:false" json:"is_read"`
	CreatedAt   time.Time    `json:"created_at"`
	Attachments []Attachment `gorm:"foreignKey:ChatMessageID" json:"attachments"`
}

// BeforeCreate 钩子：生成 UUID
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

func (cm *ChatMessage) BeforeCreate(tx *gorm.DB) (err error) {
	if cm.ID == "" {
		cm.ID = uuid.New().String()
	}
	return
}
