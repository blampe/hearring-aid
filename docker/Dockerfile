FROM linuxserver/lidarr:2.12.4-develop

RUN apk update
RUN apk add --no-cache nginx openssl

COPY nginx-init.sh /nginx-init
COPY nginx.conf /etc/nginx/nginx.conf

COPY s6-rc.d/nginx-proxy/ /etc/s6-overlay/s6-rc.d/nginx-proxy/
COPY s6-rc.d/user/contents.d/nginx-proxy /etc/s6-overlay/s6-rc.d/user/contents.d/

RUN chmod +x /nginx-init && \
    mkdir -p /etc/nginx/ssl && \
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout /etc/ssl/private/lidarr.key \
    -out /etc/ssl/certs/lidarr.pem \
    -subj "/CN=api.lidarr.audio"

