FROM golang:1.27-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /out/plugins /out/lavalink/plugins \
 && CGO_ENABLED=1 go build -o /out/ressy-bot ./src/app/cmd \
 && CGO_ENABLED=1 go build -buildmode=plugin -o /out/plugins/music.so ./plugins/music \
 && CGO_ENABLED=1 go build -buildmode=plugin -o /out/plugins/admin.so ./plugins/admin \
 && curl -fsSL https://github.com/lavalink-devs/Lavalink/releases/download/4.2.2/Lavalink.jar -o /out/lavalink/Lavalink.jar \
 && curl -fsSL https://maven.lavalink.dev/releases/dev/lavalink/youtube/youtube-plugin/1.18.2/youtube-plugin-1.18.2.jar -o /out/lavalink/plugins/youtube-plugin-1.18.2.jar

FROM debian:bookworm-slim
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates curl ffmpeg gcc libc6-dev openjdk-17-jre-headless \
 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:${PATH}" CGO_ENABLED=1
WORKDIR /app
COPY --from=builder /out/ressy-bot ./ressy-bot
COPY --from=builder /out/plugins ./plugins
COPY --from=builder /out/lavalink ./lavalink
COPY application.yml ./lavalink/application.yml
COPY docker-entrypoint.sh ./docker-entrypoint.sh
RUN mkdir -p /app/data /app/lavalink/logs \
 && chmod +x /app/docker-entrypoint.sh
ENV RESSY_PLUGIN_DIR="/app/plugins" \
    RESSY_SETTINGS_PATH="/app/data/modules.json" RESSY_PANEL_ADDR=":8080" \
    RESSY_LAVALINK_URL="http://127.0.0.1:2333"
EXPOSE 8080 2333
ENTRYPOINT ["/app/docker-entrypoint.sh"]
