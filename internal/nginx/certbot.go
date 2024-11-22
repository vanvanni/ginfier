package nginx

import (
	"fmt"
	"os"
	"os/exec"
)

func RequestCertificate(domain string) error {
	email := os.Getenv("LETSENCRYPT_EMAIL")
	cmd := exec.Command("sudo", "certbot", "certonly",
		"--nginx",
		"--non-interactive",
		"--agree-tos",
		"-d", domain,
		"-m", email)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("certbot error: %v\nOutput: %s", err, output)
	}

	return nil
}

func RenewCertificates() error {
	cmd := exec.Command("sudo", "certbot", "renew")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("renewal error: %v\nOutput: %s", err, output)
	}
	return nil
}
