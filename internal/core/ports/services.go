package ports

import (
	"context"
	"io"
	"raven/internal/core/domain"
)

type MailService interface {
	SendMail(ctx context.Context, senderID string, req SendMailRequest) (*domain.Mail, error)
	GetInbox(ctx context.Context, sessionID, userID string, page, pageSize int, query string) ([]domain.Mail, int64, error)
	GetSent(ctx context.Context, sessionID, userID string, page, pageSize int, query string) ([]domain.Mail, int64, error)
	ReadMail(ctx context.Context, sessionID, userID, mailID string) (*domain.Mail, error)
	DeleteMail(ctx context.Context, sessionID, userID, mailID string) error
	DeleteSession(ctx context.Context, sessionID string) error
	GetAttachment(ctx context.Context, sessionID, attachmentID string) (*domain.Attachment, error)
	// Notification stream
	Subscribe() chan string
	Unsubscribe(chan string)
}

type StorageService interface {
	UploadFile(ctx context.Context, sessionID, fileName string, content io.Reader) (string, error)
	GetFile(ctx context.Context, path string) (io.ReadCloser, error)
	DeleteSessionDir(ctx context.Context, sessionID string) error
}

type SendMailRequest struct {
	SessionID   string
	Subject     string
	Content     string
	ContentType string
	To          []string // UserIDs
	Cc          []string
	Bcc         []string
	Attachments []AttachmentRequest
}

type AttachmentRequest struct {
	FileName string
	Content  io.Reader
	Size     int64
	MimeType string
}
