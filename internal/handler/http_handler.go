package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"raven/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type MailHandler struct {
	service ports.MailService
	storage ports.StorageService // Need storage to download files
}

func NewMailHandler(service ports.MailService, storage ports.StorageService) *MailHandler {
	return &MailHandler{
		service: service,
		storage: storage,
	}
}

// SendMail handles sending a new mail with attachments
func (h *MailHandler) SendMail(c *gin.Context) {
	// Multipart form
	subject := c.PostForm("subject")
	content := c.PostForm("content")
	to := strings.Split(c.PostForm("to"), ",")
	cc := strings.Split(c.PostForm("cc"), ",")
	bcc := strings.Split(c.PostForm("bcc"), ",") // Optional

	// Attachments
	form, _ := c.MultipartForm()
	files := form.File["attachments"]

	var attachmentReqs []ports.AttachmentRequest
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open attachment"})
			return
		}
		defer f.Close()

		attachmentReqs = append(attachmentReqs, ports.AttachmentRequest{
			FileName: file.Filename,
			Content:  f,
			Size:     file.Size,
			MimeType: file.Header.Get("Content-Type"),
		})
	}

	req := ports.SendMailRequest{
		Subject:     subject,
		Content:     content,
		To:          filterEmpty(to),
		Cc:          filterEmpty(cc),
		Bcc:         filterEmpty(bcc),
		Attachments: attachmentReqs,
	}

	// Mock Sender ID (In real app, get from Context/Token)
	senderID := c.Query("user_id") // Temporary for simulation
	if senderID == "" {
		senderID = "user-123"
	}

	mail, err := h.service.SendMail(c.Request.Context(), senderID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mail)
}

func (h *MailHandler) GetInbox(c *gin.Context) {
	userID := c.Query("user_id")
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	mails, total, err := h.service.GetInbox(c.Request.Context(), userID, page, pageSize, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": mails, "total": total, "page": page, "page_size": pageSize})
}

func (h *MailHandler) GetSent(c *gin.Context) {
	userID := c.Query("user_id")
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	mails, total, err := h.service.GetSent(c.Request.Context(), userID, page, pageSize, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": mails, "total": total, "page": page, "page_size": pageSize})
}

func (h *MailHandler) DeleteMail(c *gin.Context) {
	id := c.Param("id")
	userID := c.Query("user_id") // In real app, from context

	if err := h.service.DeleteMail(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *MailHandler) GetMail(c *gin.Context) {
	id := c.Param("id")
	userID := c.Query("user_id")

	mail, err := h.service.ReadMail(c.Request.Context(), userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mail)
}

func (h *MailHandler) DownloadAttachment(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		// Fallback for legacy calls if any (though we are changing the contract)
		path := c.Query("path")
		if path != "" {
			// Legacy insecure mode (deprecated)
			h.downloadByPath(c, path)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Attachment ID required"})
		return
	}

	att, err := h.service.GetAttachment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	f, err := h.storage.GetFile(c.Request.Context(), att.FilePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File content not found"})
		return
	}
	defer f.Close()

	// Use original filename
	filename := att.FileName
	encodedFilename := url.QueryEscape(filename)
	// Replace + with %20 for space compatibility
	encodedFilename = strings.ReplaceAll(encodedFilename, "+", "%20")

	// Standard approach: filename="ascii_only_fallback"; filename*=UTF-8''url_encoded
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", encodedFilename, encodedFilename))

	c.Header("Content-Type", att.MimeType)
	if att.MimeType == "" {
		c.Header("Content-Type", "application/octet-stream")
	}

	if _, err := io.Copy(c.Writer, f); err != nil {
		// Log error
	}
}

func (h *MailHandler) downloadByPath(c *gin.Context, path string) {
	f, err := h.storage.GetFile(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer f.Close()

	fileName := filepath.Base(path)
	// Legacy fallback, simple disposition
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	io.Copy(c.Writer, f)
}

func filterEmpty(s []string) []string {
	var r []string
	for _, v := range s {
		if strings.TrimSpace(v) != "" {
			r = append(r, strings.TrimSpace(v))
		}
	}
	return r
}
