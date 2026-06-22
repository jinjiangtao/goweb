package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"qrcode/internal/model"
)

func (svc *Service) stats(c *gin.Context) {
	st, err := svc.Store.Stats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("统计失败: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.OK(st))
}
