package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

// 如何把配置影射成config对象

// 从toml格式的配置文件中加载
func LoadConfigFromToml(filepath string) error {
	config = NewDefaultConfig()

	// 读取toml格式的配置文件
	_, err := toml.DecodeFile(filepath, config)
	if err != nil {
		return fmt.Errorf("Load config file error, path: %s, %s", filepath, err)
	}

	// return loadGlobal()
	return nil
}

// 从环境变量中加载
func LoadConfigFromEnv() error {
	config = NewDefaultConfig()

	// 从环境变量加载配置
	err := env.Parse(config)
	if err != nil {
		return err
	}

	// return loadGlobal()
	return nil
}

// 加载DB的全局实例
func loadGlobal() (err error) {
	// 加载DB的全局实例
	db, err = config.MySQL.getDBConn()
	if err != nil {
		return err
	}

	return
}
