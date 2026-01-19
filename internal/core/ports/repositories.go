package ports

import (
	"context"
	"raven/internal/core/domain"
)

type MailRepository interface {
	Create(ctx context.Context, mail *domain.Mail) error
	GetByID(ctx context.Context, id string) (*domain.Mail, error)
	GetInbox(ctx context.Context, recipientID string, page, pageSize int, query string) ([]domain.Mail, int64, error)
	GetSent(ctx context.Context, senderID string, page, pageSize int, query string) ([]domain.Mail, int64, error)
	UpdateStatus(ctx context.Context, mailID, recipientID, status string) error
	DeleteForSender(ctx context.Context, mailID string) error
	GetAttachmentByID(ctx context.Context, id string) (*domain.Attachment, error)
}
