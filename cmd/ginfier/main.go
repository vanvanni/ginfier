package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vanvanni/ginfier/internal/config"
	"github.com/vanvanni/ginfier/internal/logger"
	"github.com/vanvanni/ginfier/internal/nginx"
	"github.com/vanvanni/ginfier/internal/templates"
)

func validateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != os.Getenv("API_SECRET") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		c.Next()
	}
}
func createHandler(c *gin.Context) {
	var cfg config.ProxyConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(400, gin.H{"code": "INVALID_BODY", "error": "Invalid request body"})
		return
	}

	// Validate required fields
	var missingFields []string
	if cfg.Domain == "" {
		missingFields = append(missingFields, "domain")
	}
	if cfg.Destination == "" {
		missingFields = append(missingFields, "destination")
	}

	if len(missingFields) > 0 {
		c.JSON(400, gin.H{"code": "INVALID_BODY", "missing_fields": missingFields})
		return
	}

	if cfg.EnableSSL {
		certManager, err := nginx.NewCertManager("/etc/letsencrypt/live")
		if err != nil {
			log.Printf("Failed to create cert manager: %v", err)
			c.JSON(500, gin.H{"code": "SERVER_ERROR", "error": "Failed to initialize SSL manager"})
			return
		}

		if err := certManager.RequestCertificate(cfg.Domain); err != nil {
			log.Printf("Failed to request SSL certificate: %v", err)
			c.JSON(500, gin.H{"code": "SERVER_ERROR", "error": "Failed to request SSL certificate"})
			return
		}
	}

	nginxPath, err := nginx.GetPath()
	if err != nil {
		logger.Fatal("Could not load NGINX path")
	}

	configPath := fmt.Sprintf("%s/%s.conf", nginxPath, cfg.Domain)
	f, err := os.Create(configPath)
	if err != nil {
		log.Printf("Failed to create config file: %v", err)
		c.JSON(500, gin.H{"code": "SERVER_ERROR", "error": "Failed to create config file"})
		return
	}
	defer f.Close()

	tpl, err := templates.ReverseHost()
	if err != nil {
		logger.Fatal("Could not load NGINX path")
	}

	if err := tpl.Execute(f, cfg); err != nil {
		c.JSON(500, gin.H{"error": "Failed to write config"})
		return
	}

	// Reload NGINX
	// Note: This requires sudo privileges or proper system configuration
	// TODO: Implement NGINX reload

	response := config.APIResponse{
		Code:    "OK",
		Message: "Proxy configuration created successfully",
	}

	c.JSON(200, response)
}

func main() {
	logger.Info("Starting: GinFier")

	required := []string{"API_SECRET"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			logger.Fatal("Missing: API_SECRET")
		}
	}

	fmt.Println(templates.ReverseHost())
	fmt.Println(nginx.GetPath())

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    "OK",
			"message": "pong",
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	authorized := r.Group("/api")
	authorized.Use(validateAPIKey())
	{
		authorized.POST("/create", createHandler)
	}

	r.Run(":3000")
}
