package app

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loonbac/Ressy-Bot/src/modules/core"
)

// VoiceStateStore keeps the bot voice state reported by Discord.
type VoiceStateStore struct {
	mu      sync.Mutex
	byGuild map[string]core.VoiceState
}

// NewVoiceStateStore creates an empty store.
func NewVoiceStateStore() *VoiceStateStore {
	return &VoiceStateStore{byGuild: make(map[string]core.VoiceState)}
}

// Attach registers the Discord voice handlers on a session.
func (v *VoiceStateStore) Attach(session *discordgo.Session) {
	session.AddHandler(v.OnStateUpdate)
	session.AddHandler(v.OnServerUpdate)
}

func (v *VoiceStateStore) Lookup(guildID string) (core.VoiceState, bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	state, found := v.byGuild[guildID]
	return state, found
}

func (v *VoiceStateStore) OnStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	if session.State == nil || session.State.User == nil || event.UserID != session.State.User.ID || event.ChannelID == "" {
		return
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	state := v.byGuild[event.GuildID]
	state.SessionID = event.SessionID
	state.ChannelID = event.ChannelID
	v.byGuild[event.GuildID] = state
}

func (v *VoiceStateStore) OnServerUpdate(_ *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	v.mu.Lock()
	defer v.mu.Unlock()
	state := v.byGuild[event.GuildID]
	state.Token = event.Token
	state.Endpoint = event.Endpoint
	v.byGuild[event.GuildID] = state
}
