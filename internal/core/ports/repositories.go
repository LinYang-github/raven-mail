package ports

import (
	"context"
	"raven/internal/core/domain"
)

type MailRepository interface {
	Create(ctx context.Context, mail *domain.Mail) error
	GetByID(ctx context.Context, sessionID, id string) (*domain.Mail, error)
	GetInbox(ctx context.Context, sessionID, recipientID string, page, pageSize int, query string) ([]domain.Mail, int64, error)
	GetSent(ctx context.Context, sessionID, senderID string, page, pageSize int, query string) ([]domain.Mail, int64, error)
	UpdateStatus(ctx context.Context, mailID, recipientID, status string) error
	DeleteForSender(ctx context.Context, mailID string) error
	DeleteSession(ctx context.Context, sessionID string) error
	GetAttachmentByID(ctx context.Context, sessionID, id string) (*domain.Attachment, error)
}
