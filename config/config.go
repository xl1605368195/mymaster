package config

import (
	"github.com/spf13/viper"
	"jrasp-master/frame/logger"
)

// jrasp-master配置文件
type MasterConfig struct {
	Env        string `json:"env"` // 运行环境: dev,test,uat,prod
	EtcdConfig `json:"etcdConfig"`
}

type EtcdConfig struct {
	Endpoints []string `json:"endpoints"`
}

func InitConfig(logger *logger.Writer) (*MasterConfig, error) {
	// 查找配置文件
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath(".")

	// 设置默认配置
	vp.SetDefault("env", "dev")
	vp.SetDefault("etcdConfig.endpoints", []string{"localhost:2379"})
	// TODO 新增配置的默认值在这里加


	err := vp.ReadInConfig()
	if err != nil {
		logger.Warning(1000, "read config failed, use default config,err:%+v", err)
	}

	// 配置输出
	var configjson MasterConfig
	if err := vp.Unmarshal(&configjson); err != nil {
		logger.Err(2000, "vp.Unmarshal failed,err:%+v", err)
		return nil, err
	}
	return &configjson, nil
}
