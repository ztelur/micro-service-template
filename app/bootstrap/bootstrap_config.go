package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.AutomaticEnv()
	initBootstrapConfig()
	//读取yaml文件
	//v := viper.New()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}
	if err := subParse("http", &HttpConfig); err != nil {
		log.Fatal("Fail to parse Http config", err)
	}
	if err := subParse("mysql", &MySQlConfig); err != nil {
		log.Fatal("Fail to parse mysql config", err)
	}

}
func initBootstrapConfig() {
	//设置读取的配置文件
	viper.SetConfigName("bootstrap")
	//添加读取的配置文件路径
	viper.AddConfigPath("./")
	//windows环境下为%GOPATH，linux环境下为$GOPATH
	viper.AddConfigPath("$GOPATH/src/")
	//设置配置文件类型
	viper.SetConfigType("yaml")
}

func subParse(key string, value interface{}) error {
	log.Printf("配置文件的前缀为：%v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
