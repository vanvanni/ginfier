package main

import (
	"fmt"

	"github.com/vanvanni/ginfier/internal/logger"
	"github.com/vanvanni/ginfier/internal/nginx"
	"github.com/vanvanni/ginfier/internal/templates"
)

func main() {
	logger.Info("Starting: GinFier")
	fmt.Println(templates.ReverseHost())
	fmt.Println(nginx.GetPath())
}
