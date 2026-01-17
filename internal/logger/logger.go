package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(env string) error {
    var config zap.Config
    
    if env == "production" {
        config = zap.NewProductionConfig()
        config.EncoderConfig.TimeKey = "timestamp"
        config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    } else {
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }
    
    config.OutputPaths = []string{"stdout", "logs/app.log"}
    config.ErrorOutputPaths = []string{"stderr", "logs/error.log"}
    
    logger, err := config.Build()
    if err != nil {
        return err
    }
    
    Log = logger
    return nil
}

func Sync() {
    if Log != nil {
        Log.Sync()
    }
}