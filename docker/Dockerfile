FROM linuxserver/lidarr:2.12.4-develop

RUN apk update
RUN apk add --no-cache nginx openssl

COPY nginx.conf /etc/nginx/nginx.conf
COPY start.sh /start.sh
RUN chmod +x /start.sh

RUN openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout /etc/ssl/private/lidarr.key \
    -out /etc/ssl/certs/lidarr.pem \
    -subj "/CN=api.lidarr.audio"
