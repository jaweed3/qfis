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
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("QFIS starting on :%s", port)
	r.Run(":" + port)
}

func serveFrontend(r *gin.Engine) {
	if f, err := frontendFS.Open("frontend/dist/index.html"); err == nil {
		f.Close()
		sub, err := fs.Sub(frontendFS, "frontend/dist")
		if err != nil {
			log.Fatalf("Failed to get frontend sub FS: %v", err)
		}
		staticFS := http.FS(sub)

		r.GET("/", func(c *gin.Context) { c.FileFromFS("index.html", staticFS) })
		r.GET("/assets/*filepath", func(c *gin.Context) {
			c.FileFromFS(strings.TrimPrefix(c.Request.URL.Path, "/"), staticFS)
		})
		r.GET("/favicon.svg", func(c *gin.Context) { c.FileFromFS("favicon.svg", staticFS) })
		r.NoRoute(func(c *gin.Context) {
			p := strings.TrimPrefix(c.Request.URL.Path, "/")
			if f, err := sub.Open(p); err == nil {
				f.Close()
				c.FileFromFS(p, staticFS)
				return
			}
			c.FileFromFS("index.html", staticFS)
		})
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
