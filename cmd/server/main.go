package main

import (
	"log"
	"raven/internal/core/domain"
	"raven/internal/handler"
	"raven/internal/infrastructure/storage"
	"raven/internal/repository"
	"raven/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. DB Setup using SQLite
	db, err := gorm.Open(sqlite.Open("raven.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto Migrate
	if err := db.AutoMigrate(&domain.Mail{}, &domain.MailRecipient{}, &domain.Attachment{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 2. Storage Setup
	store, err := storage.NewLocalStorage("./uploads")
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	// 3. Application Layers
	mailRepo := repository.NewMailRepository(db)
	mailService := service.NewMailService(mailRepo, store)
	mailHandler := handler.NewMailHandler(mailService, store)

	// 4. Gin Router
	r := gin.Default()

	// Generic CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := r.Group("/api/v1")
	{
		mails := api.Group("/mails")
		{
			mails.POST("/send", mailHandler.SendMail)
			mails.GET("/inbox", mailHandler.GetInbox)
			mails.GET("/sent", mailHandler.GetSent)
			mails.GET("/:id", mailHandler.GetMail)
			mails.GET("/download", mailHandler.DownloadAttachment)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
