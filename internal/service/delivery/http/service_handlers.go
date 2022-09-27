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

type serviceSetResponse struct {
	EncodeURL string `json:"encode_url"`
}

type serviceGetResponse struct {
	DecodeURL string `json:"decode_url"`
}

type serviceHandler struct {
	logger      *logrus.Entry
	workManager *manager.WorkerManager
}

// @Summary      Get redirect
// @Description  by known key, user can get an url
// @Tags         general
// @Produce      json
// @Param        key  path      string  true  "decoded full url"
// @Success      200  {object}  serviceGetResponse
// @Router       /{key} [Get]
func (s *serviceHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		k := c.Param("key")
		c.JSON(http.StatusOK, gin.H{
			"val": s.workManager.Get(k),
		})
	}
}

// @Summary      Set redirect
// @Description  sets key-value, where user is providing value, and gets key
// @Tags         general
// @Accept       json
// @Produce      json
// @Param        input   body      serviceSetRequest  true  "encoded short url"
// @Success      200     {object}  serviceSetResponse
// @Failure      400     {object}  error
// @Router       / [post]
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

		c.JSON(http.StatusOK, serviceSetResponse{
			EncodeURL: link,
		})
	}
}

func NewServiceHandler(wM *manager.WorkerManager, logger *logrus.Entry) *serviceHandler {
	return &serviceHandler{
		logger:      logger,
		workManager: wM,
	}
}
