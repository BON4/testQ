package server

import (
	"time"

	_ "github.com/BON4/timedQ/docs"
	serviceHttp "github.com/BON4/timedQ/internal/service/delivery/http"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Log to file
func LoggerToFile(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Start time
		startTime := time.Now()

		//Process request
		c.Next()

		//End time
		endTime := time.Now()

		//Execution time
		latencyTime := endTime.Sub(startTime)

		//Request method
		reqMethod := c.Request.Method

		//Request routing
		reqUri := c.Request.RequestURI

		// status code
		statusCode := c.Writer.Status()

		// request IP
		clientIP := c.ClientIP()

		//Log format
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

func (s *Server) MapHandlers() error {
	s.g.Use(LoggerToFile(s.logger), gin.Recovery())

	v1 := s.g.Group("/v1")

	srvHand := serviceHttp.NewServiceHandler(s.wM, s.logger.WithField("service", "service-name"))

	serviceHttp.NewServiceRoutes(v1, srvHand)

	//Swagger
	s.g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return nil
}
