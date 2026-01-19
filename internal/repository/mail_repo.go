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

func (r *MailRepository) GetByID(ctx context.Context, id string) (*domain.Mail, error) {
	var mail domain.Mail
	if err := r.db.WithContext(ctx).Preload("Attachments").Preload("Recipients").First(&mail, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &mail, nil
}

func (r *MailRepository) GetInbox(ctx context.Context, recipientID string, page, pageSize int) ([]domain.Mail, int64, error) {
	var mails []domain.Mail
	var total int64

	// Join with recipients table
	query := r.db.WithContext(ctx).
		Joins("JOIN mail_recipients ON mail_recipients.mail_id = mails.id").
		Where("mail_recipients.recipient_id = ? AND mail_recipients.status != 'deleted'", recipientID).
		Preload("Attachments")

	if err := query.Model(&domain.Mail{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("mails.created_at DESC").Limit(pageSize).Offset(offset).Find(&mails).Error; err != nil {
		return nil, 0, err
	}
	return mails, total, nil
}

func (r *MailRepository) GetSent(ctx context.Context, senderID string, page, pageSize int) ([]domain.Mail, int64, error) {
	var mails []domain.Mail
	var total int64

	query := r.db.WithContext(ctx).Where("sender_id = ?", senderID).Preload("Attachments")

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
