package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"raven"
	"raven/internal/core/domain"
	"raven/internal/handler"
	"raven/internal/infrastructure/storage"
	"raven/internal/repository"
	"raven/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	var port int
	var ooHost string
	var defUser string

	flag.IntVar(&port, "port", 8080, "server port")
	flag.StringVar(&ooHost, "oo-host", os.Getenv("ONLYOFFICE_HOST"), "OnlyOffice server host (e.g. 192.168.1.100:8090)")
	flag.StringVar(&defUser, "default-user", os.Getenv("DEFAULT_USER_ID"), "Default User ID for simulation")
	flag.Parse()

	if ooHost == "" {
		ooHost = "localhost:8090" // Generic fallback
	}
	if defUser == "" {
		defUser = "guest"
	}

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
	mailHandler := handler.NewMailHandler(mailService, store, ooHost, defUser)

	// 4. Gin Router
	r := gin.Default()

	// Generic CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Session-ID")
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
			mails.DELETE("/:id", mailHandler.DeleteMail)
			mails.GET("/download", mailHandler.DownloadAttachment)
			mails.GET("/events", mailHandler.StreamNotifications)
		}
		api.GET("/onlyoffice/template", mailHandler.ServeOnlyOfficeTemplate)
		api.POST("/onlyoffice/callback", mailHandler.OnlyOfficeCallback)
		api.POST("/onlyoffice/forcesave", mailHandler.OnlyOfficeForceSave)
		api.DELETE("/sessions/:id", mailHandler.DeleteSession)
	}

	// --- 5. 静态前端资源托管 (内嵌) ---
	// 将嵌入的 FS 子目录 web/dist 提取出来
	distFS, err := fs.Sub(raven.FrontendDist, "web/dist")
	if err == nil {
		// 1. 托管静态资源目录 (如 assets)
		// 注意：Vite 打包通常包含 assets 文件夹
		// 如果 dist 根目录下有 favicon.ico 等，可以通过 StaticFile 托管或通用逻辑处理

		// 2. SPA 全路由拦截：非 API 请求且文件不存在时，返回 index.html
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// 如果是 API 请求，直接返回 404
			if strings.HasPrefix(path, "/api") {
				return
			}

			// 尝试从嵌入文件系统中读取具体文件 (如 /favicon.ico)
			cleanPath := strings.TrimPrefix(path, "/")
			if cleanPath == "" {
				cleanPath = "index.html"
			}

			f, err := distFS.Open(cleanPath)
			if err == nil {
				defer f.Close()
				// 获取文件信息以确定 Content-Type
				stat, _ := f.Stat()
				if !stat.IsDir() {
					// 简单起见，使用 StaticFS 的逻辑包装或直接输出
					http.FileServer(http.FS(distFS)).ServeHTTP(c.Writer, c.Request)
					return
				}
			}

			// 如果文件不存在且不是 API，兜底返回 index.html (支持 Vue Router)
			fIndex, err := distFS.Open("index.html")
			if err != nil {
				c.String(404, "Frontend assets not found. Did you run 'npm run build'?")
				return
			}
			defer fIndex.Close()
			content, _ := io.ReadAll(fIndex)
			c.Data(200, "text/html; charset=utf-8", content)
		})
	}

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server starting on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
