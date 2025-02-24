package utils

import (
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func InitHMailServer(adminUser, adminPass string) *ole.IDispatch {
	// COM 初始化
// 	if err := ole.CoInitialize(0); err != nil {
// 		log.Fatal("COM 初始化失败:", err)
// 	}
	// 创建 hMailServer 对象
	unknown, err := oleutil.CreateObject("hMailServer.Application")
	if err != nil {
		log.Fatal("创建对象失败:", err)
	}

	hmail, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal("获取接口失败:", err)
	}

	// 身份验证
	if _, err := oleutil.CallMethod(hmail, "Authenticate", adminUser, adminPass); err != nil {
		log.Fatal("认证失败:", err)
	}

	return hmail
}