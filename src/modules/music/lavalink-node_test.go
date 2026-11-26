package music

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
		_, _ = w.Write([]byte(`{"loadType":"search","data":{"tracks":[{"encoded":"abc","info":{"title":"Song","uri":"https://example.test"}}]}}`))
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
