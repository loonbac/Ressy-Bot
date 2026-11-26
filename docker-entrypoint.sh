#!/usr/bin/env bash
set -euo pipefail

export LAVALINK_PASSWORD="${LAVALINK_PASSWORD:-changeme}"
export RESSY_LAVALINK_PASSWORD="$LAVALINK_PASSWORD" RESSY_BOT_USER_ID="${RESSY_BOT_USER_ID:-${DISCORD_APP_ID:-}}"
cd /app/lavalink
java -jar Lavalink.jar > logs/lavalink.log 2>&1 &
node_pid=$!
trap 'kill "$node_pid" 2>/dev/null || true' EXIT

for _ in $(seq 1 60); do
  if curl -sf -H "Authorization: $LAVALINK_PASSWORD" http://127.0.0.1:2333/version > /dev/null; then
    break
  fi
  sleep 1
done

cd /app
exec ./ressy-bot
