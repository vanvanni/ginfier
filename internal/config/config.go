package config

type (
	ProxyConfig struct {
		Domain      string `json:"domain"`
		Port        int    `json:"port"`
		Destination string `json:"destination"`
		EnableSSL   bool   `json:"enable_ssl"`
	}

	APIResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)
