package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"raven/internal/core/domain"
	"raven/internal/core/ports"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMailService_SendMail(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success without attachments", func(t *testing.T) {
		mockRepo := new(MockMailRepository)
		mockStorage := new(MockStorageService)
		svc := NewMailService(mockRepo, mockStorage)

		req := ports.SendMailRequest{
			SessionID:   "session-1",
			Subject:     "Hello",
			Content:     "World",
			ContentType: "text",
			To:          []string{"user-2"},
		}

		mockRepo.On("Create", ctx, mock.MatchedBy(func(m *domain.Mail) bool {
			return m.Subject == "Hello" && m.SessionID == "session-1" && len(m.Recipients) == 1
		})).Return(nil)

		mail, err := svc.SendMail(ctx, "user-1", req)

		assert.NoError(t, err)
		assert.NotNil(t, mail)
		assert.Equal(t, "Hello", mail.Subject)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success with attachments", func(t *testing.T) {
		mockRepo := new(MockMailRepository)
		mockStorage := new(MockStorageService)
		svc := NewMailService(mockRepo, mockStorage)

		fileContent := strings.NewReader("fake content")
		req := ports.SendMailRequest{
			SessionID: "session-1",
			To:        []string{"user-2"},
			Attachments: []ports.AttachmentRequest{
				{FileName: "test.txt", Content: fileContent, Size: 12, MimeType: "text/plain"},
			},
		}

		// Expect upload
		mockStorage.On("UploadFile", ctx, "session-1", "test.txt", fileContent).Return("path/to/test.txt", nil)

		// Expect create
		mockRepo.On("Create", ctx, mock.MatchedBy(func(m *domain.Mail) bool {
			return len(m.Attachments) == 1 && m.Attachments[0].FilePath == "path/to/test.txt"
		})).Return(nil)

		mail, err := svc.SendMail(ctx, "user-1", req)

		assert.NoError(t, err)
		assert.Len(t, mail.Attachments, 1)
		mockStorage.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Rollback on DB Failure", func(t *testing.T) {
		mockRepo := new(MockMailRepository)
		mockStorage := new(MockStorageService)
		svc := NewMailService(mockRepo, mockStorage)

		fileContent := strings.NewReader("fake content")
		req := ports.SendMailRequest{
			SessionID: "session-1",
			To:        []string{"user-2"},
			Attachments: []ports.AttachmentRequest{
				{FileName: "test.txt", Content: fileContent},
			},
		}

		// 1. Storage Upload Success
		mockStorage.On("UploadFile", ctx, "session-1", "test.txt", fileContent).Return("path/to/test.txt", nil)

		// 2. Repo Create Fails
		dbErr := errors.New("db connection lost")
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Mail")).Return(dbErr)

		// 3. Expect Rollback (DeleteFile)
		mockStorage.On("DeleteFile", ctx, "path/to/test.txt").Return(nil)

		mail, err := svc.SendMail(ctx, "user-1", req)

		assert.Error(t, err)
		assert.Nil(t, mail)
		assert.Equal(t, dbErr, err)

		mockStorage.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Upload Failure", func(t *testing.T) {
		mockRepo := new(MockMailRepository)
		mockStorage := new(MockStorageService)
		svc := NewMailService(mockRepo, mockStorage)

		fileContent := strings.NewReader("fake")
		req := ports.SendMailRequest{
			SessionID: "session-1",
			Attachments: []ports.AttachmentRequest{
				{FileName: "test.txt", Content: fileContent},
			},
		}

		// Storage Fails
		mockStorage.On("UploadFile", ctx, "session-1", "test.txt", fileContent).Return("", errors.New("disk full"))

		// Should NOT call Repo.Create
		// Should NOT call DeleteFile

		_, err := svc.SendMail(ctx, "user-1", req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "disk full")

		mockRepo.AssertNotCalled(t, "Create")
		mockStorage.AssertExpectations(t)
	})
}

func TestMailService_SyncSessions(t *testing.T) {
	ctx := context.TODO()

	t.Run("Sync and Clean", func(t *testing.T) {
		mockRepo := new(MockMailRepository)
		mockStorage := new(MockStorageService)
		svc := NewMailService(mockRepo, mockStorage)

		activeIDs := []string{"active-1"}
		orphans := []string{"orphan-1", "orphan-2"}

		// 1. Get Orphans
		mockRepo.On("GetOrphanSessionIDs", ctx, activeIDs).Return(orphans, nil)

		// 2. Delete each orphan (Repo + Storage)
		// Orphan-1
		mockRepo.On("DeleteSession", ctx, "orphan-1").Return(nil)
		mockStorage.On("DeleteSessionDir", ctx, "orphan-1").Return(nil)
		// Orphan-2
		mockRepo.On("DeleteSession", ctx, "orphan-2").Return(nil)
		mockStorage.On("DeleteSessionDir", ctx, "orphan-2").Return(nil)

		count, err := svc.SyncSessions(ctx, activeIDs)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}
