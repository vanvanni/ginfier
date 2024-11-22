package config

type (
	ProxyConfig struct {
		Domain      string `json:"domain"`
		Destination string `json:"destination"`
		EnableSSL   bool   `json:"ssl"`
	}

	APIResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)
