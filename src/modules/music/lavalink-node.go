package music

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type lavalinkNode struct {
	config    LavalinkConfig
	client    *http.Client
	mu        sync.RWMutex
	sessionID string
	conn      *websocket.Conn
}

// NewLavalinkNode creates a REST client and starts listening for node events.
func NewLavalinkNode(config LavalinkConfig) AudioNode {
	node := &lavalinkNode{config: config, client: &http.Client{Timeout: 10 * time.Second}}
	go node.followEvents()
	return node
}

func (n *lavalinkNode) session() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.sessionID
}

func (n *lavalinkNode) followEvents() {
	endpoint := strings.Replace(strings.TrimRight(n.config.BaseURL, "/"), "http", "ws", 1) + "/v4/websocket"
	headers := http.Header{
		"Authorization": []string{n.config.Password},
		"User-Id":       []string{n.config.UserID},
		"Client-Name":   []string{"Ressy/1"},
	}
	connection, _, err := websocket.DefaultDialer.Dial(endpoint, headers)
	if err != nil {
		return
	}
	n.mu.Lock()
	n.conn = connection
	n.mu.Unlock()
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			return
		}
		var event struct {
			Op        string `json:"op"`
			SessionID string `json:"sessionId"`
		}
		if json.Unmarshal(message, &event) == nil && event.Op == "ready" && event.SessionID != "" {
			n.mu.Lock()
			n.sessionID = event.SessionID
			n.mu.Unlock()
		}
	}
}

func (n *lavalinkNode) request(method, path string, body []byte) (*http.Response, error) {
	request, err := http.NewRequest(method, strings.TrimRight(n.config.BaseURL, "/")+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", n.config.Password)
	request.Header.Set("Content-Type", "application/json")
	return n.client.Do(request)
}

func (n *lavalinkNode) Resolve(query string) ([]Track, error) {
	path := "/v4/loadtracks?identifier=" + url.QueryEscape(query)
	response, err := n.request(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("load tracks returned %s", response.Status)
	}
	var result struct {
		LoadType string `json:"loadType"`
		Data     struct {
			Tracks []struct {
				Encoded string `json:"encoded"`
				Info    struct {
					Title string `json:"title"`
					URI   string `json:"uri"`
				} `json:"info"`
			} `json:"tracks"`
		} `json:"data"`
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		return nil, err
	}
	if result.LoadType == "empty" || result.LoadType == "error" {
		return nil, nil
	}
	tracks := make([]Track, 0, len(result.Data.Tracks))
	for _, item := range result.Data.Tracks {
		tracks = append(tracks, Track{Encoded: item.Encoded, Title: item.Info.Title, URI: item.Info.URI})
	}
	return tracks, nil
}

func (n *lavalinkNode) Play(guildID string, track Track, voice VoiceState) error {
	session := n.session()
	if session == "" {
		return ErrLavalinkSessionMissing
	}
	body, err := json.Marshal(map[string]any{
		"track": map[string]any{"encoded": track.Encoded},
		"voice": map[string]any{
			"token": voice.Token, "endpoint": voice.Endpoint,
			"sessionId": voice.SessionID, "channelId": voice.ChannelID,
		},
	})
	if err != nil {
		return err
	}
	path := "/v4/sessions/" + url.PathEscape(session) + "/players/" + url.PathEscape(guildID)
	response, err := n.request(http.MethodPatch, path, body)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("update player returned %s", response.Status)
	}
	return nil
}

func (n *lavalinkNode) Stop(guildID string) error {
	session := n.session()
	if session == "" {
		return ErrLavalinkSessionMissing
	}
	path := "/v4/sessions/" + url.PathEscape(session) + "/players/" + url.PathEscape(guildID)
	response, err := n.request(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusBadRequest && response.StatusCode != http.StatusNotFound {
		return fmt.Errorf("destroy player returned %s", response.Status)
	}
	return nil
}

func (n *lavalinkNode) NowPlaying(guildID string) (Track, bool) {
	session := n.session()
	if session == "" {
		return Track{}, false
	}
	path := "/v4/sessions/" + url.PathEscape(session) + "/players/" + url.PathEscape(guildID)
	response, err := n.request(http.MethodGet, path, nil)
	if err != nil {
		return Track{}, false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Track{}, false
	}
	var player struct {
		Track *struct {
			Encoded string `json:"encoded"`
			Info    struct {
				Title string `json:"title"`
				URI   string `json:"uri"`
			} `json:"info"`
		} `json:"track"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&player); err != nil || player.Track == nil {
		return Track{}, false
	}
	return Track{Encoded: player.Track.Encoded, Title: player.Track.Info.Title, URI: player.Track.Info.URI}, true
}

func (n *lavalinkNode) Close() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.conn == nil {
		return nil
	}
	return n.conn.Close()
}
