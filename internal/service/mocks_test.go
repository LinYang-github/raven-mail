package service

import (
	"context"
	"io"
	"raven/internal/core/domain"
	"raven/internal/core/ports"

	"github.com/stretchr/testify/mock"
)

// MockMailRepository
type MockMailRepository struct {
	mock.Mock
}

func (m *MockMailRepository) Create(ctx context.Context, mail *domain.Mail) error {
	args := m.Called(ctx, mail)
	return args.Error(0)
}

func (m *MockMailRepository) GetByID(ctx context.Context, sessionID, id string) (*domain.Mail, error) {
	args := m.Called(ctx, sessionID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Mail), args.Error(1)
}

func (m *MockMailRepository) GetInbox(ctx context.Context, sessionID, recipientID string, page, pageSize int, query string) ([]domain.Mail, int64, error) {
	args := m.Called(ctx, sessionID, recipientID, page, pageSize, query)
	return args.Get(0).([]domain.Mail), args.Get(1).(int64), args.Error(2)
}

func (m *MockMailRepository) GetSent(ctx context.Context, sessionID, senderID string, page, pageSize int, query string) ([]domain.Mail, int64, error) {
	args := m.Called(ctx, sessionID, senderID, page, pageSize, query)
	return args.Get(0).([]domain.Mail), args.Get(1).(int64), args.Error(2)
}

func (m *MockMailRepository) UpdateStatus(ctx context.Context, mailID, recipientID, status string) error {
	args := m.Called(ctx, mailID, recipientID, status)
	return args.Error(0)
}

func (m *MockMailRepository) DeleteForSender(ctx context.Context, mailID string) error {
	args := m.Called(ctx, mailID)
	return args.Error(0)
}

func (m *MockMailRepository) DeleteSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockMailRepository) GetAttachmentByID(ctx context.Context, sessionID, id string) (*domain.Attachment, error) {
	args := m.Called(ctx, sessionID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attachment), args.Error(1)

}
func (m *MockMailRepository) CreateChatMessage(ctx context.Context, msg *domain.ChatMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockMailRepository) GetChatHistory(ctx context.Context, sessionID, userA, userB string, limit int) ([]domain.ChatMessage, error) {
	args := m.Called(ctx, sessionID, userA, userB, limit)
	return args.Get(0).([]domain.ChatMessage), args.Error(1)
}

func (m *MockMailRepository) MarkChatAsRead(ctx context.Context, sessionID, senderID, receiverID string) error {
	args := m.Called(ctx, sessionID, senderID, receiverID)
	return args.Error(0)
}

func (m *MockMailRepository) GetUnreadMailCount(ctx context.Context, sessionID, userID string) (int64, error) {
	args := m.Called(ctx, sessionID, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMailRepository) GetIMUnreadCounts(ctx context.Context, sessionID, userID string) (map[string]int, error) {
	args := m.Called(ctx, sessionID, userID)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockMailRepository) GetOrphanSessionIDs(ctx context.Context, activeSessionIDs []string) ([]string, error) {
	args := m.Called(ctx, activeSessionIDs)
	return args.Get(0).([]string), args.Error(1)
}

// MockStorageService
type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) UploadFile(ctx context.Context, sessionID, fileName string, content io.Reader) (string, error) {
	args := m.Called(ctx, sessionID, fileName, content)
	return args.String(0), args.Error(1)
}

func (m *MockStorageService) GetFile(ctx context.Context, path string) (io.ReadCloser, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageService) DeleteSessionDir(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockStorageService) DeleteFile(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

// Ensure interfaces are implemented
var _ ports.MailRepository = (*MockMailRepository)(nil)
var _ ports.StorageService = (*MockStorageService)(nil)
