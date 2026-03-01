package server

import (
	"leet-repeat-api/internal/database/repositories/progress"
	"leet-repeat-api/internal/handlers"
	"net/http"
	"os"

	"github.com/PeterTakahashi/gin-openapi/openapiui"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	specPath := "./docs/swagger.json"
	if os.Getenv("APP_ENV") == "production" {
		specPath = "../../docs/swagger.json"
	}

	r.GET("/docs/*any", openapiui.WrapHandler(openapiui.Config{
		SpecURL:      "/docs/openapi.json",
		SpecFilePath: specPath,
		Title:        "Example API",
		Theme:        "dark",
	}))

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	progressRepo := progress.NewProgressRepository(s.db.DB())
	progressHandler := handlers.NewProgressHandler(progressRepo)

	progressGroup := r.Group("/api/progress")
	{
		progressGroup.POST("/bulk-upsert", progressHandler.BulkUpsert)
		progressGroup.GET("", progressHandler.GetAll)
		progressGroup.DELETE("/clear", progressHandler.Clear)
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
