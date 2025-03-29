package app

import (
	"context"
	"github.com/mattermost/mattermost/server/public/model"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
	"votty/internal/handlers"
	"votty/internal/mattermost"
	"votty/internal/storage/tarantool"
)

type App struct {
	log       *slog.Logger
	tarantool *tarantool.Storage
	bot       *mattermost.Bot
}

func NewApp(log *slog.Logger, tarantool *tarantool.Storage, bot *mattermost.Bot) *App {
	return &App{log, tarantool, bot}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer a.bot.WebSocketClient.Close()
	defer a.tarantool.Conn.Close()

	bot, _, err := a.bot.APIv4Client.GetUser(ctx, "me", "")
	if err != nil {
		a.log.Error("Failed to retrieve bot data", err)
		return nil
	}

	botID := bot.Id

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case event := <-a.bot.WebSocketClient.EventChannel:
			if event.EventType() == model.WebsocketEventPosted {
				handlers.PostHandler(ctx, a.tarantool, a.log, a.bot.APIv4Client, event, botID)
			}
		case sig := <-quit:
			a.log.Info("Shutting down...", slog.String("Received signal", sig.String()))
			cancel()
			return nil
		}

	}

	return nil
}
