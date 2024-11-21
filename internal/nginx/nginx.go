package nginx

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func GetPath() (string, error) {
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "/etc/nginx/sites-enabled", nil
	}

	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return "/etc/nginx/conf.d", nil
	}
	return "", errors.New("unsupported os release")
}

func Restart() error {
	// Execute systemctl restart nginx
	cmd := exec.Command("systemctl", "restart", "nginx")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to restart Nginx: %s - %w", output, err)
	}

	return nil
}

func Reload() error {
	// Execute systemctl reload nginx
	cmd := exec.Command("systemctl", "reload", "nginx")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload Nginx: %s - %w", output, err)
	}

	return nil
}
