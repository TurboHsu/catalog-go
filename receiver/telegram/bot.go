package telegram

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
)

func Run(token string, ctx context.Context) {
	log.Printf("[I] Starting Telegram bot...\n")
	opts := []bot.Option{
		bot.WithMiddlewares(updateNilCatcherMiddleware),
		bot.WithMiddlewares(mediaGroupMiddleware),
		bot.WithDefaultHandler(defaultHandler),
	}
	b, err := bot.New(token, opts...)
	if err != nil {
		log.Printf("[E] Failed to create bot: %v\n", err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/id", bot.MatchTypeExact, idHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/ping", bot.MatchTypeExact, authenticateMiddleware(pingHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/post", bot.MatchTypeExact, authenticateMiddleware(postHandler))

	b.Start(ctx)
}
