package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"text/template"

	"embed"

	"github.com/vanvanni/nginfier/internal/config"
	"github.com/vanvanni/nginfier/internal/ssl"
)

//go:embed templates/nginx.conf.tmpl
var embeddedFiles embed.FS

func validateAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != os.Getenv("API_SECRET") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func getNginxConfigPath() string {
	switch runtime.GOOS {
	case "linux":
		if _, err := os.Stat("/etc/nginx/conf.d"); err == nil {
			return "/etc/nginx/conf.d"
		}
		return "/etc/nginx/sites-enabled"
	default:
		return "/etc/nginx/conf.d"
	}
}

func createProxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config config.ProxyConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var missingFields []string
	if config.Domain == "" {
		missingFields = append(missingFields, "domain")
	}
	if config.Destination == "" {
		missingFields = append(missingFields, "destination")
	}

	if len(missingFields) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.ValidationError{MissingFields: missingFields})
		return
	}

	// Request SSL certificate if enabled
	if config.EnableSSL {
		certManager, err := ssl.NewCertManager("/etc/letsencrypt/live")
		if err != nil {
			log.Printf("Failed to create cert manager: %v", err)
			http.Error(w, "Failed to initialize SSL manager", http.StatusInternalServerError)
			return
		}

		if err := certManager.RequestCertificate(config.Domain); err != nil {
			log.Printf("Failed to request SSL certificate: %v", err)
			http.Error(w, "Failed to request SSL certificate", http.StatusInternalServerError)
			return
		}
	}

	// Load template from file
	tmpl, err := loadTemplate()
	if err != nil {
		log.Printf("Failed to load template: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Create config file path
	configPath := fmt.Sprintf("%s/%s.conf", getNginxConfigPath(), config.Domain)
	f, err := os.Create(configPath)
	if err != nil {
		log.Printf("Failed to create config file: %v", err)
		http.Error(w, "Failed to create config file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Execute template
	if err := tmpl.Execute(f, config); err != nil {
		http.Error(w, "Failed to write config", http.StatusInternalServerError)
		return
	}

	// Reload NGINX
	// Note: This requires sudo privileges or proper system configuration
	// TODO: Implement NGINX reload

	response := config.APIResponse{
		Success: true,
		Message: "Proxy configuration created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func loadTemplate() (*template.Template, error) {
	tmplContent, err := embeddedFiles.ReadFile("templates/nginx.conf.tmpl")
	if err != nil {
		return nil, err
	}

	return template.New("nginx.conf").Parse(string(tmplContent))
}

func main() {
	required := []string{"API_SECRET"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			log.Fatalf("%s environment variable is required", env)
		}
	}

	http.HandleFunc("/api/proxy", validateAPIKey(createProxyHandler))

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
