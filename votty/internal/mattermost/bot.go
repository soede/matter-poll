package mattermost

import (
	"github.com/mattermost/mattermost/server/public/model"
	"golang.org/x/exp/slog"
	"strings"
	"votty/internal/config"
)

type Bot struct {
	WebSocketClient *model.WebSocketClient
	APIv4Client     *model.Client4
}

func New(log *slog.Logger, cfg *config.Config) *Bot {
	client := model.NewAPIv4Client(cfg.MattermostURL)

	client.SetOAuthToken(cfg.BotToken)

	webSocketClient, err := model.NewWebSocketClient4(strings.Replace(cfg.MattermostURL, "http", "ws", 1), client.AuthToken)
	if err != nil {
		log.Error("Mattermost bot initialization failed", err)
		return nil
	}
	go webSocketClient.Listen()

	return &Bot{
		webSocketClient,
		client,
	}
}
