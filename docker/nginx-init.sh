#!/bin/bash
echo starting nginx init....

if ! grep -q "api.lidarr.audio" /etc/hosts; then
echo "127.0.0.1 api.lidarr.audio" >> /etc/hosts
fi

/usr/sbin/nginx

