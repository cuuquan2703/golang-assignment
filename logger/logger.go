package logger

import (
	"log"
	"os"

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
	CLogger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger")
	}
	defer CLogger.Sync()
}

func initFile() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
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
	L.CmdLogger.Info(content)
	L.FileLogger.Info(content)
}

func (L Logger) Error(content string, err error) {
	L.CmdLogger.Error(content, zap.Error(err))
	L.FileLogger.Error(content, zap.Error(err))
}
