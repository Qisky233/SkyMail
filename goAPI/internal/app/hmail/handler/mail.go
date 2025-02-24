package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"skyAPI/internal/app/hmail/service"
	"skyAPI/internal/config"
)

type MailHandler struct {
	service *service.MailService
}

func NewMailHandler(cfg *config.Config) *MailHandler {
	return &MailHandler{
		service: service.NewMailService(cfg.HMail.AdminUser, cfg.HMail.AdminPass),
	}
}

func (h *MailHandler) GetVersion(c *gin.Context) {
	version := h.service.GetVersion()
	c.JSON(http.StatusOK, gin.H{"version": version})
}

// 示例：创建域名的处理函数
func (h *MailHandler) CreateDomain(c *gin.Context) {
	// 处理业务逻辑...
}