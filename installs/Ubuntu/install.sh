#!/bin/sh
cp -v ./godownload /usr/bin
cat > "/lib/systemd/system/godownload.service" << "EOF"
[Unit]
Description=Simple download service writen on golang
After=network.target
[Service]
ExecStart=/usr/bin/godownload
User=andrew
Group=andrew
KillMode=process
Restart=always
[Install]
WantedBy=network.target
EOF
