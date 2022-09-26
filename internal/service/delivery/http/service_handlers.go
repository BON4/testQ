package http

import (
	"net/http"

	"github.com/BON4/timedQ/internal/manager"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type serviceSetRequest struct {
	Redirect string `json:"redirect" binding:"required"`
}

type serviceHandler struct {
	logger      *logrus.Entry
	workManager *manager.WorkerManager
}

func (s *serviceHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		k := c.Param("key")
		c.JSON(http.StatusOK, gin.H{
			"val": s.workManager.Get(k),
		})
	}
}

func (s *serviceHandler) Set() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &serviceSetRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// TODO: Change to shortten provided link
		link := uuid.New().String()

		s.workManager.Set(link, req.Redirect)

		c.JSON(http.StatusOK, gin.H{
			"link": link,
		})
	}
}

func NewServiceHandler(wM *manager.WorkerManager, logger *logrus.Entry) *serviceHandler {
	return &serviceHandler{
		logger:      logger,
		workManager: wM,
	}
}
