# MusicBrainz Mirror with Lidarr Metadata Server Setup Guide

This guide will help you deploy a local MusicBrainz mirror, [blampe's Lidarr Metadata Server](https://hub.docker.com/r/blampe/lidarr.metadata), and integrate it with Lidarr using the Tubifarry plugin, providing a working setup despite current Lidarr metadata issues. It walks through host and container setup, basic configuration, and validation steps to ensure the system is working properly.  If you already have a Lidarr instance but are not using lidarr-plugins, you will need to migrate to either ls.io or hot.io plugin branch (step 10 below walks through setting up a new lidarr container using ls.io)

> **Note:**  
> This guide is based on Debian 12 and Docker. It is provided as-is and without warranty, and your feedback and testing results are appreciated!

---

## Prerequisites

- Debian 12.11 server (root access)
- At least 8GB RAM, a moderately capable CPU, and 100GB of free disk space
- Basic familiarity with Docker and command line
- Internet connection

---

## 1. System Setup: Install Docker, Git, Screen, and Updates

```bash
# Add Docker's official GPG key and repository
apt-get update
apt-get install -y ca-certificates curl
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
apt-get update

# Install Docker, Docker Compose plugin, Git, Screen, and upgrade system
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin git screen
apt-get upgrade -y && apt-get dist-upgrade -y
```

---

## 2. Generate MusicBrainz Replication Token

1. Visit https://metabrainz.org/supporters/account-type and select your account type (usually "individual").
2. Go to https://metabrainz.org/profile and create an access token.
3. Save the 40-character alphanumeric token for later.

---

## 3. Setup MusicBrainz Docker Environment

```bash
mkdir -p /opt/docker && cd /opt/docker
git clone https://github.com/metabrainz/musicbrainz-docker.git
cd musicbrainz-docker
mkdir -p local/compose
```

### 3.1 Override Postgres User and Password

Create `local/compose/postgres-settings.yml`:

```yaml
services:
  musicbrainz:
    environment:
      POSTGRES_USER: "abc"
      POSTGRES_PASSWORD: "abc"
      MUSICBRAINZ_WEB_SERVER_HOST: "HOST_IP"  # Replace with your host IP
  db:
    environment:
      POSTGRES_USER: "abc"
      POSTGRES_PASSWORD: "abc"
  indexer:
    environment:
      POSTGRES_USER: "abc"
      POSTGRES_PASSWORD: "abc"
```

### 3.2 Customize Memory Settings

Create `local/compose/memory-settings.yml`:

```yaml
services:
  db:
    command: postgres -c "shared_buffers=2GB" -c "shared_preload_libraries=pg_amqp.so"
  search:
    environment:
      - SOLR_HEAP=2g
```

> These settings are a recommendation for a setup that will service 5 or fewer Lidarr instances; adjust memory allocations according to your server's capacity, and use case.

### 3.3 Customize Volume Paths

Create `local/compose/volume-settings.yml`:

```yaml
volumes:
  mqdata:
    driver_opts:
      type: none
      device: /opt/docker/musicbrainz-docker/volumes/mqdata
      o: bind
  pgdata:
    driver_opts:
      type: none
      device: /opt/docker/musicbrainz-docker/volumes/pgdata
      o: bind
  solrdata:
    driver_opts:
      type: none
      device: /opt/docker/musicbrainz-docker/volumes/solrdata
      o: bind
  dbdump:
    driver_opts:
      type: none
      device: /opt/docker/musicbrainz-docker/volumes/dbdump
      o: bind
  solrdump:
    driver_opts:
      type: none
      device: /opt/docker/musicbrainz-docker/volumes/solrdump
      o: bind
```

### 3.4 Configure Lidarr Metadata Server

Create `local/compose/lmd-settings.yml`:

```yaml
volumes:
  lmdconfig:
    driver_opts:
      type: none
      device: /opt/docker/musicbrainz-docker/volumes/lmdconfig
      o: bind
    driver: local

services:
  lmd:
    image: blampe/lidarr.metadata:70a9707
    ports:
      - 5001:5001
    environment:
      DEBUG: false
      PRODUCTION: false
      USE_CACHE: true
      ENABLE_STATS: false
      ROOT_PATH: ""
      IMAGE_CACHE_HOST: "theaudiodb.com"
      EXTERNAL_TIMEOUT: 1000
      INVALIDATE_APIKEY: ""
      REDIS_HOST: "redis"
      REDIS_PORT: 6379
      FANART_KEY: "xxx"        # Replace with your own key
      PROVIDERS__FANARTTVPROVIDER__0__0: "xxx"  # Replace with your own key
      SPOTIFY_ID: "xxx"            # Replace with your own key
      SPOTIFY_SECRET: "xxx"        # Replace with your own key
      SPOTIFY_REDIRECT_URL: "http://xxx.xxx.xxx.xxx:5001"   # set your host_ip
      PROVIDERS__SPOTIFYPROVIDER__1__CLIENT_ID: "xxx"    # Replace with your own key
      PROVIDERS__SPOTIFYPROVIDER__1__CLIENT_SECRET: "xxx" # Replace with your own key
      PROVIDERS__SPOTIFYAUTHPROVIDER__1__CLIENT_ID: "xxx" # Replace with your own key
      PROVIDERS__SPOTIFYAUTHPROVIDER__1__CLIENT_SECRET: "xxx" # Replace with your own key
      PROVIDERS__SPOTIFYAUTHPROVIDER__1__REDIRECT_URI: "http://xxx.xxx.xxx.xxx:5001"   # set your host_ip
      TADB_KEY: "2"    # Default, may need your own key for full functionality
      PROVIDERS__THEAUDIODBPROVIDER__0__0: "2"    # Default, may need your own key for full functionality
      LASTFM_KEY: "xxx"   # Replace with your own key
      LASTFM_SECRET: "xxx" # Replace with your own key
      PROVIDERS__SOLRSEARCHPROVIDER__1__SEARCH_SERVER: "http://search:8983/solr"
    restart: unless-stopped
    volumes:
      - lmdconfig:/config
    depends_on:
      - db
      - mq
      - search
      - redis
```

---

## 4. Create Volume Directories and Add Compose Overrides

```bash
mkdir -p volumes/{mqdata,pgdata,solrdata,dbdump,solrdump,lmdconfig}
admin/configure add local/compose/postgres-settings.yml local/compose/memory-settings.yml local/compose/volume-settings.yml local/compose/lmd-settings.yml
```

---

## 5. Build and Initialize MusicBrainz Database

```bash
docker compose build
docker compose run --rm musicbrainz createdb.sh -fetch   # This may take an hour or more
docker compose up -d
docker compose exec indexer python -m sir reindex --entity-type artist --entity-type release  # Indexing may take a couple hours
```

---

## 6. Schedule Weekly Index Updates

Edit `/etc/crontab` and add:

```
0 1 * * 7 root cd /opt/docker/musicbrainz-docker && /usr/bin/docker compose exec -T indexer python -m sir reindex --entity-type artist --entity-type release
```

---

## 7. Configure Replication Token and Start Replication

```bash
docker compose down
admin/set-replication-token   # Enter your replication token when prompted
admin/configure add replication-token
docker compose up -d
docker compose exec musicbrainz replication.sh   # Run initial replication; use screen to keep it running
admin/configure add replication-cron
docker compose down   # Wait for replication to finish before restarting
rm -rf volumes/dbdump/*   # Clean up, saves ~6GB
docker compose up -d
```

---

## 8. Initialize Lidarr Metadata Server Database

```bash
docker exec -it musicbrainz-docker-musicbrainz-1 /bin/bash
cd /tmp
git clone https://github.com/Lidarr/LidarrAPI.Metadata.git
psql postgres://abc:abc@db/musicbrainz_db -c 'CREATE DATABASE lm_cache_db;'
psql postgres://abc:abc@db/musicbrainz_db -f LidarrAPI.Metadata/lidarrmetadata/sql/CreateIndices.sql
exit
docker compose restart
```

---

## 9. Using the Lidarr Metadata Server

- Your Lidarr metadata server is available at: `http://host-ip:5001`

---

## 10. (IF NEEDED) Stand Up Lidarr-plugin container

```bash
cd /opt/docker && mkdir -p lidarr/volumes/lidarrconfig && cd lidarr
```

Create `docker-compose.yml`:

```yaml
services:
  lidarr:
    image: ghcr.io/linuxserver-labs/prarr:lidarr-plugins
    ports:
      - '8686:8686'
    environment:
      TZ: America/New_York
      PUID: 1000
      PGID: 1000
    volumes:
      - '/opt/docker/lidarr/volumes/lidarrconfig:/config'
      - '/mnt/media:/mnt/media'   # Adjust to your media path
    networks:
      - default

networks:
  default:
    driver: bridge
```

Start Lidarr:

```bash
docker compose up -d
```

---

### 10.1 Configure Tubifarry Plugin in Lidarr

1. Open your browser to `http://host_ip:8686` and complete initial setup.
2. Navigate to **System > Plugins**.
3. Install **Tubifarry prod plugin** by entering the URL:  
   `https://github.com/TypNull/Tubifarry`  
   and click **Install**.
4. After Lidarr restarts, go back to **System > Plugins**.
5. Install the **Tubifarry develop branch** plugin by entering:  
   `https://github.com/TypNull/Tubifarry/tree/develop`  
   and click **Install**.
6. After Lidarr restarts, log back into Lidarr and go to **Settings > Metadata**.
7. Under **Metadata Consumers**, select **Lidarr Custom**.
8. Check both boxes and enter your metadata server URL (e.g., `http://host_ip:5001`) in the **Metadata Source** field.
9. Save changes and restart Lidarr again:  
    ```bash
    docker compose restart
    ```

---

## 11. Verify and Troubleshoot

Follow these steps to test the setup and resolve common issues.

### ✅ Basic Functionality Test

- Open Lidarr and search for a new artist.  
  - If metadata loads correctly, your setup is working.
- Give Lidarr a restart to rule out any issues there.

### 🔁 Restart the Stack (Clean Restart)

If things aren't working as expected, try restarting the MusicBrainz stack:

```bash
cd /opt/docker/musicbrainz-docker
docker compose down && sleep 30 && docker compose up -d
```

> Give it a minute or two after restarting for services to fully initialize.

### 🌐 Test Metadata Server (lmd) in Browser

To confirm `lmd` is serving metadata properly, visit this URL in your browser (replace `host_ip` with your Docker host’s IP):

```
http://host_ip:5000/artist/1921c28c-ec61-4725-8e35-38dd656f7923
```

You should see a JSON response with artist details. If not:

- The `lmd` container may not be running correctly.
- The metadata may not be indexed yet.

### 📋 Check Logs for Errors

To troubleshoot further, view the `lmd` container logs:

```bash
docker logs -f musicbrainz-docker-lmd-1
```

Look for any errors or failed service messages that might indicate why the data isn't being returned.

---

## Final Notes

- Replace all placeholder keys (Spotify, Fanart.tv, Last.fm, etc.) with your own API keys.
- Adjust IP addresses and volume paths to match your environment.
- Initial DB creation, indexing, and replication can take several hours; be patient!
- Feel free to open an issue or comment if you find bugs or need help.

---

Enjoy your fully functional Lidarr instance again! 🎶🚀
