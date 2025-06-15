# syntax=docker.io/docker/dockerfile:1
ARG LIDARR_VARIANT="linuxserver-testing"
FROM ghcr.io/linuxserver/lidarr:version-2.12.4.4658 AS linuxserver-release
FROM ghcr.io/linuxserver/lidarr:develop-version-2.13.0.4664 AS linuxserver-testing
FROM ghcr.io/linuxserver/lidarr:nightly-version-2.13.1.4665 AS linuxserver-nightly
FROM ghcr.io/linuxserver-labs/prarr:lidarr-plugins-2.13.0.4661 AS linuxserver-plugins
FROM ghcr.io/hotio/lidarr:release-2.12.4.4658 AS hotio-release
FROM ghcr.io/hotio/lidarr:testing-2.13.0.4664 AS hotio-testing
FROM ghcr.io/hotio/lidarr:nightly-2.13.1.4665 AS hotio-nightly
FROM ghcr.io/hotio/lidarr:pr-plugins-2.13.0.4661 AS hotio-plugins
FROM ${LIDARR_VARIANT:?} AS main

ENV API_HOST="api.musicinfo.pro"
ENV API_PORT="443"
ENV API_SECURE="1"

SHELL ["/bin/sh", "-euc"]

RUN apk add --no-cache openssl haproxy

RUN install -DTm 0644 /dev/stdin /etc/haproxy/haproxy.cfg <<-'EOF'
	global
	  ca-base /etc/ssl/certs/
	  log stdout format raw local0

	defaults
	  timeout connect 5s
	  timeout client 10s
	  timeout server 10s
	  log global
	  option httplog

	crt-store
	  crt-base /etc/ssl/certs/
	  key-base /etc/ssl/private/
	  load crt "api.crt" key "api.key"

	frontend api-frontend
	  mode http
	  bind *:80
	  bind *:443 ssl crt "api.crt"
	  http-request redirect scheme https unless { ssl_fc }
	  default_backend api-backend

	backend api-backend
	  mode http
	  http-request set-header Host %[env(API_HOST)]:%[env(API_PORT)]
	  .if streq("${API_SECURE}",1)
	    server api ${API_HOST}:${API_PORT} ssl sni env(API_HOST) verify required ca-file ca-certificates.crt
	  .else
	    server api ${API_HOST}:${API_PORT}
	  .endif
EOF

RUN install -DTm 0644 /dev/stdin /etc/s6-overlay/s6-rc.d/haproxy/type <<-'EOF'
	longrun
EOF

RUN install -DTm 0755 /dev/stdin /etc/s6-overlay/s6-rc.d/haproxy/run <<-'EOF'
	#!/command/with-contenv /bin/sh
	set -eu

	if ! grep -Fwq 'api.lidarr.audio' /etc/hosts; then
	  printf '127.0.0.1 api.lidarr.audio\n' >> /etc/hosts
	fi

	if [ ! -f /etc/ssl/certs/api.crt ] || [ ! -f /etc/ssl/private/api.key ]; then
	  openssl ecparam -name prime256v1 -genkey -out /etc/ssl/private/api.key
	  openssl req -x509 -new -nodes -key /etc/ssl/private/api.key -days 36500 -out /etc/ssl/certs/api.crt -subj '/CN=api.lidarr.audio'
	  update-ca-certificates
	fi

	exec 2>&1
	exec haproxy -f /etc/haproxy/haproxy.cfg
EOF

RUN touch /etc/s6-overlay/s6-rc.d/user/contents.d/haproxy
