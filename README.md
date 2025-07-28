# 🦻 hearring-aid [![Discord](https://img.shields.io/discord/1367649771237675078?label=Discord)](https://discord.gg/Xykjv87yYs)

**Hear what your LiDAR is missing.**

This repo contains scripts for deploying and self-hosting a MusicBrainz instance in conjunction with L——'s metadata API.

> [!IMPORTANT]
> A hosted `api.musicinfo.pro` instance is provided as a service to the community. The software it's running is not my own, and I cannot guarantee its stability. If you'd like to help improve it, please reach out!

The process to use the provided `api.musicinfo.pro` instance varies depending on whether you're using Docker. Non-Docker users will need to update the SQLite database directly.

## 🚀 Docker Images

- The simplest way to use the hosted instance is by changing your L—— image to one of the [available tags](https://hub.docker.com/r/blampe/lidarr/tags).
  - This image is based on the upstream LinuxServer.io image, with a small nginx proxy that redirects metadata queries to `api.musicinfo.pro`.

- Alternatively, if you're already using L—— plugin builds (hot.io or LinuxServer.io), you can use the Tubifarry Plugin to set a custom metadata server.
  - [Section 10.1 of the self-hosted guide](https://github.com/blampe/hearring-aid/blob/main/docs/self-hosted-mirror-setup.md#101-configure-tubifarry-plugin-in-lidarr) provides detailed steps.
  - At step 8, set the metadata server URL to: `https://api.musicinfo.pro`

## 🛠️ Manual SQL (Non-Docker)

If you're not using Docker or don't want to use a forked image, you can override the default metadata server by inserting a row into your SQLite database:

```sql
INSERT INTO Config (Key, Value) VALUES ('metadatasource', 'https://api.musicinfo.pro/api/v0.4/');
```

To revert to the official metadata server:

```sql
DELETE FROM Config WHERE Key = 'metadatasource';
```

> [!TIP]
> You might also be interested in [rreading-glasses](https://github.com/blampe/rreading-glasses), which takes a similar approach for R——.

## 🧱 Self-Hosting

A step-by-step self-hosting guide is available using `docker-compose`. It walks through how to use the official [musicbrainz-docker](https://github.com/metabrainz/musicbrainz-docker) stack along with the L—— Metadata API image.

📖 [Read the full guide](https://github.com/blampe/hearring-aid/blob/main/docs/self-hosted-mirror-setup.md#101-configure-tubifarry-plugin-in-lidarr)

> **tl;dr:** Vanilla MusicBrainz images work just fine with the open-source metadata server — it's just a little tricky to configure.

## ⚠️ Disclaimer

This software is provided "as is", without warranty of any kind, express or implied, including but not limited to warranties of merchantability, fitness for a particular purpose, and noninfringement.

In no event shall the authors or copyright holders be liable for any claim, damages, or other liability, whether in an action of contract, tort, or otherwise, arising from or in connection with the software or the use or other dealings in the software.

This software is intended for educational and informational purposes only. It does not constitute legal, financial, or professional advice. The user assumes all responsibility for its use or misuse.

The user is free to use, modify, and distribute the software for any purpose, subject to the above disclaimers and conditions.
