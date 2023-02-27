package handler

import (
	"net/http"

	"actlabs-auth/entity"

	"github.com/gin-gonic/gin"
)

type RedisHandler struct {
	RedisService entity.RedisService
}

func NewRedisHandler(r *gin.Engine, redisService entity.RedisService) {
	handler := &RedisHandler{
		RedisService: redisService,
	}

	r.DELETE("/cache", handler.DeleteServerCache)
}

func (r *RedisHandler) DeleteServerCache(c *gin.Context) {
	if err := r.RedisService.ResetServerCache(); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}
