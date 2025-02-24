package main

import (
	"fmt"
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func main() {
	// 初始化 COM 库
	err := ole.CoInitialize(0)
	if err != nil {
		log.Fatal("COM 初始化失败:", err)
	}
	defer ole.CoUninitialize()

	// 创建 hMailServer 的 Application 对象
	unknown, err := oleutil.CreateObject("hMailServer.Application")
	if err != nil {
		log.Fatal("创建 hMailServer.Application 对象失败:", err)
	}
	defer unknown.Release()

	// 获取 IDispatch 接口
	hmail, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal("获取 IDispatch 接口失败:", err)
	}
	defer hmail.Release()

	// 调用 Authenticate 方法进行身份验证
	_, err = oleutil.CallMethod(hmail, "Authenticate", "Administrator", "你的管理员密码")
	if err != nil {
		log.Fatal("身份验证失败:", err)
	}

	// 调用 Version 方法获取服务器版本
	version, err := oleutil.GetProperty(hmail, "Version")
	if err != nil {
		log.Fatal("获取版本信息失败:", err)
	}

	// 输出服务器版本
	fmt.Println("hMailServer 版本:", version.ToString())
}