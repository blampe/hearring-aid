# 🦻 hearring-aid [![Discord](https://img.shields.io/discord/1367649771237675078?label=Discord)](https://discord.gg/Xykjv87yYs)

Hear what your LiDAR is missing.

This repo contains scripts for deploying and self-hosting a MusicBrainz instance
in conjunction with L——'s metadata API.

> [!IMPORTANT]
> A hosted `api.musicinfo.pro` instance is provided as a service to the
> community. The software it's running is not my own and I cannot guarantee its
> stability. If you'd like to help improve things please reach out!

To use the hosted instance, change your L—— image to any of the [available
tags](https://hub.docker.com/r/blampe/lidarr/tags).

This L—— image is derived from the upstream linuxserver image, with the
addition of a small nginx proxy responsible for re-directing metadata
queries to `api.musicinfo.pro`.

> [!TIP]
> You might also be interested in
> [rreading-glasses](http://github.com/blampe/rreading-glasses), which takes a
> similar approach with R——.

## Self-hosting

tl;dr: Vanilla MusicBrainz images work just fine with the open-source metadata
server, it's just a little tricky to configure.

> [!WARNING]
> Do **not** use any of the scripts, compose files, or images referenced in the
> L—— metadata source. They are outdated and incompatible with MusicBrainz
> schema changes.

* Deploy [metabrainz-docker](http://github.com/metabrainz/musicbrainz-docker)
  but use a username/password of abc/abc because the metadata server hard-codes
  it. Don't use the Solr dumps -- we don't need most of them.

  This is the hardest part, so make sure you follow the instructions closely.

  When interacting with `sir` to generate triggers and reindex data, always use
  `--entity-type artist --entity-type release` to dramatically reduce storage
  requirements from ~350GB to ~50GB.

* Deploy L—— metadata using the example in `deploy/compose` as a guide.

* Setup the L—— metadata server:
  * `psql postgres://abc:abc@localhost/musicbrainz_db -c 'CREATE DATABASE lm_cache_db;'`
  * `psql postgres://abc:abc@localhost/musicbrainz_db -f lidarrmetadata/sql/CreateIndices.sql`

The L—— server will crash repeatedly but seems to eventually stabilize (maybe).

## Disclaimer

This software is provided "as is", without warranty of any kind, express or
implied, including but not limited to the warranties of merchantability,
fitness for a particular purpose and noninfringement.

In no event shall the authors or copyright holders be liable for any claim,
damages or other liability, whether in an action of contract, tort or
otherwise, arising from, out of or in connection with the software or the use
or other dealings in the software.

This software is intended for educational and informational purposes only. It
is not intended to, and does not, constitute legal, financial, or professional
advice of any kind. The user of this software assumes all responsibility for
its use or misuse.

The user is free to use, modify, and distribute the software for any purpose,
subject to the above disclaimers and conditions.
