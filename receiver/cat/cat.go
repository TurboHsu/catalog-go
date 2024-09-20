package cat

import (
	"catalog-go/config"
	"catalog-go/database"
	"catalog-go/database/model"
	"catalog-go/database/query"
	"catalog-go/store"
	"context"
	"os"
	"time"

	"github.com/google/uuid"
)

func Place(image File, thumbnail File, caption string, ctx context.Context) (err error) {
	uuid := uuid.NewString()
	// Save file
	filename := config.CONFIG.Store.StorePath + "/" + uuid + "." + image.Type
	if _, err := os.Stat(config.CONFIG.Store.StorePath); os.IsNotExist(err) {
		os.Mkdir(config.CONFIG.Store.StorePath, 0755)
	}
	err = os.WriteFile(filename, image.Buffer, 0644)
	if err != nil {
		return err
	}

	thumbnailFilename := config.CONFIG.Store.StorePath + "/" + uuid + "." + "thumbnail" + "." + thumbnail.Type
	err = os.WriteFile(thumbnailFilename, thumbnail.Buffer, 0644)
	if err != nil {
		return err
	}

	cat := model.Cats{
		UUID:      uuid,
		Caption:   caption,
		Image:     uuid + "." + image.Type,
		Thumbnail: uuid + "." + "thumbnail" + "." + thumbnail.Type,
		CreatedAt: time.Now(),
	}
	c := query.Use(database.DB).Cats
	err = c.WithContext(ctx).Create(&cat)

	go func(filename string) {
		store.PutFileHook(filename)
		store.PutFileHook(thumbnailFilename)
	}(filename)

	return
}

type File struct {
	Buffer []byte
	Type   string // Extension
}