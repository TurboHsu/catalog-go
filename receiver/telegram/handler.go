package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var mediaMap = map[string][]string{}

func idHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("Your ID is `%d`", userId),
		ParseMode: "Markdown",
	})
}

func pingHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Pong!",
	})
}

func postHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.ReplyToMessage == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "You need to reply to cat pics to use this command",
		})
		return
	}

	if update.Message.ReplyToMessage.Photo == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "You need to reply to cat pics to use this command",
		})
		return
	}

	caption := update.Message.ReplyToMessage.Caption

	if update.Message.ReplyToMessage.MediaGroupID == "" {
		err := handleCat(b, ctx, update.Message.ReplyToMessage.Photo[len(update.Message.ReplyToMessage.Photo)-1].FileID, caption)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Failed to get file",
			})
			log.Printf("[E] Failed to get file: %v\n", err)
			return
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "A cat has been posted",
		})
	} else { // Media group
		ids, ok := mediaMap[update.Message.ReplyToMessage.MediaGroupID]
		defer delete(mediaMap, update.Message.ReplyToMessage.MediaGroupID)
		if !ok {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Failed to get file: media group not found in map",
			})
			log.Printf("[E] Failed to get file: media group not found in map\n")
			return
		}
		for _, id := range ids {
			err := handleCat(b, ctx, id, caption)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Failed to get file",
				})
				log.Printf("[E] Failed to get file: %v\n", err)
				return
			}
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("%d cats has been posted.", len(ids)),
		})
	}
}

func accessDeniedHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "The cat god has denied you access to this command",
	})
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Bonjour, but I can't help eating cat!",
	})
}

func defaultHandler(_ context.Context, _ *bot.Bot, update *models.Update) {
	log.Printf("[I] Received message: [%d] %s\n", update.Message.From.ID, update.Message.Text)
}

func mediaGroupMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message.MediaGroupID != "" && update.Message.Photo != nil {
			mediaMap[update.Message.MediaGroupID] = append(mediaMap[update.Message.MediaGroupID],
				update.Message.Photo[len(update.Message.Photo)-1].FileID)
		}
		next(ctx, b, update)
	}
}

func updateNilCatcherMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}
		next(ctx, b, update)
	}
}
