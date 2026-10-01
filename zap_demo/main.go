package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"strings"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func initLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel) // "zapcore.DebugLevel" 哪种级别的日志将被写入

	//将调用函数信息记录到日志中
	logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

// 编码器(如何写入日志)
func getEncoder() zapcore.Encoder {
	//更改时间编码
	encoderCinfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,

		EncodeTime: zapcore.ISO8601TimeEncoder,

		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	//按照 json 格式编码
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//按照 console 格式
	return zapcore.NewConsoleEncoder(encoderCinfig)
}

//指定日志将写到哪里去
// func getLogWriter() zapcore.WriteSyncer {
// 	file, _ := os.OpenFile("./test.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0744)
// 	return zapcore.AddSync(file)
// }

// 使用Lumberjack进行日志切割归档
// 添加日志切割归档功能
func getLogWriter() zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    200, //容量 MB
		MaxBackups: 5,   //最大备份数量
		MaxAge:     30,  //最大备份天数
		Compress:   false, 
	}
	return zapcore.AddSync(lumberjackLogger)
}

func simpleHttpGet(url string) {
	response, err := http.Get(url)
	if err != nil {
		sugarLogger.Error(
			"Error fetching url...",
			zap.String("url", url),
			zap.Error(err),
		)
	} else {
		sugarLogger.Info(
			"success..",
			zap.String("statusCode", response.Status),
			zap.String("url", url),
		)
		response.Body.Close()

	}
}

// Ginlogger 接受gin框架默认的日志
func Ginlogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
			start := time.Now()
			path := c.Request.URL.Path
			query := c.Request.URL.RawQuery
			c.Next()

			cost := time.Since(start)
			logger.Info(path,
				zap.Int("Status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
			)
		}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// func main() {
// 	initLogger()
// 	defer logger.Sync()
// 	simpleHttpGet("www.sogo.com")
// 	simpleHttpGet("http://www.sogo.com")
// }

func main() {
	initLogger()
	//r := gin.Default()
	r := gin.New()
	r.Use(Ginlogger(logger), GinRecovery(logger, true))
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello liwenzhou.com!")})
	r.Run()
}
 