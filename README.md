# ◀️ GinFier

> 🚧 **Note**: This project is currently under development. While it's usable, some features are still being implemented and the API might change.

A lightweight, straightforward API server for managing nginx virtual hosts, written in Go. GinFier aims to simplify the process of deploying and managing web servers through a clean REST API.

## Features
- 📘 GoLang: Fast and efficient API server
- 🚀 Simple & Quick: Deploy new virtual hosts in seconds
- 🔄 REST API: Easy to integrate with your existing tools
- 📁 Auto-Config: Generates nginx configuration files
- 🔐 SSL/TLS Management via **Certbot**
- 🔄 Reload Configuration
- 🔌 Minimal: Focus on essential features for quick deployment

## Installation

### Via binaries

### Compiling from source
```bash
# Ensure NGINX and Certbot is Installed
sudo snap install --classic certbot -y
sudp apt install nginx -y

# Clone the repository
git clone https://github.com/vanvanni/ginfier.git

# Navigate to the project directory
cd ginfier

# Build the project
go build