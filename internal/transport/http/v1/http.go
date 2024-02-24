package v1

import (
	"seabattle/internal/service"
	"github.com/evrone/go-seabattle-template/pkg/logger"
	"github.com/gin-gonic/gin"
)

type transport struct {
	s service.Service
	l logger.Interface
}

func newTransport(handler *gin.RouterGroup, t service.Service, l logger.Interface) {
	r := &transport{t, l}
	r = r
	//h := handler.Group("/translation")
	//{
	//	h.GET("/history", r.history)
	//	h.POST("/do-translate", r.doTranslate)
	//}
}
func NewRouter(handler *gin.Engine, l logger.Interface, t service.Service) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	//// Swagger
	//swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	//handler.GET("/swagger/*any", swaggerHandler)

	//// K8s probe
	//handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	//
	//// Prometheus metrics
	//handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newTransport(h, t, l)
	}
}
