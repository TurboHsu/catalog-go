package telegram

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var mediaMap = map[string][]FileID{}

type FileID struct {
	Raw       string
	Thumbnail string
	Caption   string
}

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

	if update.Message.ReplyToMessage.MediaGroupID == "" {
		caption := update.Message.ReplyToMessage.Caption
		err := handleCat(b, ctx, update.Message.ReplyToMessage.Photo[len(update.Message.ReplyToMessage.Photo)-1].FileID,
			update.Message.ReplyToMessage.Photo[2].FileID, caption)
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
		caption := update.Message.ReplyToMessage.Caption
		defer delete(mediaMap, update.Message.ReplyToMessage.MediaGroupID)
		if !ok {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Failed to get file: media group not found in map",
			})
			log.Printf("[E] Failed to get file: media group not found in map\n")
			return
		}
		var wg sync.WaitGroup
		for _, id := range ids {
			wg.Add(1)
			go func(id FileID) {
				err := handleCat(b, ctx, id.Raw, id.Thumbnail, caption)
				if err != nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "Failed to get file",
					})
					log.Printf("[E] Failed to get file: %v\n", err)
					return
				}
				wg.Done()
			}(id)
		}

		wg.Wait()

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("%d cats has been posted.", len(ids)),
		})
	}
}

func accessDeniedHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "The cat god has denied your access to this command",
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
				FileID{
					Raw:       update.Message.Photo[len(update.Message.Photo)-1].FileID,
					Thumbnail: update.Message.Photo[2].FileID,
					Caption:   update.Message.Caption,
				})
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
