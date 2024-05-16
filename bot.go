package main

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var (
	_26381discordSession *discordgo.Session
)

func initDiscordSession() {
	bot, err := discordgo.New("Bot " + DISCORD_AUTH_TOKEN)
	if err != nil {
		slog.Error("discordgo.New", err)
		return
	}

	_26381discordSession = bot
}

func discordSession() *discordgo.Session {
	if _26381discordSession == nil {
		initDiscordSession()
		if _26381discordSession == nil {
			slog.Info("initDiscordSession", slog.String("reason", "unable to init discord session"))
			panic("unable to init discord session")
		}
	}

	return _26381discordSession
}
