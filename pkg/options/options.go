package options

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Version string = "v1.0.0"

type LogOptions struct {
	Level    string `yaml:"level"`
	Path     string `yaml:"path"`
	MaxSize  int    `yaml:"maxsize"`
	MaxAge   int    `yaml:"maxage"`
	Compress bool   `yaml:"compress"`
}

type CoreOptions struct {
	Log LogOptions `yaml:"log"`
}

type TaskOptions struct {
	Interval   int    `yaml:"interval"`
	OutputPath string `yaml:"outputpath"`
}

type RedisOptions struct {
	Addr        string `yaml:"addr"`
	Password    string `yaml:"password"`
	DB          int    `yaml:"db"`
	DialTimeout int    `yaml:"dialtimeout"`
	Prefix      string `yaml:"prefix"`
}

type Options struct {
	Core  CoreOptions  `yaml:"core"`
	Task  TaskOptions  `yaml:"task"`
	Redis RedisOptions `yaml:"redis"`
}

func NewOptions() (opts Options) {
	viper.SetDefault("core.interval", 60)

	optsSource := viper.AllSettings()
	err := createOptions(optsSource, &opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "create options failed:", err)
		os.Exit(1)
	}
	return
}
