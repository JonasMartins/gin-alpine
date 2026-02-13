// Package bootstrap
package bootstrap

import (
	"context"
	"gin-alpine/src/internal/configs"
	"gin-alpine/src/pkg/utils"
	"gin-alpine/src/services/web"

	"runtime"
	"time"

	"gin-alpine/src/internal/infra/postgres"
	"gin-alpine/src/internal/infra/redis"
	"gin-alpine/src/internal/usecases"
	authWebHandler "gin-alpine/src/services/main/internal/handler/web/auth"
	"gin-alpine/src/services/main/internal/jobs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v9"
	"github.com/robfig/cron/v3"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Bootstrap struct {
	Job            *jobs.Job
	Cron           *cron.Cron
	Logger         *zap.Logger
	Renderer       *web.Renderer
	Config         *configs.Config
	Translator     *utils.Translator
	RedisDB        *redis.RedisClient
	RollingLogger  *lumberjack.Logger
	AuthWebHandler *authWebHandler.AuthHandler
}

func MustGetBootstrapInstance() *Bootstrap {
	cfg, err := configs.LoadConfig()
	if err != nil {
		utils.FatalResult("failed to load config", err)
	}
	repo, err := postgres.NewPgRepository(cfg)
	if err != nil || repo == nil {
		utils.FatalResult("failed to load sqlite repository", err)
	}
	translator := utils.NewTranslator("pt")
	mode := gin.DebugMode
	if cfg.Env == "production" {
		mode = gin.ReleaseMode
	}
	b := Bootstrap{
		Config:     cfg,
		Cron:       cron.New(),
		Translator: translator,
		Renderer:   web.NewRenderer(mode),
		RedisDB:    redis.NewRedisClient(cfg.RedisURL),
	}
	b.setUpLogger()
	b.setupJobs()

	// REPOSITORIES
	userRepository := postgres.NewUserRepository(repo)

	// USECASES
	authUsecases := usecases.NewAuthUsecases(userRepository)

	// HANDLERS
	authHandler := authWebHandler.NewAuthHandler(authUsecases, b.RedisDB, b.Renderer, b.Logger, translator)

	// set bootstrap handlers
	b.AuthWebHandler = authHandler

	return &b
}

func (b *Bootstrap) setupJobs() {
	b.Job = jobs.NewJob(
		b.Cron,
		b.Logger,
	)
}

func (b *Bootstrap) setUpLogger() {
	if b.RollingLogger == nil {
		logFile, err := utils.GetDefaultLogsFileName()
		if err != nil {
			utils.FatalResult("unable to set logs file", err)
		}
		b.RollingLogger = &lumberjack.Logger{
			Filename:   logFile, // Will rotate daily
			MaxSize:    10,      // Max size in MB
			MaxBackups: 7,       // Keep 7 old logs
			MaxAge:     30,      // Keep logs for 30 days
			Compress:   false,   // Compress old logs
		}
	}
	if b.Logger == nil {
		logWriter := zapcore.AddSync(b.RollingLogger)
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), logWriter, zap.InfoLevel)
		b.Logger = zap.New(core, zap.AddCaller())
	}
	defer func(Logger *zap.Logger) {
		err := Logger.Sync()
		if err != nil {
			utils.FatalResult("unable to defer logger sync %v", err)
		}
	}(b.Logger)
}

func (b *Bootstrap) LogSystemUsage() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		// cpuPercent, _ := cpu.Percent(0, false) // Get CPU usage
		b.Logger.Info("SYSTEM_USAGE",
			// zap.Float64("CPU (%)", cpuPercent[0]),
			zap.Uint64("Memory Alloc (MiB)", memStats.Alloc/1024/1024),
			zap.Uint64("Memory Sys (MiB)", memStats.Sys/1024/1024),
			zap.Int("Goroutines", runtime.NumGoroutine()),
		)
	}
}

func (b *Bootstrap) RunJobs() {
	b.Job.RotateLogs()
	b.Cron.Start()
}

func (b *Bootstrap) SetInitialData() {
	ctx := context.Background()
	data := utils.AppData{
		SystemLastUpdate: time.Now().Format("02/01/2006 - 15:04"),
	}
	err := b.RedisDB.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   b.RedisDB.GetAppDataKey(),
		Value: data,
		TTL:   0,
	})
	if err != nil {
		utils.FatalResult("error caching inital data", err)
	}
}
