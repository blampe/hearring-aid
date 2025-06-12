#!/bin/sh

# Point api.lidarr.audio to localhost (nginx).
echo "127.0.0.1 api.lidarr.audio" >> /etc/hosts

nginx &

exec "/init"
