package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xin-24/go/mxshop-api/user-web/mxshop-api/user-web/global"
)

type MysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	ServiceName string      `mapstructure:"name"`
	MysqlInfo   MysqlConfig `mapstructure:"mysql"`
}

// 如何将线下和线上的东西隔离
func GetEnvInfo(env string) bool { //获取计算机的变量
	viper.AutomaticEnv()
	return viper.GetBool(env)
	//设置完需要重启
}
func InitConfig() {
	debug := (GetEnvInfo("mxshop_debug"))
	configFilePreFix := "config"
	ConfigFileName := fmt.Sprintf("mxshop-api/user-web/%s_pro.yaml", configFilePreFix) //线上
	if debug {
		ConfigFileName = fmt.Sprintf("mxshop-api/user-web/%s_debug.yaml", configFilePreFix) //线下
	}
	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(ConfigFileName) //为当前文件位置
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用----全局变量

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}

	zap.S().Infof("配置信息：%v", global.ServerConfig)

	//viper的功能监测动态变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件产生变化：%s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
		zap.S().Infof("配置信息：%v", global.ServerConfig)
	})

}
