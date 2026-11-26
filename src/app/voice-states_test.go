package app

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func botSession() *discordgo.Session {
	return &discordgo.Session{State: &discordgo.State{Ready: discordgo.Ready{User: &discordgo.User{ID: "bot"}}}}
}

func TestVoiceStatesCollectSessionAndServer(t *testing.T) {
	states := NewVoiceStateStore()
	session := botSession()

	states.OnStateUpdate(session, &discordgo.VoiceStateUpdate{VoiceState: &discordgo.VoiceState{
		GuildID: "g1", UserID: "bot", ChannelID: "c1", SessionID: "session-1",
	}})
	states.OnServerUpdate(session, &discordgo.VoiceServerUpdate{
		GuildID: "g1", Token: "token-1", Endpoint: "voice.test",
	})

	state, found := states.Lookup("g1")
	if !found {
		t.Fatal("voice state was not stored")
	}
	if state.ChannelID != "c1" || state.SessionID != "session-1" || state.Token != "token-1" || state.Endpoint != "voice.test" {
		t.Fatalf("state = %+v", state)
	}
}

func TestVoiceStatesIgnoresOtherUsers(t *testing.T) {
	states := NewVoiceStateStore()
	session := botSession()

	states.OnStateUpdate(session, &discordgo.VoiceStateUpdate{VoiceState: &discordgo.VoiceState{
		GuildID: "g1", UserID: "someone-else", ChannelID: "c1", SessionID: "session-1",
	}})

	if _, found := states.Lookup("g1"); found {
		t.Fatal("voice state of another user was stored")
	}
}
