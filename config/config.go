package config

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func init() {
	// 配置文件解析
	cfgName := "./config.yaml"
	file, err := os.ReadFile(cfgName)
	if err != nil {
		log.Fatalf("open config yaml [%s] got err: %s", cfgName, err)
	}

	if err := yaml.Unmarshal(file, &defaultConfig); err != nil {
		log.Fatalf("unmarshal config file got err: %s", err)
	}
}

// 默认配置
var defaultConfig *Config

// Config 应用配置
type Config struct {
	App struct {
		Name       string `yaml:"name"`
		Namespace  string `yaml:"namespace"`
		DebugEvent bool   `yaml:"debug_event"`
	}
	Log      *LogConfig  `yaml:"log"`
	Keyboard []*Keyboard `yaml:"keyboard"`
}

// String 字符串名字
func (c *Config) String() string {
	s, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshal config got err: %s", err)
	}
	return string(s)
}

// GetConfig 获取配置
func GetConfig() *Config {
	return defaultConfig
}

// LogConfig 日志配置
type LogConfig struct {
	LogPath string `yaml:"log_path"`
	Debug   bool   `yaml:"debug"`
}

// Keyboard 按键配置
type Keyboard struct {
	Name             string   `yaml:"name"`
	ProcessName      string   `yaml:"process_name"`        // 想循环操作的PID按键
	KeyPress         string   `yaml:"key_press"`           // 循环按键内容为什么
	KeySleep         int      `yaml:"key_sleep"`           // 按键间隔
	StartStopKeyDiff bool     `yaml:"start_stop_key_diff"` // 是否启停按键不一致
	StopKey          []string `yaml:"stop_key"`            // 应用停止快捷键
	StartKey         []string `yaml:"start_key"`           // 应用启动快捷键
}

// GetLogConfig 获取日志文件路径
func GetLogConfig() *LogConfig {
	if defaultConfig.Log == nil {
		log.Fatalf("config get log path got err: defaultConfig.Log is nil")
	}

	return defaultConfig.Log
}

// GetKeyboards 按键配置
func GetKeyboards() []*Keyboard {
	if defaultConfig.Keyboard == nil {
		log.Fatalf("config get keyboard got err: defaultConfig.Keyboard is nil")
	}
	return defaultConfig.Keyboard
}
