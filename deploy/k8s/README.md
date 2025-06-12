This directory contains example manifests that deploy MusicBrainz in
addition to the L—— metadata server.

The base manifests were auto-generated from the upstream
[docker-compose.yaml](https://github.com/metabrainz/musicbrainz-docker/blob/2fba21abb8117eb9a0cb853849c463a33a5060e7/docker-compose.yml)
using
[kompose](https://kubernetes.io/docs/tasks/configure-pod-container/translate-compose-kubernetes/).

The end result is not as Kubernetes-native as it could be, but it works.

The docker-compose stack requires a lot of scripts to get everything up and
running. The musicbrainz container isn't needed to serve traffic but is useful
for adhoc commands. The indexer container is also where you'll setup triggers
etc.

Many of the scripts don't translate well to k8s but are just running some
simple commands under the hood. For example, to generate and apply triggers in
the indexer container:
    ```
    python -m sir triggers --broker-id=1
    psql postgres://abc:abc@db/musicbrainz_db -f sql/CreateTriggers.sql
    psql postgres://abc:abc@db/musicbrainz_db -f sql/CreateFunctions.sql # Use DropFunctions.sql to disable live indexing.
    ```

The indexer's config was modified slightly to account for Lidarr's expected
"abc:abc" credentials.
