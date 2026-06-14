package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"qfis/internal/api"
	"qfis/internal/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

func main() {
	// ── Production defaults ──────────────────────────────────
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ── Database ─────────────────────────────────────────────
	dbPath := "data/qfis.db"
	if p := os.Getenv("QFIS_DB_PATH"); p != "" {
		dbPath = p
	}

	// Ensure data directory exists
	if dir := strings.TrimSuffix(dbPath, "/qfis.db"); dir != dbPath {
		os.MkdirAll(dir, 0755)
	}

	if err := db.Init(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	count, _ := db.GetMerchantCount()
	if count == 0 {
		if err := db.Seed(); err != nil {
			log.Printf("Seed: %v", err)
		}
	}

	// ── Router ───────────────────────────────────────────────
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// ── CORS ─────────────────────────────────────────────────
	allowedOrigins := []string{"*"}
	if origins := os.Getenv("CORS_ORIGINS"); origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-API-Key"},
		AllowCredentials: true,
	}))

	// ── API Key Auth (optional) ──────────────────────────────
	apiKey := os.Getenv("API_KEY")
	if apiKey != "" {
		r.Use(func(c *gin.Context) {
			// Skip auth for CORS preflight and frontend static
			if c.Request.Method == "OPTIONS" || !strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.Next()
				return
			}
			key := c.GetHeader("X-API-Key")
			if key != apiKey {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
				return
			}
			c.Next()
		})
	}

	// ── Rate Limiting Hint ───────────────────────────────────
	// For production behind nginx, add:
	//   limit_req zone=qfis burst=20 nodelay;
	// Or use: github.com/ulule/limiter/v3

	// ── API Routes ───────────────────────────────────────────
	v1 := r.Group("/api/v1")
	{
		v1.POST("/merchant/check", api.CheckMerchant)
		v1.GET("/merchant/search", api.SearchMerchant)
		v1.GET("/merchant/:id", api.GetMerchantByID)

		v1.POST("/report/submit", api.SubmitReport)
		v1.GET("/report/list", api.ListReports)
		v1.GET("/report/:id", api.GetReport)

		v1.GET("/dashboard/stats", api.GetDashboardStats)
		v1.GET("/dashboard/network", api.GetMerchantNetwork)
	}

	// ── Frontend ─────────────────────────────────────────────
	serveFrontend(r)

	// ── Start ────────────────────────────────────────────────
	port := "21465"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("QFIS starting on :%s", port)
	r.Run(":" + port)
}

func serveFrontend(r *gin.Engine) {
	// Try embedded frontend first
	sub, err := fs.Sub(frontendFS, "frontend/dist")
	if err == nil {
		r.GET("/", frontendHandler(sub, ""))
		r.GET("/assets/*path", frontendHandler(sub, "assets"))
		r.GET("/favicon.svg", frontendHandler(sub, "favicon.svg"))
		r.NoRoute(frontendHandler(sub, ""))
		log.Println("Frontend: embedded")
		return
	}

	// Fallback: local dev filesystem
	localDist := "../frontend/dist"
	if info, err := os.Stat(localDist); err == nil && info.IsDir() {
		r.Static("/assets", localDist+"/assets")
		r.StaticFile("/favicon.svg", localDist+"/favicon.svg")
		r.NoRoute(func(c *gin.Context) { c.File(localDist + "/index.html") })
		log.Println("Frontend: filesystem")
	} else {
		log.Println("Frontend: none (API only)")
	}
}

func frontendHandler(sub fs.FS, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		filePath := strings.TrimPrefix(c.Request.URL.Path, "/")
		if prefix != "" && filePath == prefix {
			filePath = prefix
		}
		if prefix != "" && filePath != prefix {
			filePath = strings.TrimPrefix(filePath, prefix+"/")
			filePath = prefix + "/" + filePath
		}

		data, err := fs.ReadFile(sub, filePath)
		if err != nil {
			// SPA fallback: serve index.html
			data, err = fs.ReadFile(sub, "index.html")
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
		}

		contentType := detectContentType(filePath, data)
		c.Data(http.StatusOK, contentType, data)
	}
}

func detectContentType(path string, data []byte) string {
	switch {
	case strings.HasSuffix(path, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(path, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(path, ".js"):
		return "application/javascript"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".png"):
		return "image/png"
	case strings.HasSuffix(path, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(path, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(path, ".json"):
		return "application/json"
	default:
		return http.DetectContentType(data)
	}
}
