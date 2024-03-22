package initialize

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var logOnce sync.Once
var Logger *zap.Logger

func InitLogger() *zap.Logger {
	logOnce.Do(func() {
		//	rawJSON := []byte(`{
		// "level": "debug",
		// "encoding": "json",
		// "outputPaths": ["stdout", "./logs.txt"],
		// "errorOutputPaths": ["stderr"],
		// "encoderConfig": {
		//   "messageKey": "message",
		//   "levelKey": "level",
		//   "levelEncoder": "lowercase"
		// }
		//}`)
		//
		//	var cfg zap.Config
		//	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		//	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		//	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		//		panic(err)
		//	}
		//Logger, _ = zap.NewProduction(zap.AddCaller())

		writeSyncer := getLogWriter()
		encoder := getEncoder()
		core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

		Logger = zap.New(core)
	})

	return Logger
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./connectly_bot.log",
		MaxSize:    20,
		MaxAge:     7,
		MaxBackups: 10,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
