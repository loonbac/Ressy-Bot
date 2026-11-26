package music

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestLavalinkNodeResolvesYouTubeSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/loadtracks" {
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("identifier"); got != "ytsearch:song" {
			t.Errorf("identifier = %q", got)
		}
		if got := r.Header.Get("Authorization"); got != "password" {
			t.Errorf("authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"loadType":"search","data":[{"encoded":"abc","info":{"identifier":"abc","isSeekable":true,"author":"Artist","length":180000,"title":"Song","uri":"https://example.test"},"pluginInfo":{}}]}`))
	}))
	defer server.Close()

	node := NewLavalinkNode(LavalinkConfig{BaseURL: server.URL, Password: "password", UserID: "bot"})
	defer node.Close()

	tracks, err := node.Resolve("ytsearch:song")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if len(tracks) != 1 || tracks[0].Encoded != "abc" || tracks[0].Title != "Song" {
		t.Fatalf("tracks = %+v", tracks)
	}
}

func TestLavalinkNodeResolvesPlaylist(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"loadType":"playlist","data":{"info":{"name":"Mix","selectedTrack":-1},"pluginInfo":{},"tracks":[{"encoded":"one","info":{"identifier":"one","isSeekable":true,"author":"A","length":1000,"title":"First","uri":"https://example.test/1"},"pluginInfo":{}},{"encoded":"two","info":{"identifier":"two","isSeekable":true,"author":"B","length":2000,"title":"Second","uri":"https://example.test/2"},"pluginInfo":{}}]}}`))
	}))
	defer server.Close()

	node := NewLavalinkNode(LavalinkConfig{BaseURL: server.URL, Password: "password", UserID: "bot"})
	defer node.Close()
	tracks, err := node.Resolve("https://example.test/playlist")
	if err != nil || len(tracks) != 2 || tracks[1].Encoded != "two" || tracks[1].Title != "Second" {
		t.Fatalf("tracks = %+v, err = %v", tracks, err)
	}
}

func TestLavalinkNodeStopsRetryingAfterClose(t *testing.T) {
	var dials atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { dials.Add(1) }))
	defer server.Close()

	node := NewLavalinkNode(LavalinkConfig{BaseURL: server.URL, Password: "password", UserID: "bot"})
	for deadline := time.Now().Add(time.Second); dials.Load() == 0 && time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
	}
	node.Close()
	time.Sleep(1500 * time.Millisecond)
	if got := dials.Load(); got != 1 {
		t.Fatalf("event loop dials after Close = %d, want the single pre-Close dial", got)
	}
}
