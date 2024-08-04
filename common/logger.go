package common

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConf struct {
	Level      string `json:"level" toml:"level"`             // Level 最低日志等级，DEBUG<INFO<WARN<ERROR<FATAL
	Filename   string `json:"filename" toml:"filename"`       // Filename 日志文件位置
	MaxSize    int    `json:"max_size" toml:"max_size"`       // MaxSize 日志文件的最大大小(MB为单位)，默认为 100MB
	MaxAge     int    `json:"max_age" toml:"max_age"`         // MaxAge 保留旧日志文件的最大天数
	MaxBackups int    `json:"max_backups" toml:"max_backups"` // MaxBackups 是要保留的旧日志文件的最大数量
	Compress   bool   `json:"compress" toml:"compress"`       // Compress 是否压缩归档日志
}

func getEncoder() zapcore.Encoder {
	// 获取一个指定的的EncoderConfig，进行自定义
	encodeConf := zap.NewProductionEncoderConfig()

	// 设置时间格式
	encodeConf.EncodeTime = zapcore.RFC3339TimeEncoder

	// 将 Level 序列化为全大写字符串
	encodeConf.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewJSONEncoder(encodeConf)
}

// 负责日志写入的位置
func getLogWriter(c *LogConf) zapcore.WriteSyncer {
	l := &lumberjack.Logger{
		Filename:   c.Filename,   // 文件位置
		MaxSize:    c.MaxSize,    // 日志文件的最大大小（MB）
		MaxAge:     c.MaxAge,     // 保留旧文件的最大天数
		MaxBackups: c.MaxBackups, // 保留旧文件的最大个数
		Compress:   c.Compress,   // 是否压缩/归档旧文件
	}
	// AddSync 将 io.Writer 转换为 WriteSyncer
	return zapcore.AddSync(l)
}

// InitLogger 初始化Logger
func InitLogger(lc *LogConf) error {
	// 获取日志写入位置
	writeSyncer := getLogWriter(lc)
	// 获取日志编码格式
	encoder := getEncoder()

	// 获取日志最低等级
	level := new(zapcore.Level)
	err := level.UnmarshalText([]byte(lc.Level))
	if err != nil {
		return err
	}

	// 创建一个将日志写入 WriteSyncer 的核心
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core, zap.AddCaller())

	// 替换 zap 包中全局的 logger 实例，后续在其他包中只需使用 zap.L() 调用即可
	zap.ReplaceGlobals(logger)
	return nil
}
