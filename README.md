# Steam Exporter
Exports your playtime metrics to prometheus. Losely based on [ichaelsergio/steam-exporter](https://github.com/michaelsergio/steam-exporter), but with improvements to make it more inline with how Prometheus works and updated steam apis.

## Prometheus
Everthing listens to /metric. You will be able to access it from `localhost:8080/metrics` when laoded.

Unlike other exporters, this exporter will pull the steam stats every time prometheus polls the metrics. It's recommended to have this on a relatively larger timer to not abuse the steam API.

### What it measures
Below are some of the measurements and their labels
- total steam games
    - steam_id
- steam total playtime in minutes
    - steam_id
    - app_id
    - name (_apps name_)
    - platform (_all, windows, mac, linux, deck_)
- last played timestamp
    - steam_id
    - app_id
    - name

### Grafana
A example grafana dashboard can be found in the grafana/dashboard.json

## Docker Image
[lachee/steam-exporter](https://hub.docker.com/repository/docker/lachee/steam-exporter)

## Environment Variables

| Name | Default | Description |
|------|---------|-------------|
| `STEAM_WEB_API_KEY`* | `null` | Your [Steam WebAPI Key](https://steamcommunity.com/dev). |
| `STEAM_ID64`* | `null` | The SteamID of the targeted user. The ID should be in the Int64 form. |
| `PORT` | `8080` | The port the metrics will be hosted on |

*required