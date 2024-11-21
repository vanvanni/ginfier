package nginx

import (
	"crypto/tls"
	"fmt"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

type CertManager struct {
	manager *autocert.Manager
}

func NewCertManager(cacheDir string) (*CertManager, error) {
	if os.Getenv("LETSENCRYPT_EMAIL") == "" {
		return nil, fmt.Errorf("LETSENCRYPT_EMAIL environment variable is required")
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cert cache directory: %v", err)
	}

	manager := &autocert.Manager{
		Cache:  autocert.DirCache(cacheDir),
		Prompt: autocert.AcceptTOS,
		Email:  os.Getenv("LETSENCRYPT_EMAIL"),
	}

	return &CertManager{
		manager: manager,
	}, nil
}

func (cm *CertManager) RequestCertificate(domain string) error {
	cm.manager.HostPolicy = autocert.HostWhitelist(domain)
	_, err := cm.manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain})
	if err != nil {
		return fmt.Errorf("failed to obtain certificate: %v", err)
	}

	return nil
}
