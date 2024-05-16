package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	API_KEY = "troffcons"

	DISCORD_AUTH_TOKEN            = "MTIzODM4ODY3NTcyNjAyMDYxOA.GAyXcB.LSdGZeKJ8HH32d58tgOlVp9SpNBYoCxTd74Q38"
	DISCORD_BUG_REPORT_CHANNEL_ID = "1238391421174812693"
	DISCORD_GUILD_ID              = "1238391115871420427"
)

type JSON map[string]any

type SendNotificationRequest struct {
	DiscordNickname string `json:"discord_nickname"`
	Content         string `json:"content"`
	XApiKey         string `json:"x_api_key"`
}

type YandexCloudFunctionResponse struct {
	Body       any `json:"body"`
	StatusCode int `json:"statusCode"`
}

func YandexCloudFunctionHandler(ctx context.Context, req *SendNotificationRequest) ([]byte, error) {
	slog.Info("function is booting up", slog.Any("request", req))

	if !checkApiKey(req.XApiKey) {
		return respondJSON(http.StatusForbidden, JSON{
			"error":   "invalid request",
			"message": "x_api_key is invalid",
		})
	}

	initDiscordSession()
	slog.Info("initialized discord session")

	if req.Content == "" || req.DiscordNickname == "" {
		return respondJSON(http.StatusBadRequest, JSON{
			"error":   "invalid request",
			"message": "content and discord_nickname are required",
		})
	}

	ds := discordSession()

	members, err := ds.GuildMembers(DISCORD_GUILD_ID, "", 1000)
	if err != nil {
		slog.Error("ds.GuildMembers", err)
		return respondJSON(http.StatusInternalServerError, JSON{
			"error":   "failed to get members",
			"message": fmt.Sprintf("ds.GuildMembers: %s", err.Error()),
		})
	}

	found := false
	for _, u := range members {
		correctMember := u.User.Username == req.DiscordNickname
		if correctMember {
			found = true
			msg := fmt.Sprintf("%s\n%s", u.Mention(), req.Content)

			_, err = ds.ChannelMessageSend(DISCORD_BUG_REPORT_CHANNEL_ID, msg)
			if err != nil {
				slog.Error("ds.ChannelMessageSend", err)
				return respondJSON(http.StatusInternalServerError, JSON{
					"error":   "failed to send message to discord",
					"message": fmt.Sprintf("ds.ChannelMessageSend: %s", err.Error()),
				})
			}

			slog.Info("message info", slog.Any("sent_message", msg))
			break
		}

	}

	if !found {
		return respondJSON(http.StatusNotFound, JSON{
			"error":   "member not found",
			"message": fmt.Sprintf("member %s not found", req.DiscordNickname),
		})
	}

	return respondJSON(http.StatusOK, JSON{
		"result": "ok",
	})
}

func respondJSON(statusCode int, data JSON) ([]byte, error) {
	b, err := json.Marshal(YandexCloudFunctionResponse{
		Body:       data,
		StatusCode: statusCode,
	})
	if err != nil {
		slog.Error("json.Marshal", err)
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return b, nil
}

func checkApiKey(apiKey string) bool {
	return apiKey == API_KEY
}
