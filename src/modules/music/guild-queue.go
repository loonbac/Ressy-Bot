package music

import "sync"

type guildQueues struct {
	mu      sync.Mutex
	waiting map[string][]Track
	current map[string]Track
}

func newGuildQueues() *guildQueues {
	return &guildQueues{
		waiting: make(map[string][]Track),
		current: make(map[string]Track),
	}
}

func (q *guildQueues) Enqueue(guildID string, track Track) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.waiting[guildID] = append(q.waiting[guildID], track)
}

func (q *guildQueues) Next(guildID string) (Track, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	tracks := q.waiting[guildID]
	if len(tracks) == 0 {
		return Track{}, false
	}
	q.current[guildID] = tracks[0]
	q.waiting[guildID] = tracks[1:]
	return tracks[0], true
}
func (q *guildQueues) SetCurrent(guildID string, track Track) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.current[guildID] = track
}

func (q *guildQueues) Clear(guildID string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.waiting, guildID)
	delete(q.current, guildID)
}

func (q *guildQueues) List(guildID string) []Track {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]Track(nil), q.waiting[guildID]...)
}

func (q *guildQueues) Current(guildID string) (Track, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	track, found := q.current[guildID]
	return track, found
}
