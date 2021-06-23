package main

import (
	"github.com/spf13/viper"
)

// jrasp-master配置文件
type MasterConfig struct {
	Env        string `json:"env"` // 运行环境: dev,test,uat,prod
	EtcdConfig `json:"etcdConfig"`
}

type EtcdConfig struct {
	Endpoints []string `json:"endpoints"`
}

func (m *Master) InitConfig() error {
	// 查找配置文件
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath(".")

	// 设置默认配置
	vp.SetDefault("env", "dev")
	vp.SetDefault("etcdConfig.endpoints", []string{"localhost:2379"})

	err := vp.ReadInConfig()
	if err != nil {
		m.logger.Warning(LOGGER_INIT_ERROR, "read config failed, use default config,err:%+v", err)
	}

	// 配置输出
	var configjson MasterConfig
	if err := vp.Unmarshal(&configjson); err != nil {
		m.logger.Err(LOGGER_INIT_ERROR, "vp.Unmarshal failed,err:%+v", err)
		return err
	}

	m.config = configjson
	return nil
}
