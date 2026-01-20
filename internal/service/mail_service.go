package service

import (
	"context"
	"raven/internal/core/domain"
	"raven/internal/core/ports"
	"strings"
	"time"
)

type MailService struct {
	repo    ports.MailRepository
	storage ports.StorageService
	// Simple Notification Hub
	clients map[chan string]bool
	msgChan chan string
}

func NewMailService(repo ports.MailRepository, storage ports.StorageService) *MailService {
	s := &MailService{
		repo:    repo,
		storage: storage,
		clients: make(map[chan string]bool),
		msgChan: make(chan string),
	}
	go s.runHub()
	return s
}

func (s *MailService) runHub() {
	for msg := range s.msgChan {
		for client := range s.clients {
			client <- msg
		}
	}
}

func (s *MailService) Subscribe() chan string {
	c := make(chan string)
	s.clients[c] = true
	return c
}

func (s *MailService) Unsubscribe(c chan string) {
	delete(s.clients, c)
	close(c)
}

func (s *MailService) SendMail(ctx context.Context, senderID string, req ports.SendMailRequest) (*domain.Mail, error) {
	// Handle Attachments
	var attachments []domain.Attachment
	for _, attReq := range req.Attachments {
		path, err := s.storage.UploadFile(ctx, req.SessionID, attReq.FileName, attReq.Content)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, domain.Attachment{
			SessionID: req.SessionID,
			FileName:  attReq.FileName,
			FilePath:  path,
			FileSize:  attReq.Size,
			MimeType:  attReq.MimeType,
		})
	}

	mail := &domain.Mail{
		SessionID:   req.SessionID,
		SenderID:    senderID,
		Subject:     req.Subject,
		Content:     req.Content,
		ContentType: req.ContentType,
		Attachments: attachments,
		CreatedAt:   time.Now(),
	}

	// Recipients
	var recipients []domain.MailRecipient

	addRecipients := func(ids []string, rType string) {
		for _, id := range ids {
			recipients = append(recipients, domain.MailRecipient{
				SessionID:   req.SessionID,
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

	// Broadcast notification: target_user_ids (Maybe include session in notification too?)
	var targetIDs []string
	for _, r := range mail.Recipients {
		targetIDs = append(targetIDs, r.RecipientID)
	}
	// Format: NEW_MAIL:session_id:comma_separated_ids
	s.msgChan <- "NEW_MAIL:" + req.SessionID + ":" + strings.Join(targetIDs, ",")

	return mail, nil
}

func (s *MailService) GetInbox(ctx context.Context, sessionID, userID string, page, pageSize int, query string) ([]domain.Mail, int64, error) {
	return s.repo.GetInbox(ctx, sessionID, userID, page, pageSize, query)
}

func (s *MailService) GetSent(ctx context.Context, sessionID, userID string, page, pageSize int, query string) ([]domain.Mail, int64, error) {
	return s.repo.GetSent(ctx, sessionID, userID, page, pageSize, query)
}

func (s *MailService) ReadMail(ctx context.Context, sessionID, userID, mailID string) (*domain.Mail, error) {
	mail, err := s.repo.GetByID(ctx, sessionID, mailID)
	if err != nil {
		return nil, err
	}

	// Update status to read if it's the recipient
	// Note: Logic could be optimized to check if already read
	_ = s.repo.UpdateStatus(ctx, mailID, userID, "read")

	return mail, nil
}

func (s *MailService) DeleteMail(ctx context.Context, sessionID, userID, mailID string) error {
	mail, err := s.repo.GetByID(ctx, sessionID, mailID)
	if err != nil {
		return err
	}

	if mail.SenderID == userID {
		// User is sender -> delete for sender
		return s.repo.DeleteForSender(ctx, mailID)
	}

	// User is recipient -> delete for recipient
	return s.repo.UpdateStatus(ctx, mailID, userID, "deleted")
}

func (s *MailService) GetAttachment(ctx context.Context, sessionID, attachmentID string) (*domain.Attachment, error) {
	return s.repo.GetAttachmentByID(ctx, sessionID, attachmentID)
}
