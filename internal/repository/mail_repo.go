package repository

import (
	"context"

	"raven/internal/core/domain"

	"gorm.io/gorm"
)

type MailRepository struct {
	db *gorm.DB
}

func NewMailRepository(db *gorm.DB) *MailRepository {
	return &MailRepository{db: db}
}

func (r *MailRepository) Create(ctx context.Context, mail *domain.Mail) error {
	return r.db.WithContext(ctx).Create(mail).Error
}

func (r *MailRepository) GetByID(ctx context.Context, sessionID, id string) (*domain.Mail, error) {
	var mail domain.Mail
	if err := r.db.WithContext(ctx).Preload("Attachments").Preload("Recipients").
		Where("id = ? AND session_id = ?", id, sessionID).First(&mail).Error; err != nil {
		return nil, err
	}
	return &mail, nil
}

func (r *MailRepository) GetInbox(ctx context.Context, sessionID, recipientID string, page, pageSize int, queryStr string) ([]domain.Mail, int64, error) {
	var mails []domain.Mail
	var total int64

	// Join with recipients table
	query := r.db.WithContext(ctx).
		Joins("JOIN mail_recipients ON mail_recipients.mail_id = mails.id").
		Where("mail_recipients.session_id = ? AND mail_recipients.recipient_id = ? AND mail_recipients.status != 'deleted'", sessionID, recipientID)

	if queryStr != "" {
		like := "%" + queryStr + "%"
		query = query.Where("mails.subject LIKE ? OR mails.content LIKE ?", like, like)
	}

	query = query.Preload("Attachments")

	if err := query.Model(&domain.Mail{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("mails.created_at DESC").Limit(pageSize).Offset(offset).Find(&mails).Error; err != nil {
		return nil, 0, err
	}
	return mails, total, nil
}

func (r *MailRepository) GetSent(ctx context.Context, sessionID, senderID string, page, pageSize int, queryStr string) ([]domain.Mail, int64, error) {
	var mails []domain.Mail
	var total int64

	query := r.db.WithContext(ctx).Where("session_id = ? AND sender_id = ? AND (sender_status IS NULL OR sender_status != 'deleted')", sessionID, senderID)

	if queryStr != "" {
		like := "%" + queryStr + "%"
		query = query.Where("subject LIKE ? OR content LIKE ?", like, like)
	}

	query = query.Preload("Attachments")

	if err := query.Model(&domain.Mail{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&mails).Error; err != nil {
		return nil, 0, err
	}
	return mails, total, nil
}

func (r *MailRepository) UpdateStatus(ctx context.Context, mailID, recipientID, status string) error {
	return r.db.WithContext(ctx).Model(&domain.MailRecipient{}).
		Where("mail_id = ? AND recipient_id = ?", mailID, recipientID).
		Update("status", status).Error
}

func (r *MailRepository) DeleteForSender(ctx context.Context, mailID string) error {
	return r.db.WithContext(ctx).Model(&domain.Mail{}).
		Where("id = ?", mailID).
		Update("sender_status", "deleted").Error
}

func (r *MailRepository) DeleteSession(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Hard delete from all tables
		if err := tx.Unscoped().Where("session_id = ?", sessionID).Delete(&domain.Attachment{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("session_id = ?", sessionID).Delete(&domain.MailRecipient{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("session_id = ?", sessionID).Delete(&domain.Mail{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *MailRepository) GetAttachmentByID(ctx context.Context, sessionID, id string) (*domain.Attachment, error) {
	var att domain.Attachment
	if err := r.db.WithContext(ctx).Where("id = ? AND session_id = ?", id, sessionID).First(&att).Error; err != nil {
		return nil, err
	}
	return &att, nil
}
