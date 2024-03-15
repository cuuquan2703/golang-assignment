package logger

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	FileLogger *zap.Logger
	CmdLogger  *zap.Logger
}

var FLogger *zap.Logger
var CLogger *zap.Logger

func initCMD() {
	var err error
	CLogger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger")
	}
	defer CLogger.Sync()
}

func initFile() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(logFile),
		zapcore.InfoLevel,
	)

	FLogger = zap.New(core, zap.AddCaller())

	defer FLogger.Sync()
}

func CreateLog() *Logger {
	initCMD()
	initFile()
	return &Logger{
		FileLogger: FLogger,
		CmdLogger:  CLogger,
	}
}

func (L Logger) Info(content string) {
	t := time.Now().Format(time.RFC3339)
	L.CmdLogger.Info(t + " " + content)
	L.FileLogger.Info(t + " " + content)
}

func (L Logger) Error(content string, err error) {
	t := time.Now().Format(time.RFC3339)
	L.CmdLogger.Error(t+" "+content, zap.Error(err))
	L.FileLogger.Error(t+" "+content, zap.Error(err))
}
