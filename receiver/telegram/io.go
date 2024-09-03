package telegram

import (
	"catalog-go/receiver/cat"
	"context"
	"io"
	"net/http"

	"github.com/go-telegram/bot"
)

func downloadFileToBuffer(link string) ([]byte, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func handleCat(b *bot.Bot, ctx context.Context, fileId string, caption string) error {
	f, err := b.GetFile(ctx, &bot.GetFileParams{
		FileID: fileId,
	})
	if err != nil {
		return err
	}

	link := b.FileDownloadLink(f)
	buf, err := downloadFileToBuffer(link)
	if err != nil {
		return err
	}

	return cat.Place(cat.File{
		Buffer: buf,
		Type:   "png",
	}, caption, ctx)
}
