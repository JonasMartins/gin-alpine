// Package bootstrap
package bootstrap

import (
	"gin-alpine/src/internal/configs"
	"gin-alpine/src/pkg/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Bootstrap struct {
	Logger        *zap.Logger
	Config        *configs.Config
	RollingLogger *lumberjack.Logger
}

func MustGetBootstrapInstance() *Bootstrap {
	cfg, err := configs.LoadConfig()
	if err != nil {
		utils.FatalResult("failed to load config", err)
	}
	b := Bootstrap{
		Config: cfg,
	}
	b.setUpLogger()
	return &b
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
