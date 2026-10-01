package logger

import (
	"Go/web_app/settings"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger
 
func Init(cfg *settings.LogConfig) (err error) {
	writeSyncer := getLogWriter(
		cfg.FileName, 
		cfg.Max_size, 
		cfg.Max_backups, 
		cfg.Max_age)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil{
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	lg := zap.New(core, zap.AddCaller())
	
	zap.ReplaceGlobals(lg)
	return
}

// 编码器(如何写入日志)
func getEncoder() zapcore.Encoder {
	//更改时间编码
	encoderConfig  := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	//按照 json 格式编码
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//按照 console 格式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

//指定日志将写到哪里去
// func getLogWriter() zapcore.WriteSyncer {
// 	file, _ := os.OpenFile("./test.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0744)
// 	return zapcore.AddSync(file)
// }

// 使用Lumberjack进行日志切割归档
// 添加日志切割归档功能
func getLogWriter(filename string, maxsize, maxbackups, maxage int ) zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxsize, //容量 MB
		MaxBackups: maxbackups,   //最大备份数量
		MaxAge:     maxage,  //最大备份天数
		Compress:   false, 
	}
	return zapcore.AddSync(lumberjackLogger)
}

// Ginlogger 接受gin框架默认的日志
func Ginlogger() gin.HandlerFunc {
	return func(c *gin.Context) {
			start := time.Now()
			path := c.Request.URL.Path
			query := c.Request.URL.RawQuery
			c.Next()

			cost := time.Since(start)
			zap.L().Info(path,
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
func GinRecovery(stack bool) gin.HandlerFunc {
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
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
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

