# Ressy Bot

Modular Discord bot in Go with hot-loadable `.so` plugins, a web admin panel and a music module backed by Lavalink. See `AGENTS.md` for repository conventions and `odd/tasks/ressy-bot-modular.md` for the active feature.

## Run

```sh
docker build -t ressy-bot .
docker run --rm -p 8080:8080 \
  -e DISCORD_BOT_TOKEN -e DISCORD_APP_ID \
  -e DISCORD_CLIENT_ID -e DISCORD_CLIENT_SECRET \
  -e RESSY_SESSION_SECRET -e RESSY_ADMIN_IDS \
  -e RESSY_PANEL_REDIRECT=https://your-host/callback \
  -e LAVALINK_PASSWORD ressy-bot
```

The image runs the Go host, the compiled `admin.so` and `music.so` plugins, Lavalink 4.2.2 and the YouTube source plugin 1.18.2 in one container. Java 17+, ffmpeg and the Go toolchain are present; a new plugin source can be compiled with `go build -buildmode=plugin` inside the container and dropped into `/app/plugins` for the panel to load.

YouTube availability depends on the source plugin's current clients and may need OAuth or a poToken; see `https://github.com/lavalink-devs/youtube-source`. No credentials are stored in this repository.
