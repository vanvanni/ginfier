#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'
print_status() {
    echo -e "${GREEN}==>${NC} $1"
}

print_error() {
    echo -e "${RED}Error:${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

if [ "$EUID" -ne 0 ]; then 
    print_error "Please run as root"
    echo "Try: sudo bash $0"
    exit 1
fi

if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$NAME
    VERSION=$VERSION_ID
else
    print_error "Cannot detect OS"
    exit 1
fi

install_dependencies() {
    print_status "Installing dependencies..."
    case $OS in
        "Ubuntu"|"Debian GNU/Linux")
            apt-get update
            apt-get install -y nginx curl wget
            ;;
        "CentOS Linux"|"Red Hat Enterprise Linux")
            yum install -y epel-release
            yum install -y nginx curl wget
            ;;
        *)
            print_error "Unsupported OS: $OS"
            exit 1
            ;;
    esac
}

# Download latest release
download_app() {
    print_status "Downloading nginfier..."
    LATEST_RELEASE_URL="https://api.github.com/repos/YOUR_REPO/nginfier/releases/latest"
    DOWNLOAD_URL=$(curl -s $LATEST_RELEASE_URL | grep "browser_download_url.*386" | cut -d '"' -f 4)
    
    if [ -z "$DOWNLOAD_URL" ]; then
        print_error "Could not find download URL"
        exit 1
    }
    
    wget -O /usr/local/bin/nginfier $DOWNLOAD_URL
    chmod +x /usr/local/bin/nginfier
}

setup_service() {
    print_status "Setting up systemd service..."
    
    # Create required directories
    mkdir -p /var/log/nginfier
    mkdir -p /etc/nginfier
    
    cp examples/nginfier.service /etc/systemd/system/nginfier.service
    
    chown -R root:root /var/log/nginfier
    chmod 755 /var/log/nginfier
    chown -R root:root /etc/nginfier
    chmod 755 /etc/nginfier

    systemctl daemon-reload
    systemctl enable nginfier
    systemctl start nginfier
}

# Configure Nginx
setup_nginx() {
    print_status "Configuring Nginx..."
    
    if [ -f /etc/nginx/nginx.conf ]; then
        mv /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup
    fi

    cat > /etc/nginx/nginx.conf << 'EOL'
user www-data;
worker_processes auto;
pid /run/nginx.pid;

events {
    worker_connections 768;
}

http {
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;

    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;

    include /etc/nginx/conf.d/*.conf;
}
EOL

    systemctl enable nginx
    systemctl restart nginx
}

main() {
    print_status "Starting installation..."
    
    install_dependencies
    
    setup_nginx

    download_app
    setup_service
    
    # Final status check
    if systemctl is-active --quiet nginfier && systemctl is-active --quiet nginx; then
        print_status "Installation completed successfully!"
        echo -e "\nNginx service status:"
        systemctl status nginx --no-pager
        echo -e "\nNginfier service status:"
        systemctl status nginfier --no-pager
    else
        print_error "Installation completed with errors. Please check the logs."
        exit 1
    fi
}

main 