package http

import (
	"net/http"

	"github.com/BON4/timedQ/internal/manager"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type serviceHandler struct {
	logger      *logrus.Entry
	workManager *manager.WorkerManager
}

func (s *serviceHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		k := c.Param("key")
		c.JSON(http.StatusNotImplemented, gin.H{
			"val": s.workManager.Get(k),
		})
	}
}

func NewServiceHandler(wM *manager.WorkerManager, logger *logrus.Entry) *serviceHandler {
	return &serviceHandler{
		logger:      logger,
		workManager: wM,
	}
}
