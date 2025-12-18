package logger

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/pur1fying/GO_BAAS/internal/global_info"
	"github.com/pur1fying/GO_BAAS/internal/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _logger *zap.Logger

const _highlightLogLineLen int = 80

var _BAASLoggerHighlightLine = strings.Repeat("-", _highlightLogLineLen)

var _BAASLoggerDivider = strings.Repeat("=", _highlightLogLineLen)

func InitGlobalLogger() {

	currentTimeStr := utils.CurrentTimeString()
	outputFilename := path.Join(global_info.GO_BAAS_OUTPUT_DIR, currentTimeStr, "global_log.txt")

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   outputFilename,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    "logger",
		MessageKey: "msg",

		EncodeLevel:    zapcore.CapitalLevelEncoder, // INFO, ERROR
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,

		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},

		LineEnding: zapcore.DefaultLineEnding,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	level := zapcore.DebugLevel

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, level),                 // 写入文件
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level), // 控制台输出
	)

	_logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

}

func Flush() error {
	return _logger.Sync()
}

func BAASDebug(_strings ...string) {
	_logger.Debug(strings.Join(_strings, " "))
}

func BAASInfo(_strings ...string) {
	_logger.Info(strings.Join(_strings, " "))
}

func BAASWarn(_strings ...string) {
	_logger.Warn(strings.Join(_strings, " "))
}

func BAASError(_strings ...string) {
	_logger.Error(strings.Join(_strings, " "))
}

func BAASCritical(_strings ...string) {
	_logger.Panic(strings.Join(_strings, " "))
}

func HighLight(message string) {
	BAASInfo(_BAASLoggerHighlightLine)
	BAASInfo(generateHighlightLogMessage(message))
	BAASInfo(_BAASLoggerHighlightLine)
}

func SubTitle(message string) {
	BAASInfo("<<< ", message, " >>>")
}

func Line() {
	BAASInfo(_BAASLoggerDivider)
}

func generateHighlightLogMessage(message string) string {
	var messageLen int = len(message)
	var leftSpaceLen = (_highlightLogLineLen - 2 - messageLen) / 2
	var rightSpaceLen = _highlightLogLineLen - messageLen - leftSpaceLen - 2
	return "|" + strings.Repeat(" ", leftSpaceLen) + message + strings.Repeat(" ", rightSpaceLen) + "|"
}
