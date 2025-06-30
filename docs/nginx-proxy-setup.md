
setup a docker compose project folder and generate self signed certificates to be used by the nginx proxy

```
mkdir lidarr
cd lidarr
# create self signed certs for proxy
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
      -keyout lidarr.key \
      -out lidarr.pem \
      -subj "/CN=api.lidarr.audio"
```

create a nginx config for redirecting the lidarr metadata server requests
```
# nginx.conf
events { worker_connections 1024; }
http {
    server {
        listen 443 ssl;
        server_name api.lidarr.audio;

        # Self-signed cert just for the local connection
        ssl_certificate /etc/nginx/ssl/lidarr.pem;
        ssl_certificate_key /etc/nginx/ssl/lidarr.key;

        location / {
            proxy_pass https://api.musicinfo.pro/;
            proxy_set_header Host api.musicinfo.pro;
            proxy_ssl_server_name on;
            proxy_ssl_verify on;
            proxy_ssl_trusted_certificate /etc/ssl/certs/ca-certificates.crt;
        }
    }
}
```
if you don't want to use the blampe hosted metadata server you can update the following from the previous nginx.conf excerpt to point to your own self hosted api metadata server (not to be confused with the nginx proxy):
```
        location / {
            proxy_pass https://api.musicinfo.pro/;
            proxy_set_header Host api.musicinfo.pro;
        }
```

ie:

```
        location / {
            proxy_pass http://<my-metadata-server>/;
            proxy_set_header Host api.musicinfo.pro;
        }
```
see the lidarr metadata api repo to set that up https://github.com/Lidarr/LidarrAPI.Metadata

create a docker-compose.yaml with the following config:

```
networks:
  lidarr:

services:
  lidarr-api-proxy:
    image: nginx:latest
    volumes: 
      - ./lidarr.pem:/etc/nginx/ssl/lidarr.pem:ro
      - ./nginx.conf:/etc/nginx.conf:ro
      - ./lidarr.key:/etc/nginx/ssl/lidarr.key:ro
    networks:
      traefik:
        aliases:
          - api.lidarr.audio # override dns to our proxy for lidarr meatadata server dns record

  lidarr:
    image: lscr.io/linuxserver/lidarr:latest
    container_name: lidarr
    depends_on:
      - lidarr-api-proxy
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Etc/UTC
    volumes:
      - ./lidarr.pem:/etc/ssl/certs/lidarr.pem # add self signed lidarr pem to lidarr ssl cert trust store
      - /path/to/lidarr/config:/config
      - /path/to/music:/music #optional
      - /path/to/downloads:/downloads #optional
    ports:
      - 8686:8686
    networks:
      lidarr:
    restart: unless-stopped

```

instead of using docker networks with a dns alias like above you could instead do the following:
```
services:
  lidarr-api-proxy:
  ports:
    - 443:443

    lidarr:
      extra_hosts:
        - api.lidarr.audio:192.168.1.10
```
where 192.168.1.10 is the ip of your docker host.
