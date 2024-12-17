#!/bin/bash

# Ensure the script is run by root
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root"
    exit 1
fi

# Define variables
SERVICE_NAME="idv_host"
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"
APP_PATH="/home/user/idv_host"
APP_EXEC="$APP_PATH/idv_host"

# Create the service file
echo "[Unit]
Description=IDV Host Service
After=network.target

[Service]
ExecStart=$APP_EXEC
WorkingDirectory=$APP_PATH
Restart=always
User=$(whoami)

[Install]
WantedBy=multi-user.target" > $SERVICE_FILE

# Reload systemd to recognize the new service
systemctl daemon-reload

# Enable the service to start on boot
systemctl enable $SERVICE_NAME

# Start the service
systemctl start $SERVICE_NAME

echo "Service $SERVICE_NAME has been installed and started."