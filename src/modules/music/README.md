# music

Provides YouTube playback through an embedded Lavalink 4 node.

Public API: `music.New(Config) core.Module`, plus the command constructors `NewPlayCommand`, `NewSkipCommand`, `NewQueueCommand`, `NewStopCommand` and `NewNowPlayingCommand`. `AudioNode` is the only transport boundary.

Allowed imports: the Go standard library, `gorilla/websocket` and `src/modules/core`. The module must not import `admin` or `app`.

Invariants: each guild has its own in-memory queue guarded by a mutex; `Register` adds five commands owned by `music`; `Close` closes the node WebSocket. The Lavalink constructor starts an event listener, so call it only when creating a module instance.

The adapter uses Lavalink 4 REST endpoints and reads the session id from the `ready` WebSocket op. Search uses `ytsearch:` prefixes; playback sends the Discord voice token, endpoint, session id and channel id.

To add a command, create a new file exporting `New<Name>Command`, register it in `music-module.go`, and add a test with a fake `AudioNode`.
