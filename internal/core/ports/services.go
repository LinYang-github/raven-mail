package ports

import (
	"context"
	"io"
	"raven/internal/core/domain"
)

type MailService interface {
	SendMail(ctx context.Context, senderID string, req SendMailRequest) (*domain.Mail, error)
	GetInbox(ctx context.Context, userID string, page, pageSize int, query string) ([]domain.Mail, int64, error)
	GetSent(ctx context.Context, userID string, page, pageSize int, query string) ([]domain.Mail, int64, error)
	ReadMail(ctx context.Context, userID, mailID string) (*domain.Mail, error)
	DeleteMail(ctx context.Context, userID, mailID string) error
	GetAttachment(ctx context.Context, attachmentID string) (*domain.Attachment, error)
	// Notification stream
	Subscribe() chan string
	Unsubscribe(chan string)
}

type StorageService interface {
	UploadFile(ctx context.Context, fileName string, content io.Reader) (string, error)
	GetFile(ctx context.Context, path string) (io.ReadCloser, error)
}

type SendMailRequest struct {
	Subject     string
	Content     string
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
