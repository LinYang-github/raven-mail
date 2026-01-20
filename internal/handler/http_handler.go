package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"raven/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type MailHandler struct {
	service         ports.MailService
	storage         ports.StorageService
	OnlyOfficeHost  string
	DefaultSenderID string
}

func NewMailHandler(service ports.MailService, storage ports.StorageService, onlyOfficeHost string, defaultSenderID string) *MailHandler {
	return &MailHandler{
		service:         service,
		storage:         storage,
		OnlyOfficeHost:  onlyOfficeHost,
		DefaultSenderID: defaultSenderID,
	}
}

// SendMail handles sending a new mail with attachments
func (h *MailHandler) SendMail(c *gin.Context) {
	// Multipart form
	subject := c.PostForm("subject")
	content := c.PostForm("content")
	contentType := c.DefaultPostForm("content_type", "text")
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

	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = "default"
	}

	req := ports.SendMailRequest{
		SessionID:   sessionID,
		Subject:     subject,
		Content:     content,
		ContentType: contentType,
		To:          filterEmpty(to),
		Cc:          filterEmpty(cc),
		Bcc:         filterEmpty(bcc),
		Attachments: attachmentReqs,
	}

	// Mock Sender ID (In real app, get from Context/Token)
	senderID := c.Query("user_id") // Temporary for simulation
	if senderID == "" {
		senderID = h.DefaultSenderID
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

	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = "default"
	}

	mails, total, err := h.service.GetInbox(c.Request.Context(), sessionID, userID, page, pageSize, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": mails, "total": total, "page": page, "page_size": pageSize, "session_id": sessionID})
}

func (h *MailHandler) GetSent(c *gin.Context) {
	userID := c.Query("user_id")
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = "default"
	}

	mails, total, err := h.service.GetSent(c.Request.Context(), sessionID, userID, page, pageSize, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": mails, "total": total, "page": page, "page_size": pageSize, "session_id": sessionID})
}

func (h *MailHandler) DeleteMail(c *gin.Context) {
	id := c.Param("id")
	userID := c.Query("user_id") // In real app, from context
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = "default"
	}

	if err := h.service.DeleteMail(c.Request.Context(), sessionID, userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *MailHandler) GetMail(c *gin.Context) {
	id := c.Param("id")
	userID := c.Query("user_id")
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = "default"
	}

	mail, err := h.service.ReadMail(c.Request.Context(), sessionID, userID, id)
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

	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = "default"
	}

	att, err := h.service.GetAttachment(c.Request.Context(), sessionID, id)
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

	// Determine disposition type
	dispositionType := "attachment"
	if c.Query("disposition") == "inline" {
		dispositionType = "inline"
	}

	// Standard approach: filename="ascii_only_fallback"; filename*=UTF-8''url_encoded
	c.Header("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"; filename*=UTF-8''%s", dispositionType, encodedFilename, encodedFilename))

	c.Header("Content-Type", att.MimeType)
	// If inline and no mime type, try to guess or default to octet-stream (which browser will download anyway)
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

func (h *MailHandler) ServeOnlyOfficeTemplate(c *gin.Context) {
	key := c.Query("key")
	sessionID := c.Query("session_id") // ONLYOFFICE components will pass this via query
	if sessionID == "" {
		sessionID = c.GetHeader("X-Session-ID")
	}
	if sessionID == "" {
		sessionID = "default"
	}

	// 从针对场次隔离的数据存储目录读取
	filePath := fmt.Sprintf("./data/%s/docs/%s.docx", sessionID, key)

	fmt.Printf("[OnlyOffice] Document requested for key: %s\n", key)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", "attachment; filename=document.docx")

	// 如果该 Key 对应的已保存文件存在，则返回，否则返回原始空白模板
	if _, err := os.Stat(filePath); err == nil {
		c.File(filePath)
	} else {
		c.File("./templates/empty.docx")
	}
}

func (h *MailHandler) OnlyOfficeForceSave(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(400, gin.H{"error": "missing key"})
		return
	}

	// 从配置中获取 ONLYOFFICE 服务器地址
	cmdURL := fmt.Sprintf("http://%s/coauthoring/CommandService.ashx", h.OnlyOfficeHost)

	payload := map[string]interface{}{
		"c":   "forcesave",
		"key": key,
	}

	jsonPayload, _ := json.Marshal(payload)
	resp, err := http.Post(cmdURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("[OnlyOffice] ForceSave trigger failed (check ONLYOFFICE_HOST): %v\n", err)
		c.JSON(200, gin.H{"error": 1, "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	fmt.Printf("[OnlyOffice] ForceSave triggered for key: %s\n", key)
	c.JSON(200, gin.H{"error": 0, "message": "forcesave triggered"})
}

func (h *MailHandler) OnlyOfficeCallback(c *gin.Context) {
	var body map[string]interface{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": 1})
		return
	}

	status := body["status"].(float64)
	fmt.Printf("[OnlyOffice] Callback received. Status: %v\n", status)

	// Status 2: Ready for saving (closed)
	// Status 6: Being edited, but state saved (forcesave)
	if status == 2 || status == 6 {
		downloadURL := body["url"].(string)
		key := body["key"].(string)
		sessionID := c.Query("session_id")
		if sessionID == "" {
			sessionID = "default"
		}

		fmt.Printf("[OnlyOffice] Saving document (status %v): %s, URL: %s, Session: %s\n", status, key, downloadURL, sessionID)

		// 下载文件
		resp, err := http.Get(downloadURL)
		if err != nil {
			fmt.Printf("[OnlyOffice] Download failed: %v\n", err)
			c.JSON(200, gin.H{"error": 1})
			return
		}
		defer resp.Body.Close()

		// 确保针对场次的数据存储目录存在
		storageDir := fmt.Sprintf("./data/%s/docs", sessionID)
		os.MkdirAll(storageDir, 0755)

		// 保存到针对场次隔离的文档存储目录
		filePath := filepath.Join(storageDir, fmt.Sprintf("%s.docx", key))
		out, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("[OnlyOffice] Create file failed: %v\n", err)
			c.JSON(200, gin.H{"error": 1})
			return
		}
		defer out.Close()

		io.Copy(out, resp.Body)
		fmt.Printf("[OnlyOffice] Document %s (status %v) saved successfully to %s\n", key, status, filePath)
	}

	c.JSON(200, gin.H{"error": 0})
}

func (h *MailHandler) StreamNotifications(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	ch := h.service.Subscribe()
	defer h.service.Unsubscribe(ch)

	// Ping to keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-ch:
			if !ok {
				return false
			}
			c.SSEvent("message", msg)
			return true
		case <-ticker.C:
			c.SSEvent("ping", "keep-alive")
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}
