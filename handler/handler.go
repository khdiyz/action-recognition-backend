package handler

import (
	"action-detector-backend/config"
	"action-detector-backend/docs"
	"action-detector-backend/pkg/logger"
	"action-detector-backend/pkg/postgres"
	"action-detector-backend/usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	usecase *usecase.Usecase
	cfg     *config.Config
	logger  *logger.Logger
	DB      *postgres.Postgres
}

func NewHandler(usecase *usecase.Usecase, cfg *config.Config, logger *logger.Logger, db *postgres.Postgres) *Handler {
	return &Handler{
		usecase: usecase,
		cfg:     cfg,
		logger:  logger,
		DB:      db,
	}
}

func (h *Handler) InitRoutes(cfg *config.Config) *gin.Engine {
	if cfg.Environment == config.EnvironmentProd {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.HandleMethodNotAllowed = true
	router.Use(corsMiddleware())

	// Setup Swagger documentation
	h.setupSwagger(router)

	api := router.Group("/api")

	api.POST("/files", h.uploadVideoFile)
	api.POST("/predict", h.predictAction)
	api.GET("/actions", h.getActions)
	api.DELETE("/actions", h.deleteActions)

	return router
}

func (h *Handler) setupSwagger(router *gin.Engine) {
	router.GET("/docs/*any", func(ctx *gin.Context) {
		// Set CORS headers for Swagger UI
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Update Swagger info
		docs.SwaggerInfo.Host = ctx.Request.Host
		if ctx.Request.TLS != nil {
			docs.SwaggerInfo.Schemes = []string{"https"}
		}

		// Handle OPTIONS request
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ginSwagger.WrapHandler(swaggerFiles.Handler)(ctx)
	})
}
