package nginx

import (
	"errors"
	"os"
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
