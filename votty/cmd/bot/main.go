package main

import (
	"log/slog"
	"votty/internal/app"
	"votty/internal/config"
	"votty/internal/logger"
	"votty/internal/mattermost"
	"votty/internal/storage/tarantool"
)

const version int = 1

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info(
		"starting votty-bot...",
		slog.Int("version", version),
	)
	log.Debug("debug messages are enabled")

	storage := tarantool.New(log, cfg)
	if storage == nil {
		return
	}

	bot := mattermost.New(log, cfg)
	if bot == nil {
		return
	}

	if err := app.NewApp(log, storage, bot).Run(); err != nil {
		log.Error("failed to start votty-bot.", err)
	}

}
