package utils

import (
	"encoding/json"
	"gozinx/ziface"
	"io/ioutil"
)

type globalObj struct {
	IcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string // 当前服务器名称

	Version        string
	MaxConn        int    // 运行最大链接数
	MaxPackageSize uint32 // 当前zinx框架数据包的最大值
}

// 全局
var GlobalObject *globalObj

// 加载用户配置文件
func (g *globalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

func init() {
	// 如果配置文件没有加载，默认值
	GlobalObject = &globalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		Host:           "0.0.0.0",
		TcpPort:        8999,
		MaxConn:        10000,
		MaxPackageSize: 4096,
	}

	GlobalObject.Reload()
}
