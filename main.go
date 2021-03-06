package main

import (
	"fmt"

	"go-copyright-p1/configs"
	"go-copyright-p1/routes"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var EchoObj *echo.Echo //echo框架对象全局定义

func main() {

	fmt.Printf("get config %v ,%v\n", configs.Config.Common.Port, configs.Config.Db.Connstr)
	EchoObj = echo.New()             //创建echo对象
	EchoObj.Use(middleware.Logger()) //安装日志中间件
	EchoObj.Use(middleware.Recover())
	EchoObj.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	EchoObj.GET("/ping", routes.PingHandler)                        //路由测试函数
	EchoObj.Logger.Fatal(EchoObj.Start(configs.Config.Common.Port)) //启动服务
}
