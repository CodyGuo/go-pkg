package logger

import (
	"fmt"
	"time"

	"github.com/CodyGuo/go-pkg/fs"
	"github.com/rs/zerolog"
)

type Config struct {
	Level               string `json:"level" mapstructure:"level"`
	TimeFormat          string `json:"time_format" mapstructure:"time_format"`
	FilePath            string `json:"filepath" mapstructure:"filepath"`
	AccessFilePath      string `json:"access_filepath" mapstructure:"access_filepath"`
	MaxSize             int    `json:"max_size" mapstructure:"max_size"`
	MaxAge              int    `json:"max_age" mapstructure:"max_age"`
	MaxBackups          int    `json:"max_backups" mapstructure:"max_backups"`
	Compress            bool   `json:"compress" mapstructure:"compress"`
	UTCTime             bool   `json:"utc_time" mapstructure:"utc_time"`
	EnableFile          bool   `json:"enable_file" mapstructure:"enable_file"`
	EnableConsole       bool   `json:"enable_console" mapstructure:"enable_console"`
	EnableAccessFile    bool   `json:"enable_access_file" mapstructure:"enable_access_file"`
	EnableAccessConsole bool   `json:"enable_access_console" mapstructure:"enable_access_console"`
}

func init() {
	conf := Config{
		Level:               "info",
		TimeFormat:          TimeFormat,
		UTCTime:             false,
		EnableFile:          false,
		EnableConsole:       true,
		EnableAccessFile:    false,
		EnableAccessConsole: false,
	}
	conf.Init()
}

func (c *Config) Init() error {
	err := c.validate()
	if err != nil {
		return err
	}

	loggerConf := *c
	_logger, err = New(loggerConf)
	if err != nil {
		return err
	}

	aLoggerConf := *c
	aLoggerConf.FilePath = aLoggerConf.AccessFilePath
	aLoggerConf.EnableFile = aLoggerConf.EnableAccessFile
	aLoggerConf.EnableConsole = aLoggerConf.EnableAccessConsole
	_accessLogger, err = New(aLoggerConf)
	if err != nil {
		return err
	}

	return nil
}
func (c *Config) validate() error {
	if c.EnableFile && !fs.IsFilePathValid(c.FilePath) {
		return fmt.Errorf("log filepath (%q) invalid", c.FilePath)
	}
	if c.EnableAccessFile && !fs.IsFilePathValid(c.AccessFilePath) {
		return fmt.Errorf("access log filepath (%q) invlaid", c.AccessFilePath)
	}
	return nil
}

func SetLogTime(timeFormat string, utc bool) {
	zerolog.TimeFieldFormat = timeFormat
	if utc {
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().UTC()
		}
	} else {
		zerolog.TimestampFunc = time.Now
	}
}
