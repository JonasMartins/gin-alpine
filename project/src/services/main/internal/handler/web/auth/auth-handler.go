package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gin-alpine/src/pkg/utils"

	authDomain "gin-alpine/src/internal/domain/auth"
	"gin-alpine/src/internal/infra/redis"
	"gin-alpine/src/internal/usecases"
	"gin-alpine/src/services/web"

	handler "gin-alpine/src/services/main/internal/handler/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	csrf "github.com/utrack/gin-csrf"
	"go.uber.org/zap"
)

type AuthHandler struct {
	Logger       *zap.Logger
	Renderer     *web.Renderer
	T            *utils.Translator
	RedisDB      *redis.RedisClient
	authUsecases *usecases.AuthUsecases
}

func NewAuthHandler(
	authUseCases *usecases.AuthUsecases,
	redisDB *redis.RedisClient,
	renderer *web.Renderer,
	logger *zap.Logger,
	t *utils.Translator,
) *AuthHandler {
	return &AuthHandler{
		authUsecases: authUseCases,
		Renderer:     renderer,
		RedisDB:      redisDB,
		Logger:       logger,
		T:            t,
	}
}

func (h *AuthHandler) LoginPostWeb(c *gin.Context) {
	var loginDTO LoginInputDTO
	if err := c.ShouldBind(&loginDTO); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			messages := utils.CustomErrorTranslator(validationErrors, GetLoginCustomMessages(h.T))
			for field, msg := range messages {
				if errRender := h.Renderer.Render(c.Writer, "auth", "login", gin.H{
					"error": handler.HTTPError{
						Field:   field,
						Message: msg,
					},
					"csrf": csrf.GetToken(c),
				}); errRender != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
				return
			}
		}
		return
	}
	u, err := h.authUsecases.Login(
		c.Request.Context(),
		authDomain.LoginInput{
			Email:    loginDTO.Email,
			Password: loginDTO.Password,
		})
	if err != nil {
		code, httpErr := handler.MapAuthErrorToHTTP(err, h.T)
		h.Logger.Info("login_fail", zap.String("message", httpErr.Message), zap.String("email", loginDTO.Email))
		c.Status(code)
		if errRender := h.Renderer.Render(c.Writer, "auth", "login", gin.H{
			"error": httpErr,
			"csrf":  csrf.GetToken(c),
		}); errRender != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	session := sessions.Default(c)
	session.Clear()
	session.Set("user_id", u.ID)
	if err = session.Save(); err != nil {
		log.Printf("error saving session: %v", err)
		c.Status(http.StatusInternalServerError)
		if errRender := h.Renderer.Render(c.Writer, "auth", "login", gin.H{
			"Error": "We couldn’t start your session. Please try again.",
		}); errRender != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	err = authDomain.StoreUserAuth(
		c.Request.Context(),
		h.RedisDB,
		authDomain.UserAuth{
			ID:    u.ID,
			Email: u.Email,
			Name:  u.Name,
			Role:  authDomain.Role(int(u.RoleID)),
		},
		24*time.Hour,
	)
	if err != nil {
		h.Logger.Error("auth_store_fail", zap.String("message", err.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.Logger.Info("login_success:", zap.String("email", loginDTO.Email))
	c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthHandler) LogoutPostWeb(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	// Remove all session values
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	// clean from redis cache
	if userID != nil {
		key := fmt.Sprintf("auth:user:%v", userID)
		_ = h.RedisDB.Client.Del(c.Request.Context(), key).Err()
	}
	// Persist changes and expire cookie
	err := session.Save()
	if err != nil {
		log.Printf("error saving session %v", err)
	}
	// Redirect user
	h.Logger.Info("logout_success:", zap.Any("user_id", userID))
	c.Redirect(http.StatusFound, "/login")
}
