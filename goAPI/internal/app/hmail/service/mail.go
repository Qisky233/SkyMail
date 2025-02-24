package service

import (
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"skyAPI/internal/pkg/utils"
)

type MailService struct {
	hmail *ole.IDispatch
}

func NewMailService(adminUser, adminPass string) *MailService {
	return &MailService{
		hmail: utils.InitHMailServer(adminUser, adminPass),
	}
}

func (s *MailService) GetVersion() string {
	version, err := oleutil.GetProperty(s.hmail, "Version")
	if err != nil {
		log.Fatal("获取版本失败:", err)
	}
	return version.ToString()
}

// 添加其他方法（如 CreateDomain）