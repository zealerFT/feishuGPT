package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	options *AppConfig
)

func init() {
	cfg := &AppConfig{}
	var err error
	var file []byte
	file, err = os.ReadFile("config/app.test.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Panic().Msgf("%+v", err)
	}

	fmt.Println("-------------------------------------")

	options = cfg
}

// AppConfig store all configuration options
type AppConfig struct {
	AppName         string `yaml:"app_name" env:"APP_NAME"`
	AppEnv          string `yaml:"app_environment" env:"APP_ENVIRONMENT"`
	Debug           bool   `yaml:"debug" env:"APP_DEBUG"`
	AppRole         string `yaml:"role" env:"APP_ROLE"`
	HTTPServerAddr  string `yaml:"http_server_addr" env:"HTTP_SERVER_ADDR"`
	PROXYServerAddr string `yaml:"proxy_server_addr" env:"PROXY_SERVER_ADDR"`
	HTTPMetricsAddr string `yaml:"http_metrics_addr" env:"HTTP_METRICS_ADDR"`
	GRPCServerAddr  string `yaml:"grpc_server_addr" env:"GRPC_SERVER_ADDR"`
	GRPCMetricsAddr string `yaml:"grpc_metrics_addr" env:"GRPC_METRICS_ADDR"`

	// openai
	OpenaiSetting OpenaiSetting `yaml:"openai_setting"`
	// 飞书
	LarkSetting LarkSetting `yaml:"lark_setting"`
}

type LarkSetting struct {
	AppID             string `yaml:"app_id"`
	AppSecret         string `yaml:"app_secret"`
	EncryptKey        string `yaml:"encrypt_key"`
	VerificationToken string `yaml:"verification_token"`
	ArgocdChatId      string `yaml:"argocd_chat_id"`
	ArgocdBaseUrl     string `yaml:"argocd_base_url"`
}

type OpenaiSetting struct {
	OpenaiAddr string `yaml:"openai_addr" env:"OPENAI_ADDR"`
	Token      string `yaml:"token" env:"OPENAI_TOKEN"`
}

func SetRole(role string) {
	options.AppRole = role
}

// Options return application config options
func Options() *AppConfig {
	return options
}
