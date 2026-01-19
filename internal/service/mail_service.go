package service

import (
	"context"
	"raven/internal/core/domain"
	"raven/internal/core/ports"
	"time"
)

type MailService struct {
	repo    ports.MailRepository
	storage ports.StorageService
}

func NewMailService(repo ports.MailRepository, storage ports.StorageService) *MailService {
	return &MailService{
		repo:    repo,
		storage: storage,
	}
}

func (s *MailService) SendMail(ctx context.Context, senderID string, req ports.SendMailRequest) (*domain.Mail, error) {
	// Handle Attachments
	var attachments []domain.Attachment
	for _, attReq := range req.Attachments {
		path, err := s.storage.UploadFile(ctx, attReq.FileName, attReq.Content)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, domain.Attachment{
			FileName: attReq.FileName,
			FilePath: path,
			FileSize: attReq.Size,
			MimeType: attReq.MimeType,
		})
	}

	mail := &domain.Mail{
		SenderID:    senderID,
		Subject:     req.Subject,
		Content:     req.Content,
		Attachments: attachments,
		CreatedAt:   time.Now(),
	}

	// Recipients
	var recipients []domain.MailRecipient

	addRecipients := func(ids []string, rType string) {
		for _, id := range ids {
			recipients = append(recipients, domain.MailRecipient{
				RecipientID: id,
				Type:        rType,
				Status:      "unread",
			})
		}
	}

	addRecipients(req.To, "to")
	addRecipients(req.Cc, "cc")
	addRecipients(req.Bcc, "bcc")

	mail.Recipients = recipients

	if err := s.repo.Create(ctx, mail); err != nil {
		return nil, err
	}

	return mail, nil
}

func (s *MailService) GetInbox(ctx context.Context, userID string, page, pageSize int) ([]domain.Mail, int64, error) {
	return s.repo.GetInbox(ctx, userID, page, pageSize)
}

func (s *MailService) GetSent(ctx context.Context, userID string, page, pageSize int) ([]domain.Mail, int64, error) {
	return s.repo.GetSent(ctx, userID, page, pageSize)
}

func (s *MailService) ReadMail(ctx context.Context, userID, mailID string) (*domain.Mail, error) {
	mail, err := s.repo.GetByID(ctx, mailID)
	if err != nil {
		return nil, err
	}

	// Update status to read if it's the recipient
	// Note: Logic could be optimized to check if already read
	_ = s.repo.UpdateStatus(ctx, mailID, userID, "read")

	return mail, nil
}
