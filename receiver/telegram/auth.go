package telegram

import (
	"catalog-go/config"
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func authenticate(update *models.Update) bool {
	// Check if in the list of allowed users
	for _, user := range config.CONFIG.Receiver.TelegramBot.PermittedUsers {
		if user == update.Message.From.ID {
			return true
		}
	}
	return false
}

func authenticateMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if !authenticate(update) {
			accessDeniedHandler(ctx, b, update)
			return
		}
		next(ctx, b, update)
	}
}
