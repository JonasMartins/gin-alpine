package middleware

import (
	"fmt"
	"net/http"

	"gin-alpine/src/internal/domain/auth"
	"gin-alpine/src/internal/infra/redis"
	"gin-alpine/src/pkg/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	csrf "github.com/utrack/gin-csrf"
)

func AuthWeb(redis *redis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		key := fmt.Sprintf("auth:user:%v", userID)

		var user auth.UserAuth
		var appData utils.AppData
		err := redis.Cache.Get(c.Request.Context(), key, &user)
		errAppData := redis.Cache.Get(c.Request.Context(), redis.GetAppDataKey(), &appData)
		if err != nil || errAppData != nil {
			// must have initial data always available
			// session exists but auth cache is gone → force relogin
			session.Clear()
			_ = session.Save()
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// attach to context
		c.Set("auth_user", user)
		c.Set("app_data", appData)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) int64 {
	return c.MustGet("user_id").(int64)
}

func GetAppData(c *gin.Context) utils.AppData {
	return c.MustGet("app_data").(utils.AppData)
}

func CSRFTpl() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("csrf", csrf.GetToken(c))
		c.Next()
	}
}

// RequireRoleAtLeast usage:
// admin := web.Group("/admin")
// admin.Use(RequireRoleAtLeast(auth.RoleAdmin))
func RequireRoleAtLeast(min auth.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exists := c.Get("auth_user")
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := v.(auth.UserAuth)

		if user.Role < min {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// RequireAnyRole usage:
// manager := web.Group("/manager")
// manager.Use(RequireAnyRole(
//
//	auth.RoleManager,
//	auth.RoleAdmin,
//
// ))
// manager.GET("/reports", reportsHandler)
func RequireAnyRole(allowed ...auth.Role) gin.HandlerFunc {
	allowedSet := make(map[auth.Role]struct{}, len(allowed))
	for _, r := range allowed {
		allowedSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		v, exists := c.Get("auth_user")
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := v.(auth.UserAuth)

		if _, ok := allowedSet[user.Role]; !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
