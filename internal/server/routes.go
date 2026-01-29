package server

import (
	"leet-repeat-api/internal/feature/problem"
	"net/http"
	"os"

	"github.com/PeterTakahashi/gin-openapi/openapiui"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
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

	r.GET("/hello", Hello)

	problemRepo := problem.NewProblemRepository(s.db.DB())
	problemHandler := problem.NewHandler(problemRepo)

	problemGroup := r.Group("/problems")
	{
		problemGroup.POST("", problemHandler.Create)
		problemGroup.GET("", problemHandler.GetAll)
		problemGroup.GET("/:id", problemHandler.GetByID)
		problemGroup.PUT("/:id", problemHandler.Update)
		problemGroup.DELETE("/:id", problemHandler.Delete)
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

// Example handler
// @Summary Hello world
// @Description returns hello
// @Produce json
// @Success 200 {string} string "hello"
// @Router /hello [get]
func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, "hello")
}
