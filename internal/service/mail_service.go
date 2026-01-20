package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"raven/internal/core/domain"
	"raven/internal/core/ports"
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

	// Broadcast notification
	var targetIDs []string
	for _, r := range mail.Recipients {
		targetIDs = append(targetIDs, r.RecipientID)
	}

	payload := map[string]interface{}{
		"type":       "MAIL",
		"session_id": req.SessionID,
		"targets":    targetIDs,
		"data": map[string]interface{}{
			"id":        mail.ID,
			"subject":   mail.Subject,
			"sender_id": mail.SenderID,
		},
	}
	s.broadcast(payload)

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

	// 如果当前查看者是收件人之一，且状态还是 unread，则更新为 read
	for _, r := range mail.Recipients {
		if r.RecipientID == userID && r.Status == "unread" {
			_ = s.repo.UpdateStatus(ctx, mailID, userID, "read")
			break
		}
	}

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

func (s *MailService) DeleteSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}

	// 1. Delete DB records
	if err := s.repo.DeleteSession(ctx, sessionID); err != nil {
		return err
	}

	// 2. Delete Uploads
	_ = s.storage.DeleteSessionDir(ctx, sessionID)

	// 3. Delete OnlyOffice docs (hardcoded path in handler, but we can handle it here too)
	_ = os.RemoveAll(filepath.Join("./data", sessionID))

	return nil
}

func (s *MailService) GetAttachment(ctx context.Context, sessionID, attachmentID string) (*domain.Attachment, error) {
	return s.repo.GetAttachmentByID(ctx, sessionID, attachmentID)
}

func (s *MailService) SendChatMessage(ctx context.Context, senderID string, req ports.SendChatMessageRequest) (*domain.ChatMessage, error) {
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

	msg := &domain.ChatMessage{
		SessionID:   req.SessionID,
		SenderID:    senderID,
		ReceiverID:  req.ReceiverID,
		Content:     req.Content,
		Attachments: attachments,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateChatMessage(ctx, msg); err != nil {
		return nil, err
	}

	// Broadcast
	payload := map[string]interface{}{
		"type":       "CHAT",
		"session_id": req.SessionID,
		"targets":    []string{req.ReceiverID},
		"data":       msg,
	}
	s.broadcast(payload)

	return msg, nil
}

func (s *MailService) GetChatHistory(ctx context.Context, sessionID, userA, userB string) ([]domain.ChatMessage, error) {
	return s.repo.GetChatHistory(ctx, sessionID, userA, userB, 100)
}

func (s *MailService) MarkChatAsRead(ctx context.Context, sessionID, senderID, receiverID string) error {
	return s.repo.MarkChatAsRead(ctx, sessionID, senderID, receiverID)
}

func (s *MailService) GetUserSummary(ctx context.Context, sessionID, userID string) (*ports.UserSummary, error) {
	mailCount, err := s.repo.GetUnreadMailCount(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}

	imCounts, err := s.repo.GetIMUnreadCounts(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}

	return &ports.UserSummary{
		UnreadMailCount: mailCount,
		IMUnreadCounts:  imCounts,
	}, nil
}

func (s *MailService) broadcast(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	s.msgChan <- string(data)
}
