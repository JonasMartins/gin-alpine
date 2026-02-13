// Package router ...
package router

import (
	"net/http"
	"time"

	"gin-alpine/src/services/main/internal/bootstrap"
	"gin-alpine/src/services/main/internal/handler/middleware"

	"io/fs"
	"log"

	"gin-alpine/src/services/web"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type ChartData struct {
	Labels []string `json:"labels"`
	Values []int    `json:"values"`
}

func NewRouter(b *bootstrap.Bootstrap) *gin.Engine {
	gin.SetMode(b.Config.Env)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte("very-secret-key"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   b.Config.Env == gin.ReleaseMode, // enable in production (HTTPS)
	})

	csrfMiddleware := csrf.Middleware(csrf.Options{
		Secret: b.Config.CSRFSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(http.StatusForbidden, "CSRF token mismatch")
			c.Abort()
		},
	})

	r.Use(sessions.Sessions("session", store))
	r.Use(middleware.RateLimitMiddleware(b))
	// Then csrf
	r.Use(csrfMiddleware)
	r.Use(middleware.CSRFTpl())

	staticFS, err := fs.Sub(web.StaticFiles, "static")
	if err != nil {
		log.Fatalf("error creating static sub filesystem: %v", err)
	}
	r.StaticFS("/static", http.FS(staticFS))

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("static/favicon.ico", http.FS(staticFS))
	})

	r.NoRoute(func(c *gin.Context) {
		if err := b.Renderer.Render(c.Writer, "base", "404", gin.H{
			"Title": "Page not found",
		}); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	})
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Println("panic:", recovered)
		if err := b.Renderer.Render(c.Writer, "base", "500", gin.H{
			"Title": "Internal Server Error",
		}); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}))

	// PUBLIC ROUTES
	public := r.Group("/")
	public.GET("/login", b.Renderer.Page("auth", "login", func(c *gin.Context) gin.H {
		return gin.H{"Title": "Login"}
	}))
	public.POST("/login", b.AuthWebHandler.LoginPostWeb)

	// PROTECTED ROUTES
	protected := r.Group("/")
	protected.Use(middleware.AuthWeb(b.RedisDB))
	protected.POST("/logout", b.AuthWebHandler.LogoutPostWeb)
	protected.GET("/", b.Renderer.Page("main", "home", func(c *gin.Context) gin.H {
		return gin.H{"Title": "Home"}
	}))
	// API ROUTES
	// protected.GET("/api/patterns/:id/draws", b.PatternsHTTPHandler.ListPatternAndDrawsHTTP)

	// ONLY FOR DEV USERS
	// dev := r.Group("/api/system-admin/")
	// dev.Use(middleware.AuthWeb(b.RedisDB))
	// dev.Use(middleware.RequireRoleAtLeast(auth.RoleDev))
	// dev.POST("/draws", b.DrawsHTTPHandler.AddMostRecentDrawAndSyncHTTP)

	return r
}
